### Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o dbsetup ./cmd/db
RUN CGO_ENABLED=0 GOOS=linux go build -o vrsapp ./cmd/web

## Runtime stage
FROM alpine

# Copy only the binaries and neccessary file from the build stage to the final image
COPY --from=builder /app/vrsapp /
COPY --from=builder /app/db/migrations /db/migrations
COPY --from=builder /app/db/seeds /db/seeds
COPY --from=builder /app/dbsetup /

# Run dbsetup first, then vrsapp
CMD ["/bin/sh", "-c", "/dbsetup && /vrsapp"]