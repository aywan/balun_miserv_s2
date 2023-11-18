REPOS:=auth-server chat-server chat-client
LOCAL_BIN:=$(CURDIR)/bin
LIBS:=$(CURDIR)/shared/lib

.PHONY: install-deps
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@v2.36.1

.PHONY: tidy
tidy:
	@set +e
	@for dir in $(REPOS); do echo "tidy for $$dir"; cd $$dir; make tidy; cd - > /dev/null; done
	@for dir in $(LIBS)/*/; do echo "tidy for $$dir"; cd $$dir; make tidy; cd - > /dev/null; done

.PHONY: generate
generate:
	@set +e
	@for dir in $(REPOS); do cd $$dir; make generate; cd -; done

.PHONY: lint
lint:
	@set +e
	@for dir in $(REPOS); do echo "lint for $$dir"; cd $$dir; make lint; @cd - > /dev/null; done
	@for dir in $(LIBS)/*/; do echo "lint for $$dir"; cd $$dir; make lint; @cd - > /dev/null; done

.PHONY: tests
tests:
	@set +e
	@for dir in $(REPOS); do echo "tests for $$dir"; cd $$dir; make tests; @cd - > /dev/null; done
	@for dir in $(LIBS)/*/; do echo "tests for $$dir"; cd $$dir; make tests; @cd - > /dev/null; done

.PHONY: dc-up
dc-up:
	docker compose up -d

.PHONY: shell-auth-server
shell-auth-server:
	docker compose exec auth_service bash

.PHONY: shell-chat-server
shell-chat-server:
	docker compose exec chat_service bash
