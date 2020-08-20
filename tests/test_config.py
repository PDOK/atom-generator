# TODO: remove old tests when deprecation period is over
def test_config_old(get_config_old):
    config = get_config_old
    assert config.minio.copy_mode is True
    assert config.minio.destination_bucket == "here"
    assert config.minio.destination_prefix == "old/test/file"


def test_config(get_config):
    config = get_config
    assert config.minio.copy_mode is False
    assert config.minio.source_bucket == "source"
    assert config.minio.source_prefix == "new/test/file"


# TODO: remove old tests when deprecation period is over
def test_get_base_service_url_old(get_config_old):
    config = get_config_old
    assert config.service_url == "http://example.com/base/old/test/file/"


def test_get_base_service_url(get_config):
    config = get_config
    assert config.service_url == "http://example.com/base/test/path/"
