# Build Stage
FROM golang:alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source tree
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /investwise ./cmd/api

# Run Stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /investwise .

# Expose HTTP service port
EXPOSE 8080

# Run entry point
CMD ["./investwise"]
