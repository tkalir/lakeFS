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
import datetime

from lakefs_sdk.models.external_principal_creation import ExternalPrincipalCreation  # noqa: E501

class TestExternalPrincipalCreation(unittest.TestCase):
    """ExternalPrincipalCreation unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional) -> ExternalPrincipalCreation:
        """Test ExternalPrincipalCreation
            include_option is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # uncomment below to create an instance of `ExternalPrincipalCreation`
        """
        model = ExternalPrincipalCreation()  # noqa: E501
        if include_optional:
            return ExternalPrincipalCreation(
                settings = [
                    {
                        'key' : ''
                        }
                    ]
            )
        else:
            return ExternalPrincipalCreation(
        )
        """

    def testExternalPrincipalCreation(self):
        """Test ExternalPrincipalCreation"""
        # inst_req_only = self.make_instance(include_optional=False)
        # inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()
