FROM golang:1.25.6-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags "-s -w" -a -o metadata-service ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /build/metadata-service .
COPY assets/fallback /app/assets/fallback

COPY .env .

EXPOSE 8080

ENTRYPOINT ["./metadata-service"]
