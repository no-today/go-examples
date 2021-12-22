FROM golang:latest

WORKDIR $GOPATH/src/cathub.me/go-web-examples
COPY . $GOPATH/src/cathub.me/go-web-examples

RUN export GOPROXY=https://goproxy.cn
RUN go build .

EXPOSE 8080
ENTRYPOINT ["./go-web-examples"]