FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
COPY go.work go.work.sum ./
COPY pkg/ pkg/
# Copy all services (needed for go.work)
COPY services/ services/
COPY infrastructure/ infrastructure/

# Download dependencies
RUN go mod download

# Build the application
WORKDIR /app/services/telemetry-ingest
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/services/telemetry-ingest/main .

# Expose port
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

CMD ["./main"]