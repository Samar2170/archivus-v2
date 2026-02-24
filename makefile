
VERSION=0.1.0-beta.1
PROJECT_NAME=archivus-v2
DIST_DIR=dist
BIN_DIR=$(DIST_DIR)/bin
PKG_DIR=$(DIST_DIR)/packages

build: build-backend build-frontend

build-backend:
	@echo "Building Backend..."
	mkdir -p $(BIN_DIR)/linux_amd64
	mkdir -p $(BIN_DIR)/darwin_amd64
	mkdir -p $(BIN_DIR)/darwin_arm64
	
	cd archivus-v2 && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ../$(BIN_DIR)/linux_amd64/$(PROJECT_NAME) .
# 	cd archivus-v2 && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ../$(BIN_DIR)/darwin_amd64/$(PROJECT_NAME) .
# 	cd archivus-v2 && CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o ../$(BIN_DIR)/darwin_arm64/$(PROJECT_NAME) .
	
	cp archivus-v2/config.prod.yaml $(BIN_DIR)/linux_amd64/ || cp archivus-v2/config.yaml $(BIN_DIR)/linux_amd64/config.prod.yaml
# 	cp archivus-v2/config.prod.yaml $(BIN_DIR)/darwin_amd64/ || cp archivus-v2/config.yaml $(BIN_DIR)/darwin_amd64/config.prod.yaml
# 	cp archivus-v2/config.prod.yaml $(BIN_DIR)/darwin_arm64/ || cp archivus-v2/config.yaml $(BIN_DIR)/darwin_arm64/config.prod.yaml

build-frontend:
	@echo "Building Frontend..."
	rm -rf $(DIST_DIR)/frontend
	mkdir -p $(DIST_DIR)/frontend
	cd archivus-client && npm run build
	cp -r archivus-client/.next $(DIST_DIR)/frontend/
	cp archivus-client/package.json $(DIST_DIR)/frontend/
	cp -r archivus-client/public $(DIST_DIR)/frontend/

package: build
	@echo "Packaging Release..."
	mkdir -p $(PKG_DIR)
	# Linux amd64
	tar -czf $(PKG_DIR)/archivus-v2-$(VERSION)-linux-amd64.tar.gz -C $(BIN_DIR)/linux_amd64 . -C ../../frontend .
	# Darwin amd64
# 	tar -czf $(PKG_DIR)/archivus-v2-$(VERSION)-darwin-amd64.tar.gz -C $(BIN_DIR)/darwin_amd64 . -C ../../frontend .
	# Darwin arm64
# 	tar -czf $(PKG_DIR)/archivus-v2-$(VERSION)-darwin-arm64.tar.gz -C $(BIN_DIR)/darwin_arm64 . -C ../../frontend .

clean:
	rm -rf $(DIST_DIR)
