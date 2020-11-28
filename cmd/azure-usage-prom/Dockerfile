ARG nonroot_image=gcr.io/distroless/static:nonroot

FROM ${nonroot_image}
WORKDIR /
COPY azure-usage-prom azure-usage-prom
USER nonroot:nonroot

ENTRYPOINT ["/azure-usage-prom"]