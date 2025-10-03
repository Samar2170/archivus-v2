.PHONY: prepare server backend dev



prepare:
# 	mkdir -p bin
# 	mkdir -p bin/server

server:
	@echo "Starting server..."
	cd archivus-v2 && go build . 
	cd archivus-v2 && ./archivus-v2 server &
	cd archivus-client && npm run build &
	cd archivus-client && npm run start &
	wait

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
	

