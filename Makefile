.PHONY: generate clean test build

# Generate protobuf files for all services
generate:
	buf generate

# Generate protobuf files for specific service
generate-commute:
	buf generate --path services/commute

generate-address-wrapper:
	buf generate --path services/address_wrapper

# Clean generated files
clean:
	find . -name "*.pb.go" -delete
	find . -name "*_grpc.pb.go" -delete
	find . -name "*.pb.gw.go" -delete
	find . -name "gen" -type d -exec rm -rf {} +

# Run tests
test:
	go test ./...

# Build all services
build:
	go build ./...

# Lint protobuf files
lint:
	buf lint

# Check for breaking changes
breaking:
	buf breaking --against '.git#branch=main'
