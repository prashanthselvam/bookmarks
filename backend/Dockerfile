# Stage 1: Build
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod* go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o api ./cmd/api

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
CMD ["./api"]