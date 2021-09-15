protoc --proto_path=protoc --micro_out=rpc/user --go_out=rpc/user user.proto
protoc --proto_path=protoc --micro_out=rpc/userFriends --go_out=rpc/userFriends userFriends.proto
protoc --proto_path=protoc --micro_out=rpc/UserGroups --go_out=rpc/UserGroups userGroups.proto
protoc --proto_path=protoc --micro_out=rpc/UserGroupMembers --go_out=rpc/UserGroupMembers userGroupMembers.proto