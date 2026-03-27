.PHONY: build build-local clean install

# 默认目标：构建所有平台的二进制文件
build:
	@echo "Building binaries..."
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/yx-code-darwin-amd64 ./cmd/yx-code
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/yx-code-darwin-arm64 ./cmd/yx-code
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/yx-code-linux-amd64 ./cmd/yx-code
	@echo "Build complete!"

# 仅构建当前平台的二进制文件
build-local:
	@echo "Building for current platform..."
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/yx-code ./cmd/yx-code
	@echo "Build complete!"

# 安装到系统 PATH
install: build-local
	sudo cp bin/yx-code /usr/local/bin/yx-code
	@echo "Installed yx-code to /usr/local/bin/"

# 清理构建产物
clean:
	rm -rf bin/
	rm -f yx-code

# 发布到 npm（需要先运行 make build）
publish: build
	npm publish