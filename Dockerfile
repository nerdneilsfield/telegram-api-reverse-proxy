# Stage 1: Build the Go binary
FROM golang:1.21.2-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the source from the current directory to the working directory inside the container 
COPY ./telegram-api.go telegram-api.go

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o telegram-api ./telegram-api.go

# Stage 2: Run the binary in a minimal image
FROM alpine:latest

# Set the current working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/telegram-api .

# Expose the application on port 8080
EXPOSE 8080

# Command to run the binary
CMD ["bash"]
