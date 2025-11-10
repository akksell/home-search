FROM envoyproxy/envoy:distroless-v1.36.2 AS envoy
COPY --chmod=644 ./service-descriptor.binpb /etc/envoy/service-descriptor.binpb
COPY --chmod=644 ./configuration.yaml /etc/envoy/envoy.yaml