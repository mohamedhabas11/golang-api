# Stage 1: Build
FROM golang:1.23.0 AS builder

WORKDIR /usr/src/app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire application code
COPY . .

# Build the binary
RUN go build -o /usr/local/bin/main ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.18

# Set environment variables
ENV GIN_MODE=release
ENV APP_PORT=3000

# Copy the binary from the builder stage
COPY --from=builder /usr/src/app/main /main

# Expose the configurable application port
EXPOSE ${APP_PORT}

# Run the application
CMD ["/main"]
