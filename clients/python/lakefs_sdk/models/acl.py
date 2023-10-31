# coding: utf-8

"""
    lakeFS API

    lakeFS HTTP API

    The version of the OpenAPI document: 1.0.0
    Contact: services@treeverse.io
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


from __future__ import annotations
import pprint
import re  # noqa: F401
import json



from pydantic import BaseModel, Field, StrictStr

class ACL(BaseModel):
    """
    ACL
    """
    permission: StrictStr = Field(..., description="Permission level to give this ACL.  \"Read\", \"Write\", \"Super\" and \"Admin\" are all supported. ")
    __properties = ["permission"]

    class Config:
        """Pydantic configuration"""
        allow_population_by_field_name = True
        validate_assignment = True

    def to_str(self) -> str:
        """Returns the string representation of the model using alias"""
        return pprint.pformat(self.dict(by_alias=True))

    def to_json(self) -> str:
        """Returns the JSON representation of the model using alias"""
        return json.dumps(self.to_dict())

    @classmethod
    def from_json(cls, json_str: str) -> ACL:
        """Create an instance of ACL from a JSON string"""
        return cls.from_dict(json.loads(json_str))

    def to_dict(self):
        """Returns the dictionary representation of the model using alias"""
        _dict = self.dict(by_alias=True,
                          exclude={
                          },
                          exclude_none=True)
        return _dict

    @classmethod
    def from_dict(cls, obj: dict) -> ACL:
        """Create an instance of ACL from a dict"""
        if obj is None:
            return None

        if not isinstance(obj, dict):
            return ACL.parse_obj(obj)

        _obj = ACL.parse_obj({
            "permission": obj.get("permission")
        })
        return _obj


