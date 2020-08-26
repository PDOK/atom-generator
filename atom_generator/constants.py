from pathlib import Path

RESOURCES_DIR = Path(__file__).parent / "resources"
TEMPLATES_DIR = RESOURCES_DIR / "templates"
STATIC_FILES_DIR = RESOURCES_DIR / "static_files"
SERVICE_FEED_TEMPLATE_NAME = "service_feed.mustache"
DATA_FEED_TEMPLATE_NAME = "data_feed.mustache"

S3_SECRET_KEY = "S3_SECRET_KEY"
S3_SIGNING_REGION = "S3_SIGNING_REGION"
S3_ENDPOINT_NO_PROTOCOL = "S3_ENDPOINT_NO_PROTOCOL"
NGR_ENVIRONMENT = "NGR_ENVIRONMENT"
S3_ACCESS_KEY = "S3_ACCESS_KEY"
ENV_VARIABLES = {
    S3_ACCESS_KEY,
    S3_SECRET_KEY,
    S3_SIGNING_REGION,
    S3_ENDPOINT_NO_PROTOCOL,
    NGR_ENVIRONMENT,
}

NGR_URL = {
    "prod": "https://www.nationaalgeoregister.nl",
    "test": "https://www.ngr.test",
}

DEFAULT_CONTENT_TYPE = "application/octet-stream"
ZIP_MEDIA_TYPE = "application/zip"
INSPIRE_GML_ZIP_MEDIA_TYPE = "application/x-gmz"

# These are the more
ADDITIONAL_MIME_TYPES = {
    ".gml": "application/gml+xml",
    ".gpkg": "application/geopackage+sqlite3",
}

DEFAULT_DOWNLOAD_TITLE = "Download"
DOWNLOAD_TITLES = {
    "gpkg.zip": "Zipped GeoPackage download",
    "gml.zip": "Zipped GML download",
    "gpkg": "GeoPackage download",
    "gml": "GML download",
}

ZIP_MUST_CONTAIN_EXTENSION_CHECKLIST = {"gml", "gpkg", "xml"}

MD_NAMESPACES = {
    "gmd": "http://www.isotc211.org/2005/gmd",
    "gco": "http://www.isotc211.org/2005/gco",
    "gmx": "http://www.isotc211.org/2005/gmx",
}

PROJECTIONS = {
    "28992": "Amersfoort / RD New",
    "3035": "ETRS89 / LAEA Europe",
    "7415": "Amersfoort / RD New + NAP height",
    "4326": "WGS 84",
    "3857": "WGS 84 / Pseudo-Mercator",
    "4230": "ED50",
    "23031": "ED50 / UTM zone 31N",
    "23032": "ED50 / UTM zone 32N",
    "4258": "ETRS89",
    "4937": "ETS89 (3D)",
    "3034": "ETRS89 / LCC Europe",
    "4979": "WGS 84 (3D)",
    "32631": "WGS 84 / UTM zone 31N",
    "32632": "WGS 84 / UTM zone 32N",
}
