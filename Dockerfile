# Logic MCP Server
# Use a recent Go version that matches go.mod requirement
FROM golang:1.23-alpine AS builder

# Install SWI-Prolog dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o logic-mcp ./cmd/server

# Final stage
FROM swipl:latest

# Prolog is already installed. We can still install ca-certificates if needed.
RUN apt-get update && apt-get install -y ca-certificates

# Create non-root user and group (Debian-compatible syntax)
RUN addgroup --system --gid 1001 logicmcp && \
    adduser --system --uid 1001 --ingroup logicmcp --disabled-password logicmcp

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/logic-mcp .
COPY --from=builder /app/examples ./examples

# Set ownership for the non-root user
RUN chown -R logicmcp:logicmcp /app

# Switch to non-root user (using Debian-based user)

USER logicmcp

# Expose port for HTTP mode
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default to STDIO mode, but can be overridden
ENTRYPOINT ["./logic-mcp"]
CMD ["-mode", "stdio"]