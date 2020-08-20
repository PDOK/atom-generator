import json

from tests.conftest import TEST_DATA


def test_values_parser(values_new):
    expected = json.loads((TEST_DATA / "expected_values.json").read_text())
    values = values_new.parse().to_dict()
    del values["updated"]
    del values["datasets"][0]["updated"]
    assert values == expected


def test_values_parser_old(values_old):
    expected = json.loads((TEST_DATA / "expected_values_old.json").read_text())
    values = values_old.parse().to_dict()
    del values["updated"]
    del values["datasets"][0]["updated"]
    assert values == expected
