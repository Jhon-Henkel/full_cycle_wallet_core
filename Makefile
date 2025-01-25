wc-run:
	@echo "Running wallet_core main.go..."
	docker compose start && docker exec go_app /bin/bash -c "cd cmd/wallet_core/ && go run main.go"

wc-bash:
	@echo "Starting wallet_core bash..."
	docker compose start && docker exec -it go_app bash

.PHONY: wc-run wc-bash