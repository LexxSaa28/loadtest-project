FROM golang:1.21 AS builder
WORKDIR /app

COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o grpc-server .

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/grpc-server .
EXPOSE 9001
CMD ["./grpc-server"]