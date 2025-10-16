# Operator Rate Limiter

A Kubernetes operator that demonstrates custom rate limiting for controller reconciliation.

## Features

- Custom rate limiter with exponential backoff
- Max retry handling with callbacks
- Automatic status updates when max retries reached
- Modern controller-runtime APIs

## Usage

```bash
# Build and run
go build -o operator-rate-limiter main.go
./operator-rate-limiter
```

