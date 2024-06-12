FROM golang:1.22-alpine

# ディレクトリ指定
ENV ROOT=/go/src/app
WORKDIR ${ROOT}

# Git インストール
RUN apk update && apk add git
# ファイルコピー
COPY ./src/cmd/main.go ${ROOT}
COPY ./src/go.mod ${ROOT}

RUN go mod tidy
