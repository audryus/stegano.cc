FROM golang:1.25 AS builder

FROM scratch
COPY --from=builder /app/bin/app /app
COPY . .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 3000
ENTRYPOINT ["/app"]
