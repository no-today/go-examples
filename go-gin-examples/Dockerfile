FROM golang:alpine as builder

WORKDIR /go/src/go-gin-examples
COPY . /go/src/go-gin-examples

RUN export GOPROXY=https://goproxy.cn
RUN go build .

FROM alpine:latest as prod

WORKDIR /go/src/go-gin-examples
COPY --from=0 /go/src/go-gin-examples .
COPY --from=0 /go/src/go-gin-examples/conf/* ./conf

EXPOSE 8080
ENTRYPOINT ["./go-gin-examples"]