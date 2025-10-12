# ---------- Build stage ----------
    FROM golang:1.24-alpine AS builder

    WORKDIR /app
    
    # Download dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy source
    COPY . .
    
    # Build binary
    RUN go build -o main .
    
    # ---------- Run stage ----------
    FROM alpine:3.19
    
    WORKDIR /app
    
    # Copy compiled binary
    COPY --from=builder /app/main .
    
    # Expose reservation service port
    EXPOSE 8082
    
    CMD ["./main"]
    