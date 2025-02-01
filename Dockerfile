# Use a base image with Go pre-installed
FROM docker.io/golang:1.23.5 as my-dev-environment

# Set working directory
# WORKDIR /app

# Install system dependencies and task runner
RUN apt-get update && apt-get install -y \
    curl \
    sudo \
    git \
    sqlite3 \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Install task runner
# RUN sh -c "$(curl --silent --location https://taskfile.dev/install.sh)" -- -d

# Create a user and group named devuser
RUN groupadd -g 1000 devuser && \
    useradd -m -u 1000 -g devuser -s /bin/bash devuser

# Configure passwordless sudo for the devuser user
RUN echo "devuser ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/devuser && \
    chmod 0440 /etc/sudoers.d/devuser

# Create a directory for the application and give the devuser user ownership of the directory
RUN mkdir -p /home/devuser/app
RUN chown devuser:devuser -R /home/devuser/app

# Switch to the 'devuser' user
USER devuser

# Set up the working directory for the user
WORKDIR /home/devuser/app

# Copy go.mod and go.sum and download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# Install templ and air globally
RUN go install github.com/a-h/templ/cmd/templ@v0.3.819 && \
    go install github.com/go-task/task/v3/cmd/task@latest && \
    go install github.com/air-verse/air@latest

# Expose ports used by the application
EXPOSE 8080 7331

