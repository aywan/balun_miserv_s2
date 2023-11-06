LOCAL_BIN?=$(CURDIR)/../../../bin
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
