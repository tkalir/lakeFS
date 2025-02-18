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

from lakefs_sdk.api.branches_api import BranchesApi  # noqa: E501


class TestBranchesApi(unittest.TestCase):
    """BranchesApi unit test stubs"""

    def setUp(self) -> None:
        self.api = BranchesApi()  # noqa: E501

    def tearDown(self) -> None:
        pass

    def test_cherry_pick(self) -> None:
        """Test case for cherry_pick

        Replay the changes from the given commit on the branch  # noqa: E501
        """
        pass

    def test_create_branch(self) -> None:
        """Test case for create_branch

        create branch  # noqa: E501
        """
        pass

    def test_delete_branch(self) -> None:
        """Test case for delete_branch

        delete branch  # noqa: E501
        """
        pass

    def test_diff_branch(self) -> None:
        """Test case for diff_branch

        diff branch  # noqa: E501
        """
        pass

    def test_get_branch(self) -> None:
        """Test case for get_branch

        get branch  # noqa: E501
        """
        pass

    def test_list_branches(self) -> None:
        """Test case for list_branches

        list branches  # noqa: E501
        """
        pass

    def test_reset_branch(self) -> None:
        """Test case for reset_branch

        reset branch  # noqa: E501
        """
        pass

    def test_revert_branch(self) -> None:
        """Test case for revert_branch

        revert  # noqa: E501
        """
        pass


if __name__ == '__main__':
    unittest.main()
