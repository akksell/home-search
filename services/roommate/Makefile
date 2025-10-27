GOLANG=go

.PHONY: run-local
run-local:
	ENVIRONMENT="development" PORT=":8081" GOOGLE_PROJECT_ID="home-search-475404" ROOMMATE_STORE_DATABASE="roommate-roommate-store-dev" ADDRESS_WRAPPER_SERVICE_HOST="127.0.0.1:8080" $(GOLANG) run cmd/main.go