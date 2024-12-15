# バイナリのビルド用のステージ
FROM golang:1.21-bullseye as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

# 実行用の軽量イメージ
FROM debian:bullseye-slim

# Cloud SQL Auth Proxyのインストール
RUN apt-get update && apt-get install -y ca-certificates wget
RUN wget https://dl.google.com/cloudsql/cloud-sql-proxy.linux.amd64 -O /cloud-sql-proxy
RUN chmod +x /cloud-sql-proxy

WORKDIR /app
COPY --from=builder /app/main .
COPY firebase-credentials.json .
COPY .env .

EXPOSE 8080
CMD ["./main"]