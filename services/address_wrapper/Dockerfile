FROM gcr.io/distroless/base-debian11 as release-stage

ENV ENVIRONMENT=production
ENV GEOCODE_SECRET_NAME="projects/870628054987/secrets/address-wrapper-geocode-api-key-dev"
ENV GEOCODE_SECRET_VERSION="1"

WORKDIR /

COPY ./build/address_wrapper address_wrapper

USER nonroot:nonroot

ENTRYPOINT ["./address_wrapper"]