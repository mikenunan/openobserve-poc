# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Build static binary (multi-arch compatible)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /partyman .

# --- Runtime stage ---
FROM alpine:3.20

RUN apk --no-cache add wget ca-certificates

WORKDIR /app
COPY --from=builder /partyman .
COPY data/ data/

EXPOSE 8081

CMD ["./partyman"]
