"""
    lakeFS API

    lakeFS HTTP API  # noqa: E501

    The version of the OpenAPI document: 1.0.0
    Contact: services@treeverse.io
    Generated by: https://openapi-generator.tech
"""


import sys
import unittest

import lakefs_client
from lakefs_client.model.user import User
globals()['User'] = User
from lakefs_client.model.current_user import CurrentUser


class TestCurrentUser(unittest.TestCase):
    """CurrentUser unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def testCurrentUser(self):
        """Test CurrentUser"""
        # FIXME: construct object with mandatory attributes with example values
        # model = CurrentUser()  # noqa: E501
        pass


if __name__ == '__main__':
    unittest.main()
