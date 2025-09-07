# Build Go binary
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Production image
FROM ubuntu:latest
WORKDIR /root
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
EXPOSE 7540
CMD ["./main"]