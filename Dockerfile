# Stage 1: Build dependencies and tools
FROM golang:1.23.0 AS builder

WORKDIR /usr/src/app

# Install Air globally (leverages Docker's caching)
RUN go install github.com/air-verse/air@v1.61.5

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies (caches go mod dependencies separately)
RUN go mod download

# Copy the application code (invalidates cache only when code changes)
COPY . .

# Stage 2: Runtime (lightweight and uses prebuilt Air)
FROM golang:1.23.0

WORKDIR /usr/src/app

# Copy application files and Air binary from builder
COPY --from=builder /usr/src/app /usr/src/app
COPY --from=builder /go/bin/air /usr/local/bin/air

# Expose the application port
EXPOSE ${APP_PORT}

# Set up the command for live reloading
CMD ["air", "--config", "/usr/src/app/.air.toml"]
