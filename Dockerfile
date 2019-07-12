ARG GOLANG_TAG="1.12.7-alpine3.10"
FROM golang:${GOLANG_TAG} as Go
RUN apk add --no-cache git g++ musl-dev
ENV GO111MODULE="on"
WORKDIR /usr/src/app
ENTRYPOINT ["go"]
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY main.go ./
COPY cmd ./cmd
COPY pkg ./pkg
ARG BINARY_NAME="alertmanager-webhook-forwarder"
ARG MAIN_GO="main.go"
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
    go build -a -tags netgo -ldflags "-w" \
    -o ${BINARY_NAME} ${MAIN_GO}

FROM scratch as Cli
WORKDIR /usr/src/app
COPY --from=Go /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ARG BINARY_NAME="alertmanager-webhook-forwarder"
COPY --from=Go /usr/src/app/${BINARY_NAME} ./app
ENTRYPOINT ["./app"]
CMD ["help"]
VOLUME ["/tmp"]

FROM Go as Lint
ARG GOLANGCI_LINT_VERSION="v1.17.1"
RUN wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}
ENTRYPOINT ["golangci-lint"]
CMD ["run"]
