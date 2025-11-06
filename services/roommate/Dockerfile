FROM gcr.io/distroless/base-debian11 as release-stage

ENV ENVIRONMENT=production

WORKDIR /

COPY ./build/roommate roommate

USER nonroot:nonroot

ENTRYPOINT ["./roommate"]