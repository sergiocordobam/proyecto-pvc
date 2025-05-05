SERVICE_ACCOUNT_KEY_RELATIVE_PATH := "internal/infrastructure/gcp/service-account-key.json"

set-environment:
	@echo "Setting GOOGLE_APPLICATION_CREDENTIALS dynamically..."
	@export GOOGLE_APPLICATION_CREDENTIALS=$(shell pwd)/$(SERVICE_ACCOUNT_KEY_RELATIVE_PATH) && \
    	echo "GOOGLE_APPLICATION_CREDENTIALS set to: $$GOOGLE_APPLICATION_CREDENTIALS"
run:
	go run cmd/documents-api/main.go