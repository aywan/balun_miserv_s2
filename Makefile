REPOS:=auth-server chat-server chat-client
LOCAL_BIN:=$(CURDIR)/bin

.PHONY:install-deps
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

.PHONY:tidy
lint:
	@set +e
	@for dir in $(REPOS); do cd $$dir; make tidy; cd -; done

.PHONY:generate
lint:
	@set +e
	@for dir in $(REPOS); do cd $$dir; make generate; cd -; done

.PHONY:lint
lint:
	@set +e
	@for dir in $(REPOS); do cd $$dir; make lint; cd -; done

.PHONY:tests
tests:
	@set +e
	@for dir in $(REPOS); do cd $$dir; make tests; cd -; done
