import logging
import sys

import click
import click_log

logger = logging.getLogger(__name__)
click_log.basic_config(logger)

from atom_generator.core import generate_atom_service
from atom_generator.config import Config, validate_model
from atom_generator.constants import (
    TEMPLATES_DIR,
    SERVICE_FEED_TEMPLATE_NAME,
    DATA_FEED_TEMPLATE_NAME,
)
from atom_generator.models import ServiceFeed, Dataset
from atom_generator.error import AppError, AppConfigError


@click.group()
def cli():
    pass


@cli.command(name="gen-atom-service")
@click_log.simple_verbosity_option(logger)
@click.argument("locations", nargs=-1)
@click.argument("config-path", type=click.Path(exists=True))
@click.argument("base-url")
@click.option("--force/--no-force", default=False, help="Overwrite existing atom feed")
@click.option(
    "--path",
    default=None,
    help="Path where atom xml should be stored locally",
    type=click.Path(exists=True),
)
@click.option("--service_path", default="", help="prefix path for atom url")
def generate_atom_service_command(
    locations, config_path, base_url, force, path, service_path
):
    """Generate Atom service feed.

    [locations]... should be filled with the following
    parameters (* is optional):

    source_bucket       bucket in minio that contains the raw files
    source_path         path to files inside the minio bucket
    *destination_bucket optional new location for the raw files
    *destination_path   path to files iniside the minio bucket

    example: `generate-atom source_bucket /source_path destination_bucket /destination_path conf.json http://example.com`
    """

    if len(locations) == 2:
        source_bucket, source_path = locations
        destination_bucket, destination_path = None, None
    elif len(locations) == 4:
        source_bucket, source_path, destination_bucket, destination_path = locations
    else:
        logger.error("Locations takes only 2 or 4 arguments")
        sys.exit(1)

    try:
        app_config = Config(
            config_path,
            base_url,
            force,
            path,
            service_path,
            source_bucket,
            source_path,
            destination_bucket,
            destination_path,
        )
    except AppConfigError:
        logger.exception("Atom generator config failed:")
        sys.exit(1)

    logger.info(
        "\nGenerating atom with following configuration:\n\n %s", str(app_config)
    )

    try:
        generate_atom_service(app_config)
    except AppError:
        logger.exception("Atom generator failed:")
        sys.exit(1)


@cli.command(name="validate_models")
@click_log.simple_verbosity_option(logger)
def validate_models():
    """
    Validate models. Check if the variables in the models correspond with the templates.
    """
    no_errors = True
    compare_these = {
        "service feed": (
            (TEMPLATES_DIR / SERVICE_FEED_TEMPLATE_NAME).read_text(),
            ServiceFeed,
        ),
        "data feed": ((TEMPLATES_DIR / DATA_FEED_TEMPLATE_NAME).read_text(), Dataset),
    }

    for name, (template, model) in compare_these.items():
        try:
            validate_model(model, template)
        except AssertionError as e:
            logger.error("%s: %s", name, e)
            no_errors = False
    if no_errors:
        logger.info("models are in sync")


if __name__ == "__main__":
    cli()
