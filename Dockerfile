FROM golang:1.27.0 AS my-dev-environment

# Enable CGO (required for some C-based Go libraries)
ENV CGO_ENABLED=1

# Install necessary system dependencies
RUN apt-get update && \
    apt-get install -y curl sudo git sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create a user group named 'devuser' with group ID 1000
RUN groupadd -g 1000 devuser

# Create a user named 'devuser' with user ID 1000, assign it to the 'devuser' group, and set its shell to bash
RUN useradd -m -u 1000 -g devuser -s /bin/bash devuser

# Configure passwordless sudo access for the 'devuser' user
RUN echo "devuser ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/devuser && \
    chmod 0440 /etc/sudoers.d/devuser

# Set up the Go cache directories for the 'devuser' user
ENV GOMODCACHE=/home/devuser/go/pkg/mod
ENV GOCACHE=/home/devuser/.cache/go-build

# Set the working directory for the application
WORKDIR /home/devuser/app

# Create the Go cache directories and set ownership to 'devuser'
RUN mkdir -p /home/devuser/go/pkg/mod /home/devuser/.cache/go-build && \
    chown -R devuser:devuser /home/devuser/go /home/devuser/.cache

# Copy go.mod and go.sum with proper ownership
COPY --chown=devuser:devuser go.mod go.sum ./

# Compile the project-pinned tools straight into GOTOOLDIR so the existing
# `go tool` commands resolve them directly instead of rebuilding them. The
# temporary caches are removed in the same layer: module and build caches live
# in the compose volumes, so keeping them here would store them twice.
RUN tool_dir="$(go env GOTOOLDIR)" && \
    GOBIN="$tool_dir" GOCACHE=/tmp/tool-gocache GOMODCACHE=/tmp/tool-gomodcache go install \
    github.com/a-h/templ/cmd/templ \
    github.com/air-verse/air \
    github.com/go-task/task/v3/cmd/task && \
    rm -rf /tmp/tool-gocache /tmp/tool-gomodcache

USER devuser
RUN go tool templ version && \
    go tool air -v && \
    go tool task --version

# Caddy is the public Docker endpoint; the Go server stays internal.
EXPOSE 7332
