FROM golang:1.25 AS builder

FROM scratch
COPY ./bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./config/config.yaml ./config/config.yaml
EXPOSE 3000
ENTRYPOINT ["/app"]
