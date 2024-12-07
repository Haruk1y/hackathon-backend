FROM golang:1.23-alpine AS builder

WORKDIR /app

# ビルドに必要なパッケージのインストール
RUN apk add --no-cache git

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# 実行ステージ
FROM alpine:latest

WORKDIR /app

# タイムゾーンの設定
RUN apk add --no-cache tzdata
ENV TZ=Asia/Tokyo

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main .

# 必要なディレクトリの作成とパーミッションの設定
RUN chmod +x /app/main

# ポート設定
ENV PORT=8080
EXPOSE 8080

# アプリケーションの実行
CMD ["./main"]