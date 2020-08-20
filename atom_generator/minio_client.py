import warnings
import io
import mimetypes
import os
from datetime import timedelta

import minio
from remotezip import RemoteZip

from atom_generator.constants import (
    S3_SECRET_KEY,
    S3_SIGNING_REGION,
    S3_ENDPOINT_NO_PROTOCOL,
    S3_ACCESS_KEY,
    ADDITIONAL_MIME_TYPES,
    DEFAULT_CONTENT_TYPE,
    ZIP_MEDIA_TYPE,
    INSPIRE_GML_ZIP_MEDIA_TYPE,
    ZIP_MUST_CONTAIN_EXTENSION_CHECKLIST,
)
from atom_generator.util import build_uri
from atom_generator.error import AppError


def __add_mimetypes():
    """Adds configured mimetypes to """
    for extension, mimetype in ADDITIONAL_MIME_TYPES.items():
        mimetypes.add_type(mimetype, extension)


# run at import
__add_mimetypes()


class MinioDAO(object):
    """
    A Minio Database Access Object relative to a source and destination.
    """

    def __init__(
        self,
        source_bucket,
        source_prefix,
        destination_bucket=None,
        destination_prefix=None,
    ):
        self.source_bucket = source_bucket.strip("/")
        self.source_prefix = build_uri(source_prefix)
        self._client = minio.Minio(
            endpoint=os.environ[S3_ENDPOINT_NO_PROTOCOL],
            access_key=os.environ[S3_ACCESS_KEY],
            secret_key=os.environ[S3_SECRET_KEY],
            secure=False,
            region=os.environ[S3_SIGNING_REGION],
        )

        # TODO: remove destination and copy_mode after deprecation period
        self.destination_bucket = destination_bucket and destination_bucket.strip("/")
        self.destination_prefix = destination_prefix and destination_prefix.strip("/")
        self.copy_mode = destination_bucket is not None
        if self.copy_mode:
            # base_service is overwritten for easy refactoring when this is
            # deprecated.
            warnings.warn(
                "methods on a destination bucket will be deprecated in a future "
                "release",
                DeprecationWarning,
            )

    def source_media_type(self, filename, default=None):
        """
        Determine the source media type.

        Atom generator specific media types are also taken into account.

        Args:
            filename: filename in the minio source bucket without the source prefix.
        Returns:
            The guessed media type.
        """
        # guess based on the extension of the file.
        guess = mimetypes.guess_type(filename)[0]

        if guess is None:
            if default is not None:
                return default
            raise AppError(f"unknown filetype: {filename}")
        if guess == ZIP_MEDIA_TYPE:
            guess = self._handle_zip_file(guess, filename)

        return guess

    def _handle_zip_file(self, guess, filename):
        zip_contents = self._get_zip_file_list(filename)
        extension = self._source_in_zip_extension(zip_contents)

        if extension == "gml":
            return INSPIRE_GML_ZIP_MEDIA_TYPE

        return guess

    def _get_zip_file_list(self, filename):
        minio_key = self.source_prefix + "/" + filename
        url = self._client.presigned_get_object(
            self.source_bucket, minio_key, expires=timedelta(minutes=10)
        )

        with RemoteZip(url) as rz:
            file_list = rz.infolist()

        return file_list

    def _source_in_zip_extension(self, file_list):

        found_extensions = {s.filename.split(".")[-1] for s in file_list}
        intersection = found_extensions.intersection(
            ZIP_MUST_CONTAIN_EXTENSION_CHECKLIST
        )
        count = len(intersection)

        if count == 0:
            raise AppError(
                f"There must be one of the following file types in the zip package: {ZIP_MUST_CONTAIN_EXTENSION_CHECKLIST}"
            )
        elif count > 1:
            raise AppError(
                f"Found the following types '{intersection}', there must be only one of the following types in the "
                f"zip package: {ZIP_MUST_CONTAIN_EXTENSION_CHECKLIST} "
            )

        return intersection.pop()

    def source_object_size(self, filename):
        """
        Determines filesize for a file on the given minio source.

        The file will be taken from the bucket and prepended with the given
        source_prefix.

        Args:
            filename: filename in the minio source bucket without the source prefix.

        Returns:
            int: the size of the file in bytes.
        """
        if not filename.strip():
            raise ValueError("object_path can not be empty")

        key = build_uri(self.source_prefix, filename)
        object_stat = self._client.stat_object(self.source_bucket, key)

        if not object_stat:
            raise AppError(f"could not find: {key} in {self.source_bucket}")
        if object_stat.is_dir:
            raise AppError(
                f"object_path '{key} in {self.source_bucket}' points to directory -> should point to file"
            )
        return object_stat.size

    def __str__(self):
        return (
            f"Minio DAO from: {{ {self.source_bucket} }}/{{ {self.source_prefix} }}, "
            f"to: {{ {self.destination_bucket} }}/{{ {self.destination_prefix} }}"
        )

    # TODO methods below will be deprecated
    def destination_exists(self):
        """Checks if the given destination exists."""
        if self.destination_bucket is None:
            raise AppError("minio destination unknown")
        try:
            object_stat = self._client.stat_object(
                self.destination_bucket, self.destination_prefix
            )
        except minio.error.NoSuchKey:
            return False

        return object_stat is not None

    def rm_destination_tree(self):
        """Removes the destination from minio."""
        if self.destination_bucket is None:
            raise AppError("minio destination unknown")
        objects = self._client.list_objects(
            bucket_name=self.destination_bucket,
            prefix=self.destination_prefix,
            recursive=True,
        )
        for obj in objects:
            key = build_uri(self.destination_prefix, obj.object_name)
            self._client.remove_object(
                bucket_name=self.destination_bucket, object_name=key
            )

    def copy_from_source_to_destination(
        self, filename, source_prefix="", destination_prefix=""
    ):
        """
        Copies a file from source to destination.

        Both source and destination keys will be prepended with the MinioDAOs instance's
        given prefix (the MinioDAO's destination_prefix) and the prefix given to this
        function in case this is a nonempty string.

        Args:
            filename (str): filename in the minio source bucket given prefixes.
            source_prefix (str): additional prefix that is prepended to the filename
            destination_prefix (str): additional prefix that is prepended to the
                filename
        """
        if self.destination_bucket is None:
            raise AppError("minio destination unknown")
        source_object = build_uri(
            self.source_bucket, self.source_prefix, source_prefix, filename
        )
        destination_object = build_uri(
            self.destination_prefix, destination_prefix, filename
        )

        try:
            self._client.copy_object(
                self.destination_bucket, destination_object, source_object
            )
        except minio.error.NoSuchKey:
            raise AppError(
                f"failed copying minio object from {source_object} to "
                f"{destination_object}"
            )

    def save_to_destination(self, content, filename):
        """
        Stores content to Minio destination.

        The file is stored in the destination bucket. When stored, the filename is
        prefixed with the destination prefix.

        Args:
            content (str|bytes): the content to write to file.
            filename (str): the filename on minio this will be prefixed with this DAOs
                destination_prefix
        """
        output = io.BytesIO()
        output.write(content)
        output.seek(0, os.SEEK_END)
        target_length = output.tell()
        output.seek(0)

        self._client.put_object(
            self.destination_bucket,
            build_uri(self.destination_prefix, filename),
            output,
            target_length,
            content_type=self.source_media_type(filename, DEFAULT_CONTENT_TYPE),
        )
