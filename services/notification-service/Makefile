.PHONY: compile-proto lint format

compile-proto:
	python -m grpc_tools.protoc \
	  -I ../../common/proto \
	  --python_out=src/notification/proto \
	  --grpc_python_out=src/notification/proto \
	  ../../common/proto/notification.proto ../../common/proto/trajectory.proto

lint:
	ruff check . --fix

format:
	ruff format .
