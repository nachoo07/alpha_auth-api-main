FROM golang:1.23.1 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/app ./cmd/api/main.go

# Final stage
FROM alpine:3.18

WORKDIR /usr/src/app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

ENV SSL_MODE=require
ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/app"]
