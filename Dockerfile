FROM golang:alpine as builder
RUN apk -U add ca-certificates
RUN apk update && apk upgrade && apk add git bash build-base sudo alpine-sdk
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .
RUN mkdir /app
RUN mv /build/main /app/
WORKDIR /app
CMD ["./main"]