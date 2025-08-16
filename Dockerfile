# Use official Golang image to build the app
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY *.go ./

# Build the Go app (statically linked, so it runs in alpine)
RUN go build -o main .

# Second stage: Use lightweight Alpine Linux
FROM alpine:latest

# Install CA certificates (for HTTPS/database connections)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Command to run the binary
CMD ["./main"]