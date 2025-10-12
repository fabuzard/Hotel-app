# ---------- Build stage ----------
    FROM golang:1.24-alpine AS builder

    WORKDIR /app
    
    # Download dependencies early for caching
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the rest of the source code
    COPY . .
    
    # Build the binary
    RUN go build -o main .
    
    # ---------- Run stage ----------
    FROM alpine:3.19
    
    WORKDIR /app
    
    # Copy binary from builder
    COPY --from=builder /app/main .
    
    # Copy env if needed (optional for debugging)
    # COPY .env .env
    
    # Expose service port
    EXPOSE 8083
    
    # Run the application
    CMD ["./main"]
    