# Use Go's official image to create a build artifact
FROM golang:1.23-alpine AS builder

# Install timezone data
RUN apk add --no-cache tzdata

ENV TZ=America/New_York

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary
RUN go build -o main .

# Use a minimal base image for running
FROM alpine:latest

# Install timezone data for runtime
RUN apk add --no-cache tzdata

WORKDIR /app

# Set the timezone environment variable
ENV TZ=America/New_York

# Copy binary from builder
COPY --from=builder /app/main .

# Expose the port your app runs on
EXPOSE 8080

# Command to run your binary
CMD ["./main"]
