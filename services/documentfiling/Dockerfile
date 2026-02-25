# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Build static binary (multi-arch compatible)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /documentfiling .

# --- Runtime stage ---
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /documentfiling .

EXPOSE 8083

CMD ["./documentfiling"]
