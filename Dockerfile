FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o /api-monitor

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /api-monitor .
EXPOSE 8080
CMD ["./api-monitor"]
