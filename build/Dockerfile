FROM scratch
COPY --from=builder /app/bin/app /app
COPY . .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 3000
ENTRYPOINT ["/app"]
