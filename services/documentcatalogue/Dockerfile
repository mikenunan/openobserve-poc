# --- Build stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache dependency downloads
COPY go.mod go.sum ./
RUN go mod download

# Build static binary (multi-arch compatible)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /documentcatalogue .

# --- Runtime stage ---
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /documentcatalogue .
COPY data/ data/

EXPOSE 8082

CMD ["./documentcatalogue"]
