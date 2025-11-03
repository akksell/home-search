GOLANG=go

.PHONY: build
build:
	${GOLANG} build -o build/address_wrapper ./cmd/*.go

.PHONY: run-local
run-local:
	ENVIRONMENT=development GEOCODE_SECRET_NAME="projects/870628054987/secrets/address-wrapper-geocode-api-key-dev" GEOCODE_SECRET_VERSION="1" ${GOLANG} run cmd/main.go