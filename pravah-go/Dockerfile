# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies (ca-certificates and tzdata for the final image)
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user in the builder stage
RUN addgroup -S pravah && adduser -S pravah -G pravah

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o server ./cmd/server/main.go

# Final stage - START FROM SCRATCH
FROM scratch

# Import the user/group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Import SSL certificates and Timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

# Copy the binary and assets from the builder
# We use the numeric UID (100) for security and to avoid looking up names in /etc/passwd
COPY --from=builder --chown=pravah:pravah /app/server .
COPY --from=builder --chown=pravah:pravah /app/public ./public

# Use the non-root user
USER pravah

EXPOSE 3000

# Start the application
ENTRYPOINT ["./server"]
