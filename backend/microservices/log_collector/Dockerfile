# Use the official Go image as a base image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Create a minimal image for running the Go application
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the builder image
COPY --from=builder /app/main .
COPY .env .

# Expose the port your application will run on
EXPOSE 8080

# Run the Go application
CMD ["./main"]
