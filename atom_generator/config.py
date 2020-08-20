import os
from pathlib import Path
import re
from dataclasses import fields
from typing import get_args

from atom_generator.error import AppConfigError
from atom_generator.minio_client import MinioDAO
from atom_generator.util import build_uri
from atom_generator.constants import CSW_URL, CSW_ENVIRONMENT, ENV_VARIABLES


PARAMETER_REGEX = re.compile(r"\{\{\s?[\/\^\#\>]?\s?(\w+)\s?\}\}")


class Config(object):
    """Class to manage application configuration.

    Raises:
        EnvironmentError: Raised when expected env var not set
    """

    def __init__(
        self,
        config_path,
        base_url,
        force,
        path,
        service_path,
        source_bucket,
        source_prefix,
        destination_bucket=None,
        destination_prefix=None,
    ):
        self._check_environment()

        csw_env = os.environ[CSW_ENVIRONMENT]
        try:
            self.csw_base_url = CSW_URL[csw_env]
        except KeyError:
            raise AppConfigError(f"invalid csr environment { csw_env }")

        self.config_path = Path(config_path)
        if not self.config_path.exists():
            raise AppConfigError("config_path does not exist")

        self.service_url = (
            build_uri(base_url, service_path, ending_slash=True) if service_path else ""
        )
        self.force = force

        self.minio = MinioDAO(
            source_bucket=source_bucket,
            source_prefix=source_prefix,
            destination_bucket=destination_bucket,
            destination_prefix=destination_prefix,
        )

        if not self.minio.copy_mode:
            try:
                self.path = Path(path)
            except TypeError:
                raise AppConfigError("path does not exist")
            if not self.path.exists():
                raise AppConfigError("path does not exist")
        else:  # legacy, remove this in the future
            self.path = None
            self.service_url = build_uri(
                base_url, destination_prefix, ending_slash=True
            )

    @staticmethod
    def _check_environment():
        unset_vars = ENV_VARIABLES - os.environ.keys()
        if unset_vars:
            raise AppConfigError(
                f"environment variable(s) not set: {', '.join(unset_vars)}"
            )

    def __str__(self):
        return (
            "AtomGeneratorConfig:\n"
            f"  config path = {self.config_path}\n"
            f"  service url = {self.service_url}\n"
            f"  force = {self.force}\n"
            f"  path = {self.path}\n"
            f"  minio = {self.minio}"
        )


def validate_model(model, template):
    """
    Validates the input parameters in values.json to a json schema for a service type.

    Raises MapfileGeneratorConfigError when input_json is invalid.

    Args:
        service_type (ServiceType)
    """
    template_params = _template_fields(template)
    model_params = _model_fields(model)

    missing_from_model = [f"- {param}" for param in (template_params - model_params)]
    extra_in_model = [f"+ {param}" for param in (model_params - template_params)]
    message = "difference between model and template parameters:"
    if missing_from_model:
        message += "\n\nmissing from model:\n" + "\n".join(missing_from_model)
    if extra_in_model:
        message += "\n\nextra in model:\n" + "\n".join(extra_in_model)

    assert len(missing_from_model + extra_in_model) == 0, message


def _model_fields(klass):
    try:
        data_fields = fields(klass)
    except TypeError:
        return {}

    root_fields = [f.name for f in data_fields if f.repr]
    property_fields = _property_names(klass)
    model_class_fields = _flatten(_model_fields(f.type) for f in data_fields if f.repr)

    child_fields = _flatten(
        _model_fields(child_class)
        for field in data_fields
        for child_class in get_args(field.type)
    )
    return {
        f
        for f in set(root_fields + property_fields)
        .union(child_fields)
        .union(model_class_fields)
        if not f.startswith("_")
    }


def _template_fields(template):
    return set(PARAMETER_REGEX.findall(template))


def _flatten(iterable_of_iterables):
    return {x for y in iterable_of_iterables for x in y}


def _property_names(klass):
    return [p for p in dir(klass) if isinstance(getattr(klass, p), property)]
