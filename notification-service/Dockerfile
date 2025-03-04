FROM golang:alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine
COPY --from=builder /app/main /app/main
COPY /config/local.yaml /app/local.yaml
CMD ["/app/main"]
EXPOSE 50051
EXPOSE 587