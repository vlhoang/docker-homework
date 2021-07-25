FROM golang:alpine as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o bookapi

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bookapi .
COPY config/app/.env .
COPY config/app/app.conf .

CMD ./bookapi