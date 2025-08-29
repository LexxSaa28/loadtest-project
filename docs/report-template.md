# Performance Comparison Report

## Test Setup
- Date: {{date}}
- Total Requests: {{total_requests}}
- Concurrency: {{concurrency}}
- Payload Size: {{payload_size}} KB
- Duration: {{duration}}

## Results Summary

| Protocol | Requests/sec | Avg Response Time | Error Rate |
|----------|-------------|------------------|------------|
| REST     | {{rest_rps}} | {{rest_avg}}     | {{rest_err}}% |
| gRPC     | {{grpc_rps}} | {{grpc_avg}}     | {{grpc_err}}% |
| WebSocket| {{ws_rps}}   | {{ws_avg}}       | {{ws_err}}% |

## Detailed Analysis

### REST API
{{rest_analysis}}

### gRPC
{{grpc_analysis}}

### WebSocket
{{ws_analysis}}

## Conclusion
{{conclusion}}