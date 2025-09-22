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
WORKDIR /app/services/risk-engine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/services/risk-engine/main .

# No port exposure needed for background service

CMD ["./main"]