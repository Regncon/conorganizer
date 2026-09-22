$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$certificateDirectory = Join-Path $projectRoot "tmp/caddy"
$certificatePath = Join-Path $certificateDirectory "conorganizer-caddy-root.crt"

New-Item -ItemType Directory -Force -Path $certificateDirectory | Out-Null

Push-Location $projectRoot
try {
    docker compose cp caddy:/data/caddy/pki/authorities/local/root.crt $certificatePath
    if ($LASTEXITCODE -ne 0) {
        throw "Could not copy the Caddy root certificate. Start Docker Compose first."
    }

    Import-Certificate -FilePath $certificatePath -CertStoreLocation "Cert:\CurrentUser\Root" | Out-Null
    Write-Host "Trusted the Conorganizer Docker certificate for the current Windows user."
    Write-Host "Open https://localhost on the port configured by HTTPS_PORT (default: 7331)."
}
finally {
    Pop-Location
}
