GOLANG=go
CGO_ENABLED=0
GOOS=linux

.PHONY: build
build:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} ${GOLANG} build -o build/commute ./cmd/*.go

.PHONY: run-local
run-local:
	ENVIRONMENT="development" PORT=":8082" GOOGLE_PROJECT_ID="home-search-475404" COMMUTE_STORE_DATABASE="dev-commute-service-poi-store" ADDRESS_WRAPPER_SERVICE_HOST="hs-address-wrapper-dev-870628054987.us-south1.run.app:443" ROOMMATE_SERVICE_HOST="127.0.0.1:8081" $(GOLANG) run cmd/main.go