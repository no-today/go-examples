## Quick Start

HTTP Request see [requests.http](requests.http)

### IDE Run

```shell
docker-compose up -d mongo redis
```

### Docker Run

TODO: application.yml localhost -> container name

```shell
docker build -t go-web-examples .
docker-compose up -d
```

## Generate Google Protobuf

```shell
protoc --go_out=. ./protobuf/*.proto
```