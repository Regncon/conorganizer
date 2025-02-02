# Define variables
$imageName = "my-dev-environment"
$dockerFilePath = "Dockerfile" # Specify the correct Dockerfile
$workDir = (Get-Location).Path # Use the current directory
$containerName = "dev-container"
$containerWorkDir = "/home/devuser/app"

# Check if the Docker image exists
$imageExists = docker images --format "{{.Repository}}:{{.Tag}}" | Find-Str "$imageName:latest"

if (-not $imageExists) {
    Write-Host "Docker image '$imageName' does not exist. Building the image using $dockerFilePath..."
    docker build -t $imageName -f $dockerFilePath $workDir
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to build Docker image. Exiting..."
        exit 1
    }
    Write-Host "Docker image '$imageName' built successfully."
}

# Run the Docker container
docker run -it --rm `
    --name $containerName `
    --hostname $containerName `
    -v "${workDir}:${containerWorkDir}" `
    -p 8080:8080 `
    -p 7331:7331 `
    --user 1000:1000 ` # Run the container as devuser
    $imageName
