FROM golang:1.12 as builder
RUN git clone https://github.com/edenhill/librdkafka.git && cd librdkafka && ./configure --prefix /usr && make && make install
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
RUN mkdir /app
RUN mv /build/main /app/
WORKDIR /app
CMD ["./main"]