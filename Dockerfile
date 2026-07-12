# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auction_app cmd/auction/main.go

# Stage 2: Final
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/auction_app .

# Expose the port the app runs on
EXPOSE 8080

# Run the application
ENTRYPOINT ["./auction_app"]