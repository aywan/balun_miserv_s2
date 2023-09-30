LOCAL_BIN:=$(CURDIR)/../bin
GO_EXE?=go

.Phony:tidy
tidy:
	$(GO_EXE) mod tidy

.Phony:lint
lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config ../.golangci.pipeline.yaml

.Phony:tests
tests:
	$(GO_EXE) test ./...

.Phony:generate
generate:
	true
