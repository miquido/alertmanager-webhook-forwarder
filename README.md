# Alertmanager Webhook Forwarder

- [Alertmanager Webhook Forwarder](#Alertmanager-Webhook-Forwarder)
  - [Docker Guide](#Docker-Guide)
    - [Build](#Build)
    - [Run](#Run)
    - [Lint](#Lint)
  - [Go Guide](#Go-Guide)
    - [Init](#Init)
    - [Run](#Run-1)
    - [Build](#Build-1)

## Docker Guide

### Build

```sh
docker-compose build --pull
```

### Run

```sh
docker-compose run --rm cli help
```

### Lint

```sh
docker-compose run --rm go
```

## Go Guide

### Init

```sh
export GO111MODULE=on
go mod init
go mod download
go mod verify
```

### Run

```sh
go run main.go help
```

### Build

```sh
go build -o alertmanager-webhook-forwarder main.go
./alertmanager-webhook-forwarder help
```
