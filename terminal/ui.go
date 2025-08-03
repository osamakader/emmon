package terminal

import (
	"fmt"
	"strings"
	"time"

	"emmon/monitor"

	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
)

// TerminalUI handles the terminal interface
type TerminalUI struct {
	screen  tcell.Screen
	monitor *monitor.SystemMonitor
	log     *logrus.Logger
	quit    chan struct{}
}

// NewTerminalUI creates a new terminal UI instance
func NewTerminalUI(monitor *monitor.SystemMonitor, log *logrus.Logger) *TerminalUI {
	return &TerminalUI{
		monitor: monitor,
		log:     log,
		quit:    make(chan struct{}),
	}
}

// Start starts the terminal UI
func (tui *TerminalUI) Start() error {
	// Initialize tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("failed to create screen: %v", err)
	}

	if err := screen.Init(); err != nil {
		return fmt.Errorf("failed to initialize screen: %v", err)
	}

	tui.screen = screen
	defer screen.Fini()

	// Set up event handling
	go tui.handleEvents()

	// Main render loop
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tui.render()
		case <-tui.quit:
			return nil
		}
	}
}

// handleEvents handles keyboard and mouse events
func (tui *TerminalUI) handleEvents() {
	for {
		event := tui.screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				close(tui.quit)
				return
			}
		case *tcell.EventResize:
			tui.screen.Sync()
		}
	}
}

// render renders the current system stats
func (tui *TerminalUI) render() {
	tui.screen.Clear()

	// Get system stats
	stats, err := tui.monitor.GetSystemStats()
	if err != nil {
		tui.log.Errorf("Failed to get system stats: %v", err)
		return
	}

	// Get screen dimensions
	width, height := tui.screen.Size()

	// Draw header
	tui.drawHeader(width)

	// Draw CPU section
	tui.drawCPU(stats.CPU, 0, 3, width)

	// Draw Memory section
	tui.drawMemory(stats.Memory, 0, 12, width)

	// Draw Disk section
	tui.drawDisk(stats.Disk, 0, 21, width)

	// Draw Temperature section
	tui.drawTemperature(stats.Temperature, width/2, 3, width/2)

	// Draw GPIO section
	tui.drawGPIO(stats.GPIO, width/2, 12, width/2)

	// Draw footer
	tui.drawFooter(width, height)

	// Show the screen
	tui.screen.Show()
}

// drawHeader draws the application header
func (tui *TerminalUI) drawHeader(width int) {
	title := "🧠 Embedded Linux Monitor"
	subtitle := "Press ESC or Ctrl+C to exit"

	// Center the title
	titleX := (width - len(title)) / 2
	if titleX < 0 {
		titleX = 0
	}

	// Center the subtitle
	subtitleX := (width - len(subtitle)) / 2
	if subtitleX < 0 {
		subtitleX = 0
	}

	tui.drawText(titleX, 0, title, tcell.ColorGreen, tcell.ColorDefault, tcell.StyleDefault.Bold(true))
	tui.drawText(subtitleX, 1, subtitle, tcell.ColorGray, tcell.ColorDefault, tcell.StyleDefault)

	// Draw separator line
	separator := strings.Repeat("─", width)
	tui.drawText(0, 2, separator, tcell.ColorGray, tcell.ColorDefault, tcell.StyleDefault)
}

