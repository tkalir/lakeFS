# coding: utf-8

"""
    lakeFS API

    lakeFS HTTP API

    The version of the OpenAPI document: 1.0.0
    Contact: services@treeverse.io
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


import unittest

from lakefs_sdk.api.import_api import ImportApi  # noqa: E501


class TestImportApi(unittest.TestCase):
    """ImportApi unit test stubs"""

    def setUp(self) -> None:
        self.api = ImportApi()  # noqa: E501

    def tearDown(self) -> None:
        pass

    def test_import_cancel(self) -> None:
        """Test case for import_cancel

        cancel ongoing import  # noqa: E501
        """
        pass

    def test_import_start(self) -> None:
        """Test case for import_start

        import data from object store  # noqa: E501
        """
        pass

    def test_import_status(self) -> None:
        """Test case for import_status

        get import status  # noqa: E501
        """
        pass


if __name__ == '__main__':
    unittest.main()
