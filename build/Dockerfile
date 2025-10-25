FROM golang:1.25 AS builder

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./ ./
EXPOSE 3000
ENTRYPOINT ["/bin/app"]
