FROM golang:alpine AS builder

WORKDIR /app
# Install gcc and musl-dev for sqlite3 cgo dependency
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/app

FROM alpine:latest
WORKDIR /app

# Required for sqlite
RUN apk add --no-cache sqlite-libs tzdata

COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
COPY --from=builder /app/.env.example ./.env

# Create database directory
RUN mkdir -p database

EXPOSE 8080
CMD ["./main"]
