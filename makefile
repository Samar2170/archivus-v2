.PHONY: prepare server backend dev



prepare:
# 	mkdir -p bin
# 	mkdir -p bin/server

server:
	@echo "Starting server..."
	cd archivus-v2 && go build . 
	cd archivus-v2 && ./archivus-v2
	cd archivus-client && npm run build
	cd archivus-client && npm run start

backend:
	@echo "Starting server..."
	cd archivus-v2 && go build . 
	cd archivus-v2 && ./archivus-v2

dev:
	@echo "Starting server..."
	cd archivus-v2 && go run .
	cd archivus-client && npm run dev
	

