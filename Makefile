BIN      := e-gov-lawapi-v2
CMD_PKG  := ./cmd/$(BIN)
GOFLAGS  ?=

.PHONY: all build run test vet fmt tidy clean install help

all: vet test build

build: ## バイナリをリポジトリ直下に出力
	go build $(GOFLAGS) -o $(BIN) $(CMD_PKG)

run: ## ビルド済みバイナリを stdio で起動
	go run $(CMD_PKG)

test: ## 全パッケージの単体テスト
	go test ./...

vet: ## go vet
	go vet ./...

fmt: ## gofmt（書き換えなし、差分があれば失敗）
	@diff=$$(gofmt -l .); \
	if [ -n "$$diff" ]; then echo "gofmt needs to run on:"; echo "$$diff"; exit 1; fi

tidy: ## go.mod / go.sum を整理
	go mod tidy

install: ## $GOBIN にインストール
	go install $(CMD_PKG)

clean:
	rm -f $(BIN)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
