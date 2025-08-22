# ----------- Build Stage -----------
FROM golang:1.24.1-alpine AS builder

# Enable Go modules and ensure no CGO (for a static binary)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Cache dependencies
COPY go.mod ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN go build -o server ./main.go

# ----------- Final Stage -----------
FROM alpine:3.19

WORKDIR /app

# Copy binary
COPY --from=builder /app/server .

# Copy static files
COPY static ./static

# Copy templates
COPY templates ./templates

# Expose port
EXPOSE 8080

# Run app
CMD ["./server"]
