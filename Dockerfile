# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app

# Copy the go.mod and go.sum files to utilize Go modules.
COPY go.mod go.sum ./

# Download the Go modules.
RUN go mod download

# Copy the application's source code into the container.
COPY . .

# Set environment variables for cross-compilation.
# These lines can be omitted if compiling for the same architecture as the host machine.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Compile the application. Adjust the source path if necessary.
# The entry point for the Go application is located at cmd/main.go.
RUN go build -o /car-pooling-challenge cmd/main.go

# Final stage: Create a lightweight Alpine image for running the compiled binary.
FROM alpine:3.19

# Install necessary dependencies for the final image.
RUN apk --no-cache add ca-certificates libc6-compat

# Copy the compiled binary from the builder stage to the final image.
COPY --from=builder /car-pooling-challenge /car-pooling-challenge

# Application listens to port:
EXPOSE 9091

# Configure the container to run the application.
ENTRYPOINT ["/car-pooling-challenge"]
