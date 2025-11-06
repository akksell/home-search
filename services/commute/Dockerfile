FROM gcr.io/distroless/base-debian11 as release-stage

ENV ENVIRONMENT=production

WORKDIR /

COPY ./build/commute commute

USER nonroot:nonroot

ENTRYPOINT ["./commute"]