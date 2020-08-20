import json

import dacite

from atom_generator import models


class ValuesParser:
    """Parses Atom feed input and sets dynamic fields."""

    MODEL = models.ServiceFeed

    def __init__(self, config):
        self.minio = config.minio
        self.path = config.config_path
        self.csw_base_url = config.csw_base_url
        self.service_url = config.service_url

    def parse(self):
        """
        Parse the contents of a values.json to a service dataclass.

        Args:
            values (dict): a python dictionary with the data of the service
        Returns
            dataclass: a python dataclass objects with all relevant data to render a
                service
        """
        values = json.loads(self.path.read_text())

        return dacite.from_dict(
            data_class=self.MODEL,
            data=dict(
                _csw_base_url=self.csw_base_url,
                _service_url=self.service_url,
                _minio=self.minio,
                **values
            ),
            config=dacite.Config(check_types=False, strict=False),
        )
