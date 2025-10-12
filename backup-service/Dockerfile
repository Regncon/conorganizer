FROM golang:1.24.5 AS builder
# Enable CGO to improve SQLite compatibility
ENV CGO_ENABLED=1

# Install required libs
RUN apt-get update -y && apt-get install -y ca-certificates sqlite3 tzdata

# Set workdir for building application
WORKDIR /usr/src/codebase
COPY . .

# Install templ and generate required files
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

# Copy over & install required dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Compile binary for use in next step
RUN go build -o /usr/local/bin/backup-service .


# Compile flyctl binary from official image
FROM ghcr.io/superfly/flyctl:latest AS flyctl


FROM debian:bookworm AS server
# Create user for server
RUN groupadd -r regncon && useradd -r -g regncon -d /home/regncon -s /usr/sbin/nologin regncon

# Install required libs
RUN apt-get update -y && apt-get install -y ca-certificates sqlite3 tzdata

# Copy compiled binary from previous step
COPY --from=builder /usr/local/bin/backup-service /home/regncon/backup-service

# Copy over sql init file & static files
COPY --from=builder /usr/src/codebase/database/initialize.sql /usr/local/share/regncon/initialize.sql
COPY --from=builder /usr/src/codebase/web/static /home/regncon/web/static

# Copy flyctl binary
COPY --from=flyctl /flyctl /usr/local/bin/flyctl
RUN chmod 0755 /usr/local/bin/flyctl

# Create and set permissions for folders
RUN mkdir -p /home/regncon/.fly /home/regncon/.cache/flyctl /data/regncon/{tmp,logs,backup} \
    && chown -R regncon:regncon /home/regncon /data/regncon \
    # Optional: setgid so new files/dirs inherit the group
    && chmod 700 /home/regncon/.fly /home/regncon/.cache /home/regncon/.cache/flyctl \
    && chmod 2775 /data/regncon /data/regncon/{tmp,logs,backup}

# Switch workdir
WORKDIR /home/regncon

# Switch to non-root user
USER regncon
ENV HOME=/home/regncon

# Expose port for dashboard
EXPOSE 8080

# Run the app
CMD ["/home/regncon/backup-service"]
