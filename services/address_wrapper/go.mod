module homesearch.axel.to/address_wrapper

go 1.24.0

toolchain go1.24.9

tool (
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)

require google.golang.org/grpc v1.76.0

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251014184007-4626949a642f // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
