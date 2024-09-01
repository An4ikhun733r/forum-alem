# Stage 1: Build
FROM golang:1.20.1-alpine3.16 as builder

# Install necessary packages for building
RUN apk add --no-cache build-base

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go application
RUN go build -o forum ./cmd/web/

# Stage 2: Final Image
FROM alpine:3.16

# Install necessary runtime packages
RUN apk add --no-cache ca-certificates

# Set the working directory in the container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/forum .

# Copy the static and template files
COPY ui /app/ui

# Copy the database file
COPY data/database.db /app/data/database.db

# Set the command to run the application
CMD ["./forum"]
