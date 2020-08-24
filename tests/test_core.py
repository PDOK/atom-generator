from atom_generator.core import render_to_file, render, generate_atom_service
from atom_generator.constants import (
    DATA_FEED_TEMPLATE_NAME,
    SERVICE_FEED_TEMPLATE_NAME,
    TEMPLATES_DIR,
)


def test_render_object_property():
    class O:
        @property
        def dynamic_property(self):
            return "test"

    obj = O()
    template = "{{ dynamic_property }}"
    result = render(template, obj)
    assert result == "test"


def test_render_data_feed(values_new):
    values = values_new.parse()
    template = (TEMPLATES_DIR / DATA_FEED_TEMPLATE_NAME).read_text()
    empty_render = render(template, {})
    result = render(template, values.datasets[0])
    assert bool(result)
    assert result != template
    assert result != empty_render


def test_render_data_feed_described_by_link_non_inspire(values_new):
    values = values_new.parse()
    values.datasets[0].dataset_inspire_data_theme = ""

    template = (TEMPLATES_DIR / DATA_FEED_TEMPLATE_NAME).read_text()
    result = render(template, values.datasets[0])
    assert bool(result)

    described_by_link = (
        '<link rel="describedby" href="https://www.nationaalgeoregister.nl/geonetwork/srv/dut/csw?service=CSW&amp'
        ";version=2.0.2&amp;request=GetRecordById&amp;outputschema=http://www.isotc211.org/2005/gmd&amp"
        ';elementsetname=full&amp;id=81ff84ec-42a4-4481-840b-12713bbb5d38" type="text/xml"/>'
    )
    assert described_by_link in result


def test_render_data_feed_described_by_link_inspire(values_new):
    values = values_new.parse()

    template = (TEMPLATES_DIR / DATA_FEED_TEMPLATE_NAME).read_text()
    result = render(template, values.datasets[0])
    assert bool(result)

    described_by_link = '<link rel="describedby" href="https://inspire.ec.europa.eu/theme/hh" type="text/html"/>'
    assert described_by_link in result


def test_render_service_feed(values_new):
    values = values_new.parse()
    template = (TEMPLATES_DIR / SERVICE_FEED_TEMPLATE_NAME).read_text()
    empty_render = render(template, {})
    result = render(template, values)
    assert bool(result)
    assert result != empty_render
    assert result != template


def test_render_data_feed_to_file(values_new, tmp_path):
    values = values_new.parse()
    testfile = tmp_path / "test_datafeed"
    template = (TEMPLATES_DIR / DATA_FEED_TEMPLATE_NAME).read_text()
    render_to_file(testfile, template, values)
    result = testfile.read_text()
    assert bool(result)
    assert result != template


def test_render_service_feed_to_file(values_new, tmp_path):
    values = values_new.parse()
    testfile = tmp_path / "test_service_feed"
    template = (TEMPLATES_DIR / SERVICE_FEED_TEMPLATE_NAME).read_text()
    render_to_file(testfile, template, values)
    result = testfile.read_text()
    assert bool(result)
    assert result != template


def test_generate_atomfeed(get_config, tmpdir):
    expected = ["assets", "index.xml", "style", "top10nlv2_geografische_namen.xml"]
    generate_atom_service(get_config)
    assert sorted(f.basename for f in tmpdir.listdir()) == expected


def test_generate_old_style_atomfeed(get_config_old, mocker):
    save_spy = mocker.spy(get_config_old.minio, "save_to_destination")
    copy_spy = mocker.spy(get_config_old.minio, "copy_from_source_to_destination")

    generate_atom_service(get_config_old)
    assert save_spy.call_count == 90
    assert copy_spy.call_count == 1

    index_content, index_filename = save_spy.call_args_list[0].args
    assert index_filename == "index.xml"
    assert bool(index_content)

    dataset_content, dataset_filename = save_spy.call_args_list[1].args
    assert dataset_filename == "top10nlv2_geografische_namen.xml"
    assert bool(dataset_content)


def test_generate_old_style_atomfeed_without_optional_parameters(
    get_config_old_without_optional_parameters, mocker
):
    save_spy = mocker.spy(
        get_config_old_without_optional_parameters.minio, "save_to_destination"
    )
    copy_spy = mocker.spy(
        get_config_old_without_optional_parameters.minio,
        "copy_from_source_to_destination",
    )

    generate_atom_service(get_config_old_without_optional_parameters)
    assert save_spy.call_count == 90
    assert copy_spy.call_count == 1

    index_content, index_filename = save_spy.call_args_list[0].args
    assert index_filename == "index.xml"
    assert bool(index_content)

    dataset_content, dataset_filename = save_spy.call_args_list[1].args
    assert dataset_filename == "top10nlv2_geografische_namen.xml"
    assert bool(dataset_content)
