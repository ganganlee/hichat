protoc --proto_path=protoc --micro_out=rpc/user --go_out=rpc/user user.proto
protoc --proto_path=protoc --micro_out=rpc/userFriends --go_out=rpc/userFriends userFriends.proto