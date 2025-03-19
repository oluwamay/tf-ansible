FROM golang:1.24-alpine3.20 AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-server .

# Create a minimal production image
FROM alpine:3.20

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api-server .

# Expose port
EXPOSE 8080

# Set entry point
ENTRYPOINT ["./api-server"]