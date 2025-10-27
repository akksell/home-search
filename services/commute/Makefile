GOLANG=go

.PHONY: run-local
run-local:
	ENVIRONMENT="development" PORT=":8082" GOOGLE_PROJECT_ID="home-search-475404" COMMUTE_STORE_DATABASE="dev-commute-service-poi-store" ADDRESS_WRAPPER_SERVICE_HOST="127.0.0.1:8080" ROOMMATE_SERVICE_HOST="127.0.0.1:8081" $(GOLANG) run cmd/main.go