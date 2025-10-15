# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /app/s4s-backend ./cmd/

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates netcat-openbsd libc6-compat

# Create app directory and set permissions
RUN mkdir -p /app && chmod -R 777 /app

# Copy the binary from builder
COPY --from=builder /app/s4s-backend /app/s4s-backend

# Make the binary executable
RUN chmod +x /app/s4s-backend

# Copy and set up the wait script
COPY wait-for-postgres.sh /app/wait-for-postgres.sh
RUN chmod +x /app/wait-for-postgres.sh

# Verify the binary
RUN ls -la /app/ && \
    file /app/s4s-backend && \
    ldd /app/s4s-backend 2>&1 || true

# Command to run the application
CMD ["/app/s4s-backend"]
