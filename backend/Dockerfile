# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod ./
# Note: we don't have a go.sum yet, so we will generate it during build or let go mod tidy run
RUN go mod tidy

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o tutora-backend main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/tutora-backend .
COPY --from=builder /app/.env ./

EXPOSE 8080

CMD ["./tutora-backend"]
