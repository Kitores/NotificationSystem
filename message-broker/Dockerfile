FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest
COPY --from=builder /app/main /usr/local/bin/main
CMD ["main"]
EXPOSE 50051