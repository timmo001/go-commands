#!/bin/bash

# Create working directory if it doesn't exist
mkdir -p ~/.local/go-commands

# Copy .env file to working directory
cp ~/.env ~/.local/go-commands/.env

# Ensure the systemd user directory exists
mkdir -p ~/.config/systemd/user/

# Copy the service file
cp go-commands.service ~/.config/systemd/user/

# Reload systemd to recognize the new service
systemctl --user daemon-reload

# Enable and start the service
systemctl --user enable --now go-commands.service

echo "Service installed and started successfully!"
