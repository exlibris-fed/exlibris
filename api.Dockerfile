FROM golang:1.14-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates build-base sudo alpine-sdk

RUN mkdir /build
WORKDIR /build/
COPY . /build/

RUN go get
RUN go get github.com/githubnemo/CompileDaemon
CMD ["CompileDaemon", "-build=go build -o app", "-directory=/build", "-command=/build/app"]

EXPOSE 8080
