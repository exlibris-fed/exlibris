FROM golang:1.12-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates build-base sudo alpine-sdk
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN GOOS=linux GOARCH=amd64 go build -o main .
RUN mkdir /app
RUN mv /build/main /app/
WORKDIR /app
CMD ["./main"]
