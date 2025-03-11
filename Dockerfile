# Use a base image with Go pre-installed
FROM docker.io/golang:1.23.5 as my-dev-environment

# enable when using lib
# ENV CGO_ENABLED=1

# Install system dependencies and task runner
RUN apt-get update
RUN apt-get install -y curl sudo git sqlite3
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

# Install task runner
# RUN sh -c "$(curl --silent --location https://taskfile.dev/install.sh)" -- -d

# Copy over and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install project deeps
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/air-verse/air@latest
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Create a user and group named devuser
RUN groupadd -g 1000 devuser
RUN useradd -m -u 1000 -g devuser -s /bin/bash devuser

# Configure passwordless sudo for the devuser user
RUN echo "devuser ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/devuser
RUN chmod 0440 /etc/sudoers.d/devuser

# Switch to the 'devuser' user
USER devuser

# Set up the working directory for the user
WORKDIR /home/devuser/app

# Expose ports used by the application
EXPOSE 8080 7331