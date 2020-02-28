FROM golang:1.12-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh ca-certificates build-base sudo alpine-sdk
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install
RUN rm -rf librdkafka

RUN mkdir /build
WORKDIR /build/
COPY . /build/

RUN go get
RUN go get github.com/githubnemo/CompileDaemon
CMD ["CompileDaemon", "-build=go build -o app", "-directory=/build", "-command=/build/app"]

EXPOSE 8080
