run:
	echo "run"

build_proto:
	protoc --go_out=plugins=grpc:proto ./proto/chatMessage.proto
