# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main /app/cmd/migrate/main.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/cmd/migrate/
COPY env env
COPY migration migration

CMD ["/app/cmd/migrate/main", "up"]