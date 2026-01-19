# Stage 1: Builder
# Use a full Go image for building the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go application as a static executable
RUN CGO_ENABLED=0 go build -ldflags='-s -w' -o quoter .

# Stage 2: Production
FROM scratch

WORKDIR /

COPY --from=builder /app/quoter .

EXPOSE 8080

CMD ["/quoter", "server", "start"]
