.PHONY: prepare backend frontend build dev

prepare:
# 	mkdir -p bin
# 	mkdir -p bin/server

backend:
	@echo "Starting server..."
	cd archivus-v2 && go build . 
	cd archivus-v2 && ./archivus-v2 server
# 	integrate flag

frontend:
	@echo "Starting client..."
	cd archivus-client && npm run dev

dev:
	@echo "Starting server..."
	cd archivus-v2 && go run . server &
	cd archivus-client && npm run dev &
	wait
	
build:
	@echo "Building Project..."
	

	cd archivus-v2 && go build -o ../dist/bin/linux_amd64/ .
	cp archivus-v2/config.prod.yaml dist/bin/linux_amd64/config.prod.yaml

	rm -rf dist/frontend
	mkdir -p dist/frontend
	cd archivus-client && npm run build
	cp -r archivus-client/.next dist/frontend/.next && cp archivus-client/package.json dist/frontend/ && cp -r archivus-client/public dist/frontend/public


