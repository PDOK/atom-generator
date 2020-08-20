import pystache
import logging
import os
from shutil import copytree
from nested_dataclasses import ValidationError

from atom_generator.constants import (
    DATA_FEED_TEMPLATE_NAME,
    SERVICE_FEED_TEMPLATE_NAME,
    TEMPLATES_DIR,
    STATIC_FILES_DIR,
)
from atom_generator.error import AppError, AppConfigError
from atom_generator.parser import ValuesParser


logger = logging.getLogger(__name__)


def render(template, values):
    """
    Render a mustache template.

    Args:
        template (str): a file-like object or a string containing the template
        values (dataclass): a python dataclass with the data scope to render a template

    Returns:
        str: the rendered result.
    """
    renderer = pystache.Renderer()
    return renderer.render(template, values)


def render_to_file(path, template, values):
    """
    Render a mustache template.

    Args:
        path (pathlib.Path): path to write to
        template (str): A file-like object or a string containing the template
        values (dict): A python dictionary with the data scope
    """
    path.write_text(render(template, values))


def generate_atom_service(config):
    """
    Generate Atom service feed.

    Stores the rendered atom index and entry xmls in the config.path.

    When the config is initialized with a destination bucket and prefix, the source
    files are copied from source to destination in minio. This behaviour will be
    deprecated in a future release.

    Args:
        config(atom_generator.config.Config): atom generator config.
    """
    service_feed_template = (TEMPLATES_DIR / SERVICE_FEED_TEMPLATE_NAME).read_text()
    data_feed_template = (TEMPLATES_DIR / DATA_FEED_TEMPLATE_NAME).read_text()
    parser = ValuesParser(config)
    service_feed = parser.parse()
    try:
        service_feed.validate()
    except ValidationError as e:
        raise AppConfigError(e)

    if not config.minio.copy_mode:
        render_to_file(
            path=config.path / "index.xml",
            values=service_feed,
            template=service_feed_template,
        )
        copytree(str(STATIC_FILES_DIR), str(config.path), dirs_exist_ok=True)
        for entry in service_feed.datasets:
            render_to_file(
                path=config.path / f"{entry.datafeed_name}.xml",
                values=entry,
                template=data_feed_template,
            )
    else:
        # TODO this will be deprecated.
        index_xml = render(service_feed_template, service_feed).encode("utf-8")
        logger.warning("using old style atom-generator input")
        if config.minio.destination_exists():
            if config.force:
                config.minio.rm_destination_tree()
            else:
                raise AppError(
                    f"atom {config.minio.destination_prefix} already exists, "
                    f"use --force to overwrite the existing atom "
                )

        config.minio.save_to_destination(index_xml, "index.xml")
        for entry in service_feed.datasets:
            datafeed_filename = f"{entry.datafeed_name}.xml"
            xml = render(data_feed_template, entry).encode("utf-8")
            config.minio.save_to_destination(xml, datafeed_filename)
            for download in entry.downloads:
                config.minio.copy_from_source_to_destination(
                    download.download_file, destination_prefix="downloads"
                )
            here = os.getcwd()
            os.chdir(str(STATIC_FILES_DIR))
            for root, _, files in os.walk("."):
                for file in files:
                    # root[:2] is always "./"
                    file_path = os.path.join(root[2:], file)
                    with open(file_path, "rb") as f:
                        content = f.read()
                        config.minio.save_to_destination(content, file_path)
            os.chdir(here)

    logger.info("created atom feed for %s", service_feed.service_index_url)
