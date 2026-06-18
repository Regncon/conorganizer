FROM golang:1.26.3 AS my-dev-environment

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

# Switch to the 'devuser' user for subsequent commands
USER devuser

# Copy go.mod and go.sum with proper ownership
COPY --chown=devuser:devuser go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Warm project-pinned Go tools from go.mod.
RUN go tool templ version && \
    go tool air -v && \
    go tool task --version

# Expose ports used by the application:
# - 8080: Application's main HTTP server
# - 7331: Air live reload proxy
EXPOSE 8080 7331
