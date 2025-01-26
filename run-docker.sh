#!/bin/bash

# Define variables
imageName="my-dev-environment"
dockerFilePath="Dockerfile" # Specify the correct Dockerfile
workDir=$(pwd) # Use the current directory
containerName="dev-container"
containerWorkDir="/home/devuser/app"

# Check if the Docker image exists
imageExists=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep -w "$imageName:latest")

if [[ -z "$imageExists" ]]; then
    echo "Docker image '$imageName' does not exist. Building the image using $dockerFilePath..."
    docker build -t "$imageName" -f "$dockerFilePath" "$workDir"
    if [[ $? -ne 0 ]]; then
        echo "Failed to build Docker image. Exiting..." >&2
        exit 1
    fi
    echo "Docker image '$imageName' built successfully."
fi

# Run the Docker container
docker run -it --rm \
    --name "$containerName" \
    --hostname "$containerName" \
    -v "${workDir}:${containerWorkDir}" \
    -p 8080:8080 \
    -p 7331:7331 \
    "$imageName"
