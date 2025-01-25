run:
	@echo "Running main.go..."
	docker exec go_app /bin/bash -c "cd cmd/wallet_core/ && go run main.go"

.PHONY: run