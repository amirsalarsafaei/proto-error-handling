# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: helloworld/helloworld.proto
# Protobuf Python Version: 5.29.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    2,
    '',
    'helloworld/helloworld.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x1bhelloworld/helloworld.proto\x12\x0bhello_world\"E\n\x11\x43reateUserRequest\x12\x1a\n\x08username\x18\x01 \x01(\tR\x08username\x12\x14\n\x05\x65mail\x18\x02 \x01(\tR\x05\x65mail\"^\n\x12\x43reateUserResponse\x12\x17\n\x07user_id\x18\x01 \x01(\tR\x06userId\x12/\n\x06status\x18\x02 \x01(\x0e\x32\x17.hello_world.UserStatusR\x06status\"\x87\x01\n\x15\x43reateUserAltResponse\x12\x31\n\x07success\x18\x01 \x01(\x0b\x32\x15.hello_world.UserDataH\x00R\x07success\x12\x31\n\x05\x65rror\x18\x02 \x01(\x0b\x32\x19.hello_world.ErrorDetailsH\x00R\x05\x65rrorB\x08\n\x06result\"T\n\x08UserData\x12\x17\n\x07user_id\x18\x01 \x01(\tR\x06userId\x12/\n\x06status\x18\x02 \x01(\x0e\x32\x17.hello_world.UserStatusR\x06status\"<\n\x0c\x45rrorDetails\x12\x12\n\x04\x63ode\x18\x01 \x01(\tR\x04\x63ode\x12\x18\n\x07message\x18\x02 \x01(\tR\x07message*Z\n\nUserStatus\x12\x1b\n\x17USER_STATUS_UNSPECIFIED\x10\x00\x12\x16\n\x12USER_STATUS_ACTIVE\x10\x01\x12\x17\n\x13USER_STATUS_PENDING\x10\x02\x32\xb5\x01\n\x0bUserService\x12O\n\nCreateUser\x12\x1e.hello_world.CreateUserRequest\x1a\x1f.hello_world.CreateUserResponse\"\x00\x12U\n\rCreateUserAlt\x12\x1e.hello_world.CreateUserRequest\x1a\".hello_world.CreateUserAltResponse\"\x00\x42\xb4\x01\n\x0f\x63om.hello_worldB\x0fHelloworldProtoP\x01ZHgithub.com/amirsalarsafaei/proto-error-handling/go/helloworld;helloworld\xa2\x02\x03HXX\xaa\x02\nHelloWorld\xca\x02\nHelloWorld\xe2\x02\x16HelloWorld\\GPBMetadata\xea\x02\nHelloWorldb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'helloworld.helloworld_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\017com.hello_worldB\017HelloworldProtoP\001ZHgithub.com/amirsalarsafaei/proto-error-handling/go/helloworld;helloworld\242\002\003HXX\252\002\nHelloWorld\312\002\nHelloWorld\342\002\026HelloWorld\\GPBMetadata\352\002\nHelloWorld'
  _globals['_USERSTATUS']._serialized_start=497
  _globals['_USERSTATUS']._serialized_end=587
  _globals['_CREATEUSERREQUEST']._serialized_start=44
  _globals['_CREATEUSERREQUEST']._serialized_end=113
  _globals['_CREATEUSERRESPONSE']._serialized_start=115
  _globals['_CREATEUSERRESPONSE']._serialized_end=209
  _globals['_CREATEUSERALTRESPONSE']._serialized_start=212
  _globals['_CREATEUSERALTRESPONSE']._serialized_end=347
  _globals['_USERDATA']._serialized_start=349
  _globals['_USERDATA']._serialized_end=433
  _globals['_ERRORDETAILS']._serialized_start=435
  _globals['_ERRORDETAILS']._serialized_end=495
  _globals['_USERSERVICE']._serialized_start=590
  _globals['_USERSERVICE']._serialized_end=771
# @@protoc_insertion_point(module_scope)
