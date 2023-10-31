"""
    lakeFS API

    lakeFS HTTP API  # noqa: E501

    The version of the OpenAPI document: 1.0.0
    Contact: services@treeverse.io
    Generated by: https://openapi-generator.tech
"""


import unittest

import lakefs_client
from lakefs_client.api.metadata_api import MetadataApi  # noqa: E501


class TestMetadataApi(unittest.TestCase):
    """MetadataApi unit test stubs"""

    def setUp(self):
        self.api = MetadataApi()  # noqa: E501

    def tearDown(self):
        pass

    def test_get_meta_range(self):
        """Test case for get_meta_range

        return URI to a meta-range file  # noqa: E501
        """
        pass

    def test_get_range(self):
        """Test case for get_range

        return URI to a range file  # noqa: E501
        """
        pass


if __name__ == '__main__':
    unittest.main()
