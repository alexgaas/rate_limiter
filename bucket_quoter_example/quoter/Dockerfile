# Stage 1: Builder
# Use a full Go image for building the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ vendor/
# RUN go mod download

COPY . .

# Build the Go application as a static executable
RUN go build -mod=vendor -o quoter .

# Stage 2: Production
FROM scratch

WORKDIR /

COPY --from=builder /app/quoter .

EXPOSE 8443

CMD ["./quoter", "server", "start"]
