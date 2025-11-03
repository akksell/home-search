GOLANG=go
CGO_ENABLED=0
GOOS=linux

.PHONY: build
build:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} ${GOLANG} build -o build/address_wrapper ./cmd/*.go

.PHONY: run-local
run-local:
	ENVIRONMENT=development GEOCODE_SECRET_NAME="projects/870628054987/secrets/address-wrapper-geocode-api-key-dev" GEOCODE_SECRET_VERSION="1" ${GOLANG} run cmd/main.go