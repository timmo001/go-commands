# Start Go Commands Service
$env:PATH += ";$env:USERPROFILE\go\bin"
$WorkingDir = "$env:USERPROFILE\.local\go-commands"

# Create working directory if it doesn't exist
if (-not (Test-Path -Path $WorkingDir)) {
    New-Item -ItemType Directory -Path $WorkingDir -Force
}

# Copy .env file to working directory
Copy-Item -Path .env -Destination $WorkingDir -Force

# Set working directory
Set-Location -Path $WorkingDir

# Start the application
Start-Process -FilePath "go-commands" -WorkingDirectory $WorkingDir -NoNewWindow