// drawCPU draws CPU information
func (tui *TerminalUI) drawCPU(cpu monitor.CPUStats, x, y, width int) {
	tui.drawText(x, y, "CPU", tcell.ColorYellow, tcell.ColorDefault, tcell.StyleDefault.Bold(true))

	// CPU Usage
	usageText := fmt.Sprintf("Usage: %6.1f%%", cpu.UsagePercent)
	tui.drawText(x, y+1, usageText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	// CPU Usage bar
	tui.drawProgressBar(x+15, y+1, cpu.UsagePercent, 20)

	// Load averages
	if len(cpu.LoadAverage) >= 3 {
		loadText := fmt.Sprintf("Load: %5.2f, %5.2f, %5.2f",
			cpu.LoadAverage[0], cpu.LoadAverage[1], cpu.LoadAverage[2])
		tui.drawText(x, y+2, loadText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)
	}

	// CPU Frequency
	freqText := fmt.Sprintf("Freq: %6.1f GHz", cpu.Frequency/1000)
	tui.drawText(x, y+3, freqText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)
}

// drawMemory draws memory information
func (tui *TerminalUI) drawMemory(mem monitor.MemStats, x, y, width int) {
	tui.drawText(x, y, "Memory", tcell.ColorYellow, tcell.ColorDefault, tcell.StyleDefault.Bold(true))

	// Memory Usage
	usageText := fmt.Sprintf("Usage: %6.1f%%", mem.UsagePercent)
	tui.drawText(x, y+1, usageText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	// Memory Usage bar
	tui.drawProgressBar(x+15, y+1, mem.UsagePercent, 20)

	// Memory details
	totalText := fmt.Sprintf("Total: %s", tui.formatBytes(mem.Total))
	tui.drawText(x, y+2, totalText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	usedText := fmt.Sprintf("Used:  %s", tui.formatBytes(mem.Used))
	tui.drawText(x, y+3, usedText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	freeText := fmt.Sprintf("Free:  %s", tui.formatBytes(mem.Free))
	tui.drawText(x, y+4, freeText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	availText := fmt.Sprintf("Avail: %s", tui.formatBytes(mem.Available))
	tui.drawText(x, y+5, availText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)
}

// drawDisk draws disk information
func (tui *TerminalUI) drawDisk(disk monitor.DiskStats, x, y, width int) {
	tui.drawText(x, y, "Disk", tcell.ColorYellow, tcell.ColorDefault, tcell.StyleDefault.Bold(true))

	// Disk Usage
	usageText := fmt.Sprintf("Usage: %6.1f%%", disk.UsagePercent)
	tui.drawText(x, y+1, usageText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	// Disk Usage bar
	tui.drawProgressBar(x+15, y+1, disk.UsagePercent, 20)

	// Disk details
	totalText := fmt.Sprintf("Total: %s", tui.formatBytes(disk.Total))
	tui.drawText(x, y+2, totalText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	usedText := fmt.Sprintf("Used:  %s", tui.formatBytes(disk.Used))
	tui.drawText(x, y+3, usedText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	freeText := fmt.Sprintf("Free:  %s", tui.formatBytes(disk.Free))
	tui.drawText(x, y+4, freeText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)

	ioText := fmt.Sprintf("I/O:   R:%s W:%s",
		tui.formatBytes(disk.IORead), tui.formatBytes(disk.IOWrite))
	tui.drawText(x, y+5, ioText, tcell.ColorWhite, tcell.ColorDefault, tcell.StyleDefault)
}

// drawTemperature draws temperature information
func (tui *TerminalUI) drawTemperature(temp monitor.TempStats, x, y, width int) {
	tui.drawText(x, y, "Temperature", tcell.ColorYellow, tcell.ColorDefault, tcell.StyleDefault.Bold(true))

	// CPU Temperature
	if temp.CPU > 0 {
		cpuTempText := fmt.Sprintf("CPU: %6.1f°C", temp.CPU)
		tui.drawText(x, y+1, cpuTempText, tui.getTempColor(temp.CPU), tcell.ColorDefault, tcell.StyleDefault)
	}

	// GPU Temperature
	if temp.GPU > 0 {
		gpuTempText := fmt.Sprintf("GPU: %6.1f°C", temp.GPU)
		tui.drawText(x, y+2, gpuTempText, tui.getTempColor(temp.GPU), tcell.ColorDefault, tcell.StyleDefault)
	}

	// Board Temperature
	if temp.Board > 0 {
		boardTempText := fmt.Sprintf("Board: %6.1f°C", temp.Board)
		tui.drawText(x, y+3, boardTempText, tui.getTempColor(temp.Board), tcell.ColorDefault, tcell.StyleDefault)
	}

	// Ambient Temperature
	if temp.Ambient > 0 {
		ambientTempText := fmt.Sprintf("Ambient: %6.1f°C", temp.Ambient)
		tui.drawText(x, y+4, ambientTempText, tui.getTempColor(temp.Ambient), tcell.ColorDefault, tcell.StyleDefault)
	}
}

// drawGPIO draws GPIO information
func (tui *TerminalUI) drawGPIO(gpio monitor.GPIOStats, x, y, width int) {
	tui.drawText(x, y, "GPIO Status", tcell.ColorYellow, tcell.ColorDefault, tcell.StyleDefault.Bold(true))

	if len(gpio.Pins) == 0 {
		tui.drawText(x, y+1, "No GPIO data", tcell.ColorGray, tcell.ColorDefault, tcell.StyleDefault)
		return
	}

	row := 1
	for pinName, pinData := range gpio.Pins {
		if row >= 8 { // Limit display to 8 pins
			break
		}

		pinText := fmt.Sprintf("%s: %d (%s)", pinName, pinData.Value, pinData.Mode)
		color := tcell.ColorRed
		if pinData.Value == 1 {
			color = tcell.ColorGreen
		}

		tui.drawText(x, y+row, pinText, color, tcell.ColorDefault, tcell.StyleDefault)
		row++
	}
}

// drawFooter draws the footer with timestamp
func (tui *TerminalUI) drawFooter(width, height int) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	timestampText := fmt.Sprintf("Last updated: %s", timestamp)

	// Draw at bottom of screen
	tui.drawText(0, height-1, timestampText, tcell.ColorGray, tcell.ColorDefault, tcell.StyleDefault)
}

// drawText draws text at the specified position
func (tui *TerminalUI) drawText(x, y int, text string, fg, bg tcell.Color, style tcell.Style) {
	style = style.Foreground(fg).Background(bg)

	for i, char := range text {
		if x+i >= 0 {
			tui.screen.SetContent(x+i, y, char, nil, style)
		}
	}
}

// drawProgressBar draws a progress bar
func (tui *TerminalUI) drawProgressBar(x, y int, percentage float64, width int) {
	filled := int((percentage / 100.0) * float64(width))

	// Draw filled part
	for i := 0; i < filled && i < width; i++ {
		tui.screen.SetContent(x+i, y, '█', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
	}

	// Draw empty part
	for i := filled; i < width; i++ {
		tui.screen.SetContent(x+i, y, '░', nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
	}
}

// getTempColor returns color based on temperature
func (tui *TerminalUI) getTempColor(temp float64) tcell.Color {
	switch {
	case temp < 40:
		return tcell.ColorGreen
	case temp < 60:
		return tcell.ColorYellow
	case temp < 80:
		return tcell.ColorOrange
	default:
		return tcell.ColorRed
	}
}

// formatBytes formats bytes into human readable format
func (tui *TerminalUI) formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
