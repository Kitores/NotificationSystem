FROM golang:1.22.1-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o NotificationSystem ./cmd/NotificationSystem/main.go
COPY ./config/local.yaml ./config/local.yaml

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/NotificationSystem .

ENV BOT_TOKEN=
COPY ./config/local.yaml ./config/local.yaml
EXPOSE 50051

CMD ["./NotificationSystem"]