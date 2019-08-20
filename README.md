# Alertmanager Webhook Forwarder

- [Alertmanager Webhook Forwarder](#alertmanager-webhook-forwarder)
  - [Build Matrix](#build-matrix)
  - [Docker Guide](#docker-guide)
    - [Build](#build)
    - [Run](#run)
    - [Lint](#lint)
  - [Go Guide](#go-guide)
    - [Init](#init)
    - [Run](#run-1)
    - [Build](#build-1)

---

## Build Matrix

| CI Job | Branch [`develop`](https://github.com/miquido/alertmanager-webhook-forwarder/tree/develop)                                                                                            | Branch [`master`](https://github.com/miquido/alertmanager-webhook-forwarder/tree/master)                                                                                            |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Circle | [![CircleCI](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/develop.svg?style=svg)](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/develop) | [![CircleCI](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/master.svg?style=svg)](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/master) |

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
