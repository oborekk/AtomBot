FROM alpine:latest
WORKDIR /app
COPY AtomBot .
CMD ["./AtomBot"]
