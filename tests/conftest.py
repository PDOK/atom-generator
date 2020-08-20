import pytest
import minio
from pathlib import Path

from atom_generator.config import Config
from atom_generator.minio_client import MinioDAO
from atom_generator.parser import ValuesParser


TEST_DATA = Path(__file__).parent / "data"


class MockMinio:
    class MockObject:
        def __init__(self, object_name="test.xml", is_dir=False, size=3):
            self.object_name = object_name
            self.is_dir = is_dir
            self.size = size

    mock_object = None

    def __init__(self, *args, **kwargs):
        pass

    @staticmethod
    def put_object(*args, **kwargs):
        pass

    def stat_object(self, bucket_name, object_name):
        return self.mock_object

    def list_objects(self, bucket_name, prefix="", recursive=False):
        return [self.mock_object]

    @staticmethod
    def remove_object(bucket_name, object_name):
        pass

    @staticmethod
    def copy_object(bucket_name, object_name, object_source):
        pass


local_env = {
    "S3_ACCESS_KEY": "miniostorage",
    "S3_SECRET_KEY": "miniostorage",
    "S3_SIGNING_REGION": "us-east-1",
    "S3_ENDPOINT_NO_PROTOCOL": "minio:9000",
    "NGR_ENVIRONMENT": "http://test",
}


@pytest.fixture
def _env_vars(monkeypatch):
    monkeypatch.setenv("S3_ACCESS_KEY", "access")
    monkeypatch.setenv("S3_SECRET_KEY", "secret")
    monkeypatch.setenv("S3_SIGNING_REGION", "test-region")
    monkeypatch.setenv("S3_ENDPOINT_NO_PROTOCOL", "endpoint")
    monkeypatch.setenv("NGR_ENVIRONMENT", "test")


@pytest.fixture
def _prepare_minio(_env_vars):
    return {"source_bucket": "source", "source_prefix": "/new/test/file"}


@pytest.fixture
def _prepare_minio_old(_prepare_minio):
    arguments = {"destination_bucket": "here", "destination_prefix": "/old/test/file"}
    arguments.update(_prepare_minio)
    return arguments


@pytest.fixture
def minio_new(monkeypatch, _prepare_minio):
    MockMinio.mock_object = MockMinio.MockObject()
    monkeypatch.setattr(minio, "Minio", MockMinio)
    return MinioDAO(**_prepare_minio)


@pytest.fixture
def minio_old(monkeypatch, _prepare_minio_old):
    MockMinio.mock_object = MockMinio.MockObject()
    monkeypatch.setattr(minio, "Minio", MockMinio)
    return MinioDAO(**_prepare_minio_old)


@pytest.fixture
def minio_object_is_dir(monkeypatch, _prepare_minio_old):
    MockMinio.mock_object = MockMinio.MockObject(is_dir=True)
    monkeypatch.setattr(minio, "Minio", MockMinio)
    return MinioDAO(**_prepare_minio_old)


@pytest.fixture
def minio_object_other_name(monkeypatch, _prepare_minio_old):
    MockMinio.stat_object = MockMinio.MockObject(object_name="other")
    monkeypatch.setattr(minio, "Minio", MockMinio)
    return MinioDAO(**_prepare_minio)


@pytest.fixture
def minio_empty(monkeypatch, _prepare_minio_old):
    MockMinio.mock_object = None
    monkeypatch.setattr(minio, "Minio", MockMinio)
    return MinioDAO(**_prepare_minio_old)


@pytest.fixture
def _prepare(tmpdir, _prepare_minio):
    tmp_path = str(tmpdir)

    arguments = {
        "config_path": str(TEST_DATA / "values.json"),
        "base_url": "http://example.com/base/",
        "force": False,
        "path": tmp_path,
        "service_path": "/test/path",
    }

    arguments.update(_prepare_minio)

    return arguments


@pytest.fixture
def _prepare_without_optional_parameters(tmpdir, _prepare_minio):
    arguments = {
        "config_path": str(TEST_DATA / "values.json"),
        "base_url": "http://example.com/base/",
        "force": False,
        "path": None,
        "service_path": "",
    }

    arguments.update(_prepare_minio)

    return arguments


@pytest.fixture
def get_config(_prepare, minio_new):
    config = Config(**_prepare)
    config.minio = minio_new
    return config


@pytest.fixture
def _prepare_old(_prepare, _prepare_minio_old):
    _prepare.update(_prepare_minio_old)
    _prepare["force"] = True
    return _prepare


@pytest.fixture
def _prepare_old_without_parameters(
    _prepare_without_optional_parameters, _prepare_minio_old
):
    _prepare_without_optional_parameters.update(_prepare_minio_old)
    _prepare_without_optional_parameters["force"] = True
    return _prepare_without_optional_parameters


@pytest.fixture
def get_config_old(_prepare_old, minio_old):
    config = Config(**_prepare_old)
    config.minio = minio_old
    return config


@pytest.fixture
def get_config_old_without_optional_parameters(
    _prepare_old_without_parameters, minio_old
):
    config = Config(**_prepare_old_without_parameters)
    config.minio = minio_old
    return config


@pytest.fixture
def values_new(get_config):
    return ValuesParser(get_config)


@pytest.fixture
def values_old(get_config_old):
    return ValuesParser(get_config_old)
