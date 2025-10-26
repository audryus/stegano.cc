FROM golang:1.25 AS builder

FROM scratch
WORKDIR /app
COPY . .
COPY ./bin/app ./bin/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 3000
ENTRYPOINT ["/app/bin/app"]
