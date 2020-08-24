import os
import datetime
from dataclasses import dataclass, field
from typing import List, Dict, Optional

from nested_dataclasses import nested

from atom_generator.minio_client import MinioDAO
from atom_generator.util import build_uri
from atom_generator.constants import (
    PROJECTIONS,
    DEFAULT_DOWNLOAD_TITLE,
    DOWNLOAD_TITLES,
)


@nested
@dataclass
class Download:
    download_file: str
    download_espg: str
    download_content: str = ""

    @property
    def has_download_content(self):
        return bool(self.download_content)

    @property
    def _minio(self):
        return self.parent.parent._minio

    @property
    def download_url(self):
        if self._minio.copy_mode:
            # TODO: deprecated
            return build_uri(
                self.parent.parent._service_url,
                "downloads",
                os.path.basename(self.download_file),
            )
        return build_uri(
            self.parent.parent._service_url, "downloads", self.download_file
        )

    @property
    def download_id(self):
        return self.download_url

    @property
    def download_title(self):
        # fallback value
        download_ext_str = DEFAULT_DOWNLOAD_TITLE
        for key in DOWNLOAD_TITLES:
            if self.download_file.endswith(key):
                download_ext_str = DOWNLOAD_TITLES[key]
        download_epsg = self.download_espg
        dataset_title = self.parent.datafeed_title_nl
        return f"{dataset_title} - {download_ext_str} (EPSG:{download_epsg})"

    @property
    def crs_uri(self):
        return f"http://www.opengis.net/def/crs/EPSG/0/{self.download_espg}"

    @property
    def crs_label(self):
        return PROJECTIONS.get(self.download_espg, "")

    @property
    def download_mimetype(self):
        return self._minio.source_media_type(self.download_file.lstrip("/"))

    @property
    def download_length(self):
        return self._minio.source_object_size(self.download_file.lstrip("/"))


@nested
@dataclass
class Dataset:
    datafeed_name: str
    datafeed_summary_nl: str
    datafeed_title_nl: str
    dataset_bbox: Dict[str, str]
    dataset_metadata_identifier: str
    dataset_inspire_data_theme: Optional[str]
    dataset_source_id: str
    dataset_source_id_ns: str
    dataset_rights: str
    datafeed_subtitle_nl: str
    downloads: List[Download]

    @property
    def dataset_polygon(self):
        return (
            "{miny} {minx} {miny} {maxx} {maxy} {maxx} {maxy} {minx} {miny} {minx}"
            "".format(**self.dataset_bbox)
        )

    @property
    def updated(self):
        return self.parent.updated

    @property
    def datafeed_url(self):
        return self.parent.feed_url(self.datafeed_name)

    @property
    def service_index_url(self):
        return self.parent.service_index_url

    @property
    def dataset_metadata_url(self):
        return self.parent.metadata_url(self.dataset_metadata_identifier)

    @property
    def dataset_metadata_web_url(self):
        return self.parent.metadata_web_url(self.dataset_metadata_identifier)


@nested
@dataclass
class ServiceFeed:
    service_title: str
    service_subtitle: str
    service_rights: str
    service_metadata_identifier: str
    datasets: List[Dataset]
    _ngr_base_url: str = field(repr=False)
    _service_url: str = field(repr=False)
    _minio: MinioDAO = field(repr=False)
    __updated: Optional[str] = field(repr=False, default=None)

    @property
    def updated(self):
        if self.__updated is None:
            self.__updated = f"{datetime.datetime.utcnow().isoformat()}Z"
        return self.__updated

    @property
    def service_index_url(self):
        return self.feed_url("index")

    @property
    def service_metadata_url(self):
        return self.metadata_url(self.service_metadata_identifier)

    @property
    def service_opensearch_url(self):
        return (
            f"{self._ngr_base_url}/geonetwork/opensearch/dut/"
            f"{self.service_metadata_identifier}/OpenSearchDescription.xml"
        )

    @property
    def service_metadata_web_url(self):
        return self.metadata_web_url(self.service_metadata_identifier)

    def feed_url(self, entry="index"):
        return build_uri(self._service_url, f"{entry}.xml")

    def metadata_url(self, metadata_id):
        return (
            f"{self._ngr_base_url}/geonetwork/srv/dut/csw?"
            f"service=CSW&"
            f"version=2.0.2&"
            f"request=GetRecordById&"
            f"outputschema=http://www.isotc211.org/2005/gmd&"
            f"elementsetname=full&"
            f"id={metadata_id}"
        )

    def metadata_web_url(self, metadata_id):
        return f"{self._ngr_base_url}/geonetwork/srv/dut/catalog.search#/metadata/{metadata_id}"
