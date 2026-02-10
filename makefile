OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)

ifeq ($(ARCH),x86_64)
    ARCH_NAME := amd64
else ifeq ($(ARCH),arm64)
    ARCH_NAME := arm64
else
    ARCH_NAME := $(ARCH)
endif

ifeq ($(OS),darwin)
    BINARY_DIR := macos_$(ARCH_NAME)
else
    BINARY_DIR := linux_$(ARCH_NAME)
endif

.PHONY: prepare backend frontend build dev

backend:
	@echo "Starting server..."
	cd archivus-v2 && go build . 
	cd archivus-v2 && ./archivus-v2 server

frontend:
	@echo "Starting client..."
	cd archivus-client && npm run dev

dev:
	@echo "Starting server..."
	cd archivus-v2 && go run . server &
	cd archivus-client && npm run dev &
	wait
	
build:
	@echo "Building Project for $(OS)_$(ARCH_NAME)..."
	mkdir -p dist/bin/$(BINARY_DIR)
	cd archivus-v2 && go build -o ../dist/bin/$(BINARY_DIR)/archivus-v2 .
	cp archivus-v2/config.prod.yaml dist/bin/$(BINARY_DIR)/config.prod.yaml

	rm -rf dist/frontend
	mkdir -p dist/frontend
	cd archivus-client && npm run build
	cp -r archivus-client/.next dist/frontend/.next
	cp archivus-client/package.json dist/frontend/
	cp -r archivus-client/public dist/frontend/public


