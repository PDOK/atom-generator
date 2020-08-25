from datetime import datetime
import pytest
import mimetypes

from atom_generator.error import AppError


mimetypes.add_type("application/vnd.sqlite3", ".test_sqlite")
mimetypes.add_type("application/x-sqlite3", ".test_sqlite_legacy")


def test_source_media_type(minio_new):
    minio_dao = minio_new
    media_type = minio_dao.source_media_type("test.gpkg")
    assert media_type == "application/geopackage+sqlite3"


def test_source_media_type_vendor(minio_new):
    minio_dao = minio_new
    media_type = minio_dao.source_media_type("test.test_sqlite")
    assert media_type == "application/vnd.sqlite3"


def test_source_media_type_legacy_media_type(minio_new):
    minio_dao = minio_new
    media_type = minio_dao.source_media_type("test.test_sqlite_legacy")
    assert media_type == "application/x-sqlite3"


def test_source_media_type_unknown_raises_app_error(minio_empty):
    minio_dao = minio_empty
    pytest.raises(AppError, minio_dao.source_media_type, "test.does_not_exist")


def test_source_object_date(minio_new):
    minio_dao = minio_new
    last_modified = minio_dao.source_object_date("test.gpkg")
    assert last_modified == "2020-08-25T13:45:00Z"


def test_source_object_date_format(minio_new):
    minio_dao = minio_new
    last_modified = minio_dao.source_object_date("test.gpkg")
    datetime.strptime(last_modified, "%Y-%m-%dT%H:%M:%SZ")


def test_source_object_size(minio_new):
    minio_dao = minio_new
    size = minio_dao.source_object_size("test.gpkg")
    assert size == 3


def test_source_object_size_empty(minio_empty):
    minio_dao = minio_empty
    pytest.raises(AppError, minio_dao.source_object_size, "test.gpkg")


def test_source_object_size_is_dir(minio_object_is_dir,):
    minio_dao = minio_object_is_dir
    pytest.raises(AppError, minio_dao.source_object_size, "test/")


# TODO tests below will be deprecated
def test_save_to_destination(minio_old, mocker):
    minio_dao = minio_old
    _input = (b"test", "test.xml")
    spy = mocker.spy(minio_dao._client, "put_object")
    minio_dao.save_to_destination(*_input)
    bucket_name, object_name, _, length = spy.call_args.args
    assert bucket_name == "here"
    assert object_name == "old/test/file/test.xml"
    assert length == 4
    assert spy.call_args.kwargs["content_type"] == "application/xml"


def test_destination_exists(minio_old):
    minio_dao = minio_old
    exists = minio_dao.destination_exists()
    assert exists is True


def test_destination_does_not_exist_on_minio(minio_empty):
    minio_dao = minio_empty
    exists = minio_dao.destination_exists()
    assert exists is False


def test_rm_destination_tree(minio_old, mocker):
    minio_dao = minio_old
    spy = mocker.spy(minio_dao._client, "remove_object")
    minio_dao.rm_destination_tree()
    assert spy.call_args.kwargs["bucket_name"] == "here"
    assert spy.call_args.kwargs["object_name"] == "old/test/file/test.xml"


def test_copy_from_source_to_destination(minio_old, mocker):
    minio_dao = minio_old
    spy = mocker.spy(minio_dao._client, "copy_object")
    minio_dao.copy_from_source_to_destination("test.xml")
    bucket_name, object_name, object_source = spy.call_args.args
    assert bucket_name == "here"
    assert object_name == "old/test/file/test.xml"
    assert object_source == "source/new/test/file/test.xml"


def test_destination_is_not_given_raises_app_error(minio_new):
    minio_dao = minio_new
    pytest.raises(AppError, minio_dao.destination_exists)
    pytest.raises(AppError, minio_dao.rm_destination_tree)
    pytest.raises(AppError, minio_dao.copy_from_source_to_destination, "test")
