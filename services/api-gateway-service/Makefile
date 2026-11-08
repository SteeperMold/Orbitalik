PROTO_DIR := ../../common/proto

.PHONY: compile-proto
compile-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=. \
		--go-grpc_out=. \
		$(PROTO_DIR)/metadata.proto \
		$(PROTO_DIR)/tle.proto \
		$(PROTO_DIR)/trajectory.proto