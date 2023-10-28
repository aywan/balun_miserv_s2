FROM alpine:3.18.4

ARG GOOSE_VERSION="v3.15.1"

RUN apk add --no-cache bash
ADD https://github.com/pressly/goose/releases/download/${GOOSE_VERSION}/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose && mkdir -p /migrations

COPY entrypoint.sh /entrypoint.sh

WORKDIR /migrations
ENTRYPOINT /entrypoint.sh
