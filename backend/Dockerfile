# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/akritas \
    ./cmd

FROM alpine:3.23

RUN apk add --no-cache ca-certificates \
    && addgroup -S akritas \
    && adduser -S -G akritas akritas

COPY --from=builder /out/akritas /usr/local/bin/akritas

USER akritas
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/akritas"]
