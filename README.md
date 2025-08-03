# 🧠 Embedded Linux Monitor (emmon)

A lightweight system monitor for embedded Linux devices, built in Go. Like `htop`, but tailored specifically for embedded systems with GPIO monitoring capabilities.

## Features

- **Real-time System Monitoring**
  - CPU usage, load averages, and frequency
  - Memory usage and availability
  - Disk usage and I/O statistics
  - Temperature monitoring (CPU, GPU, Board, Ambient)
  - GPIO pin status monitoring

- **Multiple Interfaces**
  - **Web UI**: Modern web interface with WebSocket real-time updates
  - **Terminal UI**: Full-screen terminal interface using tcell (like htop)

- **Lightweight & Optimized**
  - Minimal resource usage for embedded systems
  - Direct access to `/proc`, `/sys`, and ioctl calls
  - No heavy dependencies

## Installation

### Prerequisites

- Go 1.21 or later
- Linux system (for system monitoring features)

### Build

```bash
# Clone the repository
git clone <repository-url>
cd emmon

# Build the application
go build -o emmon .

# Or use make
make build
```

### Cross-compilation for ARM

```bash
# For ARM64 (Raspberry Pi, etc.)
GOOS=linux GOARCH=arm64 go build -o emmon-arm64 .

# For ARM32
GOOS=linux GOARCH=arm go build -o emmon-arm .
```

## Usage

### Web Interface

Start the web interface (default port 8080):

```bash
./emmon web
```

Or specify a custom port:

```bash
./emmon web --port 9090
```

Access the web interface at `http://localhost:8080`

### Terminal Interface

Start the terminal interface:

```bash
./emmon terminal
```

Use `ESC` or `Ctrl+C` to exit.

### Configuration

Create a configuration file `~/.emmon.yaml`:

```yaml
log:
  level: info  # debug, info, warn, error

web:
  port: 8080
```

## System Requirements

### Linux Kernel Features

The monitor uses the following Linux kernel interfaces:

- `/proc/loadavg` - Load averages
- `/proc/cpuinfo` - CPU information
- `/proc/meminfo` - Memory statistics
- `/proc/diskstats` - Disk I/O statistics
- `/sys/class/thermal/thermal_zone*/temp` - Temperature sensors
- `/sys/class/gpio/*` - GPIO pin status

### GPIO Access

For GPIO monitoring, ensure:

1. GPIO sysfs interface is enabled
2. User has read access to `/sys/class/gpio/`
3. GPIO pins are exported (if needed)

```bash
# Example: Export GPIO pin 18
echo 18 > /sys/class/gpio/export
echo in > /sys/class/gpio/gpio18/direction
```

## Architecture

```
emmon/
├── main.go              # CLI entry point
├── monitor/
│   └── system.go        # Core system monitoring
├── web/
│   ├── server.go        # Web server with WebSocket
│   └── templates.go     # HTML templates
├── terminal/
│   └── ui.go           # Terminal UI with tcell
└── go.mod              # Go module dependencies
```

### Key Components

- **SystemMonitor**: Core monitoring logic using `/proc` and `/sys`
- **WebServer**: HTTP server with WebSocket for real-time updates
- **TerminalUI**: Full-screen terminal interface using tcell

## Development

### Dependencies

- `github.com/shirou/gopsutil/v3` - System statistics
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/gdamore/tcell/v2` - Terminal UI
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `github.com/sirupsen/logrus` - Logging

### Building

```bash
# Install dependencies
go mod tidy

# Build
go build -o emmon .

# Run tests
go test ./...

# Run with race detection
go test -race ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific package
go test ./monitor

# Run with verbose output
go test -v ./...
```

## Deployment

### Systemd Service

Create `/etc/systemd/system/emmon.service`:

```ini
[Unit]
Description=Embedded Linux Monitor
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/emmon web --port 8080
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl enable emmon
sudo systemctl start emmon
```

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o emmon .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/emmon .
EXPOSE 8080
CMD ["./emmon", "web"]
```

## Troubleshooting

### Common Issues

1. **Permission Denied**: Ensure the user has read access to `/proc` and `/sys`
2. **GPIO Not Found**: Check if GPIO sysfs interface is enabled
3. **Temperature Sensors**: Verify thermal zones exist in `/sys/class/thermal/`

### Debug Mode

Run with debug logging:

```bash
./emmon web --log-level debug
```

### Port Already in Use

```bash
# Check what's using the port
sudo netstat -tlnp | grep :8080

# Use a different port
./emmon web --port 9090
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by `htop` and similar system monitors
- Uses `gopsutil` for cross-platform system statistics
- Built with modern Go practices and libraries 