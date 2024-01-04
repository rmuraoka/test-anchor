FROM golang:latest

WORKDIR /app
COPY . .

# gitをインストール
RUN apt-get update && apt-get install -y git

# CompileDaemonをインストール
RUN go install github.com/githubnemo/CompileDaemon@latest

# CompileDaemonを実行
ENTRYPOINT ["/go/bin/CompileDaemon", "-log-prefix=false", "-build=go build -o /app/main", "-command=/app/main"]
