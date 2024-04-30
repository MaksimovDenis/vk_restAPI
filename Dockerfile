FROM golang:1.21.6-alpine AS builder

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apk update && apk add --no-cache postgresql-client

RUN go mod download
RUN go build -o vk_restapi ./cmd/main.go


CMD ["./vk_restapi"]
