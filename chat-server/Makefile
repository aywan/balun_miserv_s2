LOCAL_BIN:=$(CURDIR)/../bin
GO_EXE?=go

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config ../.golangci.pipeline.yaml

tests:
	$(GO_EXE) test ./...
