.PHONY: help install fmt vet build clean mod

help: ## 显示命令帮助
	@awk 'BEGIN {FS = ":.*##"; printf "goweb-core 可用命令：\n"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## 安装 Go 依赖
	@echo "📦 安装 Go 依赖"
	@go mod download

fmt: ## 格式化 Go 代码
	@echo "🧹 gofmt"
	@gofmt -w $$(find . -name '*.go' -print)

vet: ## 运行 go vet
	@echo "🔎 go vet"
	@go vet ./...

build: ## 构建 core 包
	@echo "🔨 构建 goweb-core"
	@go build ./...

clean: ## 清理构建产物
	@echo "🧽 清理构建产物"
	@rm -rf bin dist tmp

mod: ## 整理 Go module
	@echo "🧩 go mod tidy"
	@go mod tidy
