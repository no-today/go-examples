## Quick Start

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

## CURL

### Register

```shell
curl --location --request POST 'localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "no-today",
    "password": "changeme",
    "email": "web_cheng@163.com"
}'
```

### Activation

```shell
curl --location --request GET 'localhost:8080/api/v1/activation/:activation_code'
```

### Authenticate

```shell
curl --location --request POST 'localhost:8080/api/v1/authenticate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "principal": "no-today",
    "credentials": "changeme"
}'
```

### ResendActivateEmail

```shell
curl --location --request POST 'localhost:8080/api/v1/resendActivateEmail' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "web_cheng@163.com"
}'
```