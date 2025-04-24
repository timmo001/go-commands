# Create Startup Shortcut for Go Commands Service

$WorkingDir = "$env:USERPROFILE\.local\go-commands"

# Create working directory if it doesn't exist
if (-not (Test-Path -Path $WorkingDir)) {
    New-Item -ItemType Directory -Path $WorkingDir -Force
}

# Get the startup folder path
$StartupFolder = [System.Environment]::GetFolderPath('Startup')
$ScriptPath = Join-Path $PSScriptRoot "start.ps1"
$ShortcutPath = Join-Path $StartupFolder "Go Commands.lnk"

# Create the shortcut
$WScriptShell = New-Object -ComObject WScript.Shell
$Shortcut = $WScriptShell.CreateShortcut($ShortcutPath)
$Shortcut.TargetPath = "powershell.exe"
$Shortcut.Arguments = "-ExecutionPolicy Bypass -WindowStyle Hidden -File `"$ScriptPath`""
$Shortcut.WorkingDirectory = $WorkingDir
$Shortcut.Description = "Start Go Commands Service"
$Shortcut.Save()

Write-Host "Startup shortcut created at: $ShortcutPath"
