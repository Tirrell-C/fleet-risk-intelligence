# Base Go builder image
FROM golang:1.23-alpine AS base

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /workspace

# Copy go mod files
COPY go.mod go.sum go.work go.work.sum ./
COPY services/api/go.mod services/api/go.sum* ./services/api/
COPY services/risk-engine/go.mod services/risk-engine/go.sum* ./services/risk-engine/
COPY services/telemetry-ingest/go.mod services/telemetry-ingest/go.sum* ./services/telemetry-ingest/
COPY services/websocket/go.mod services/websocket/go.sum* ./services/websocket/

# Download dependencies
RUN go mod download
RUN go work sync

# Copy source code
COPY pkg/ ./pkg/
COPY services/ ./services/