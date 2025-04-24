# Go Commands

Run commands on your machine via MQTT and Home Assistant.

## Features

- Use MQTT to run commands on your machine.
- With MQTT discovery, automatically add commands to Home Assistant.

### Commands

#### Power

- Shutdown
- Restart
- Sleep
- Hibernate
- Lock Screen
- Logout
- Restart to Windows (Linux only)

## Installation

1. Install [Go](https://go.dev/doc/install).
1. Add go to your PATH/env.
1. Setup an MQTT server (usually via [Home Assistant](https://www.home-assistant.io/integrations/mqtt)).
1. Clone the repository.

```bash
git clone https://github.com/timmo001/go-commands.git
cd go-commands
```

1. Install the app

```bash
go install .
```

### Linux

1. Install the service using this script:

```bash
sudo ./setup-service.sh
```

### Windows

1. Create a startup shortcut to the executable using this script:

```powershell
.\setup-startup.ps1
```

## Usage

Once the app is installed and running, you can send commands to it via MQTT.

For Home Assistant, the MQTT integration will automatically discover the commands and add them to the UI.

[![Open your Home Assistant instance and show your integrations.](https://my.home-assistant.io/badges/integrations.svg)](https://my.home-assistant.io/redirect/integrations/)
