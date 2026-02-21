# Dockerfile for AetherBus-Tachyon

# --- Build Stage ---
# Use the official Go image as a builder.
# Specify the Go version.
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first.
# This leverages Docker layer caching.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# Copy the entire source code.
COPY . .

# Build the Go app, creating a static binary.
# CGO_ENABLED=0 is important for creating a static binary that can run in a minimal image.
# -o sets the output file name.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /aetherbus-tachyon ./cmd/tachyon

# --- Final Stage ---
# Use a minimal base image.
FROM alpine:latest

# Set the working directory.
WORKDIR /root/

# Copy the built binary from the builder stage.
COPY --from=builder /aetherbus-tachyon .

# Expose the ports the application uses.
# 5555 for the ROUTER socket, 5556 for the PUB socket.
EXPOSE 5555
EXPOSE 5556

# The command to run the application.
CMD ["./aetherbus-tachyon"]
