FROM golang:1.21 AS builder
WORKDIR /app

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ws-server .

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/ws-server .
EXPOSE 9002
CMD ["./ws-server"]