FROM golang:1.24.6-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags "-s -w" -a -o tle-service ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /build/tle-service .
COPY .env .

EXPOSE 8080

ENTRYPOINT ["./tle-service"]
