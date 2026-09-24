$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
$targetPath = Join-Path $projectRoot ".agents/skills"
$linkPath = Join-Path $projectRoot ".claude/skills"

if (-not (Test-Path $targetPath)) {
    throw "Missing $targetPath"
}

New-Item -ItemType Directory -Force -Path (Join-Path $projectRoot ".claude") | Out-Null

if (Test-Path $linkPath) {
    Write-Host ".claude/skills already exists. Nothing to do."
    return
}

New-Item -ItemType Junction -Path $linkPath -Target $targetPath | Out-Null
Write-Host "Linked .claude/skills to .agents/skills"
