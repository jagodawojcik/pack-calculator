# ----- Build Stage ---------
FROM golang:1.25 AS builder

WORKDIR /app

# Dependencies
COPY go.mod go.sum* ./
RUN go mod download

# Build the Go binary
COPY . .
RUN go build -o server ./cmd/server

# ---- Runtime ---------
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
