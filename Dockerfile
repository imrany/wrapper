# Multi-stage build for smaller final image
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wrapper .

# Final stage - minimal image
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/wrapper .

# Expose port
EXPOSE 8000

# Allow overriding flags via CMD or ENTRYPOINT
ENTRYPOINT ["./wrapper"]
