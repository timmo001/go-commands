#!/bin/bash

# Install efibootmgr

# Arch
if [ -f /etc/arch-release ]; then
    sudo pacman -S efibootmgr
fi

# Debian/Ubuntu
if [ -f /etc/debian_version ]; then
    sudo apt-get install efibootmgr
fi

# Install go-commands
go install .

# Create working directory if it doesn't exist
mkdir -p ~/.local/go-commands

# Copy .env file to working directory
cp .env ~/.local/go-commands

# Configure sudo privileges for efibootmgr
echo "# Allow go-commands to use efibootmgr without password
$USER ALL=(ALL) NOPASSWD: /usr/bin/efibootmgr" | sudo tee /etc/sudoers.d/go-commands

# Ensure the systemd user directory exists
mkdir -p ~/.config/systemd/user/

# Copy the service file
cp go-commands.service ~/.config/systemd/user/

# Reload systemd to recognize the new service
systemctl --user daemon-reload

# Enable and start the service
systemctl --user enable --now go-commands.service

echo "Service installed and started successfully!"

# Check the service is running
systemctl --user status go-commands.service
