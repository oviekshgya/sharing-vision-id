# Stage 1 - Build Golang
FROM golang:1.23-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o shagya-tech-payment


# Stage 2 - Final Runtime
FROM alpine:3.17 AS runner

WORKDIR /app

RUN apk add --no-cache libstdc++ libx11

# Copy binary golang
COPY --from=builder /app/shagya-tech-payment /app/shagya-tech-payment

# Copy file-file python / asset
COPY ./pkg/face/facev2.py /app/pkg/face/facev2.py
COPY public/storage/img /app/public/storage/img
COPY public/views /app/public/views
COPY .env /app/.env
COPY credentials.json /app/credentials.json

RUN chmod +x /app/shagya-tech-payment

CMD ["/app/shagya-tech-payment"]
