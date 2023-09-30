LOCAL_BIN:=$(CURDIR)/../bin
GO_EXE?=go

.PHONY: tidy
tidy:
	$(GO_EXE) mod tidy

.PHONY: lint
lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config ../.golangci.pipeline.yaml

.PHONY: tests
tests:
	$(GO_EXE) test ./...

.PHONY: generate
generate:
	make generate-api-user-v1

.PHONY: generate-api-user-v1
generate-api-user-v1:
	mkdir -p pkg/grpc/v1/user_v1
	protoc --proto_path api/v1/user \
		--go_out=pkg/grpc/v1/user_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go-grpc_out=pkg/grpc/v1/user_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		api/v1/user/user.proto
