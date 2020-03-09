FROM golang:1.12-alpine as gobuilder

RUN apk update && apk upgrade && apk add --no-cache bash git openssh ca-certificates build-base sudo alpine-sdk
WORKDIR /build
ADD . /build/
RUN GOOS=linux GOARCH=amd64 go build -o main .

FROM node:lts-alpine as nodebuilder

WORKDIR /build
COPY . .
RUN npm i
RUN npm run build

FROM alpine as runtime
RUN apk update && apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=gobuilder /build/main /app/main
COPY --from=nodebuilder /build/dist /app/dist
CMD ["./main"]
