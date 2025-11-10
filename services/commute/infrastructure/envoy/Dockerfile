FROM envoyproxy/envoy:distroless-v1.36.2 AS envoy
COPY ./service-descriptor.binpb /etc/envoy/service-descriptor.binpb
COPY ./configuration.yaml /etc/envoy/envoy.yaml
RUN chmod go+r /etc/envoy/envoy.yaml