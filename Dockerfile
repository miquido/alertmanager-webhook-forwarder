ARG ALPINE_VERSION="3.11"
ARG GOLANG_VERSION="1.13.6"
ARG GOLANGCI_LINT_VERSION="1.22.2"

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as Go
RUN apk add --no-cache git g++ musl-dev
ENV GO111MODULE="on"
WORKDIR /usr/src/app
ENTRYPOINT ["go"]
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY cmd ./cmd
COPY pkg ./pkg

FROM Go as GoBuilder
ARG MAIN_GO="cmd/alertmanager-webhook-forwarder/main.go"
ARG BINARY_NAME="dist/alertmanager-webhook-forwarder"
RUN CGO_ENABLED="0" GOARCH="amd64" GOOS="linux" \
    go build -a -tags netgo -ldflags "-w" \
    -o ${BINARY_NAME} ${MAIN_GO}

FROM alpine:${ALPINE_VERSION} as RunnerBase
RUN apk add --no-cache ca-certificates
RUN addgroup -g 1000 -S runner && \
    adduser -u 1000 -S app -G runner

FROM scratch as Runner
WORKDIR /usr/src/app
COPY --from=RunnerBase /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=RunnerBase /etc/passwd /etc/passwd
ARG BINARY_NAME="dist/alertmanager-webhook-forwarder"
COPY --from=GoBuilder /usr/src/app/${BINARY_NAME} ./alertmanager-webhook-forwarder
USER app
ENTRYPOINT ["/usr/src/app/alertmanager-webhook-forwarder"]
CMD ["help"]
VOLUME ["/tmp"]

FROM golangci/golangci-lint:v${GOLANGCI_LINT_VERSION}-alpine as GolangCI-Lint

FROM Go as Lint
COPY --from=GolangCI-Lint /usr/bin/golangci-lint /usr/bin/golangci-lint
COPY .golangci.yml ./
ENTRYPOINT ["golangci-lint"]
CMD ["run"]
