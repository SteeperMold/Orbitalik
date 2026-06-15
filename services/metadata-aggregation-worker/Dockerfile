FROM golang:1.26.2-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags "-s -w" -a -o metadata-aggregation-worker ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/metadata-aggregation-worker .

COPY .env .

RUN chmod +x /app/*

EXPOSE 8080

ENTRYPOINT ["./metadata-aggregation-worker"]
