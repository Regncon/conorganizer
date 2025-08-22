FROM golang:1.24.5 AS builder
# Enable CGO to improve SQLite compatibility
ENV CGO_ENABLED=1

# Install required libs
RUN apt-get update -y && apt-get install -y ca-certificates sqlite3 tzdata

# Set workdir for building application
WORKDIR /usr/src/codebase
COPY . .

# Copy over & install required dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Compile binary for use in next step
RUN go build -o /usr/local/bin/backup-service .



FROM debian:bookworm AS server
# Create user for server
RUN groupadd -r regncon && useradd -r -g regncon -d /srv/regncon -s /usr/sbin/nologin regncon

# Install required libs
RUN apt-get update -y && apt-get install -y ca-certificates sqlite3 tzdata

# Copy compiled binary from previous step
COPY --from=builder /usr/local/bin/backup-service /home/regncon/backup-service

# Copy over sql init file
COPY --from=builder /usr/src/codebase/database/initialize.sql /usr/local/share/regncon/initialize.sql

RUN mkdir -p /data/regncon/{tmp,logs,backup} && \
    chown -R regncon:regncon /data/regncon/{tmp,logs,backup}

# Switch to non-root user
USER regncon

# Expose port for dashboard
EXPOSE 8080

# Run the app
CMD ["/home/regncon/backup-service"]
