# Alertmanager Webhook Forwarder

- [Alertmanager Webhook Forwarder](#Alertmanager-Webhook-Forwarder)
  - [Build Matrix](#Build-Matrix)
  - [Docker Guide](#Docker-Guide)
    - [Build](#Build)
    - [Run](#Run)
    - [Lint](#Lint)
  - [Go Guide](#Go-Guide)
    - [Init](#Init)
    - [Run](#Run-1)
    - [Build](#Build-1)

---

## Build Matrix

| CI       | Branch                                                                                | Status                                                                                                                                                                                  |
| -------- | ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| CircleCI | [`develop`](https://github.com/miquido/alertmanager-webhook-forwarder/tree/develop)   | [![CircleCI](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/develop.svg?style=svg)](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/develop)   |
| CircleCI | [`master`](https://github.com/miquido/alertmanager-webhook-forwarder/tree/master)     | [![CircleCI](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/master.svg?style=svg)](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/master)     |
| CircleCI | [`gh-pages`](https://github.com/miquido/alertmanager-webhook-forwarder/tree/gh-pages) | [![CircleCI](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/gh-pages.svg?style=svg)](https://circleci.com/gh/miquido/alertmanager-webhook-forwarder/tree/gh-pages) |

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
