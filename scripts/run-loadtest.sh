#!/bin/bash

# Build all services
docker-compose build

# Start infrastructure
docker-compose up -d statsd grafana rest-server grpc-server ws-server

# Wait for services to start
sleep 10

# Run load tests for each protocol
protocols=("rest" "grpc" "ws")

for protocol in "${protocols[@]}"; do
  echo "Running load test for $protocol protocol..."
  
  docker-compose run --rm load-test ./loadtest \
    -p $protocol \
    -r 10000 \
    -c 100 \
    -s 1 \
    -d 30s
  
  # Wait between tests
  sleep 10
done

# Generate report
echo "Generating performance report..."
docker-compose run --rm load-test python generate_report.py

# Stop services
docker-compose down