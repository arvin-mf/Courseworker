FROM golang:1.23.2-alpine as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /main ./cmd

FROM alpine:latest as run

COPY --from=build /main /main

WORKDIR /

CMD ["/main"]