package monitor

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
)

// SystemStats represents the current system statistics
type SystemStats struct {
	Timestamp   time.Time `json:"timestamp"`
	CPU         CPUStats  `json:"cpu"`
	Memory      MemStats  `json:"memory"`
	Disk        DiskStats `json:"disk"`
	Temperature TempStats `json:"temperature"`
	GPIO        GPIOStats `json:"gpio"`
}

// CPUStats represents CPU information
type CPUStats struct {
	UsagePercent float64   `json:"usage_percent"`
	LoadAverage  []float64 `json:"load_average"`
	Temperature  float64   `json:"temperature"`
	Frequency    float64   `json:"frequency"`
}

// MemStats represents memory information
type MemStats struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	Available    uint64  `json:"available"`
	UsagePercent float64 `json:"usage_percent"`
}

// DiskStats represents disk information
type DiskStats struct {
	Total        uint64  `json:"total"`
	Used         uint64  `json:"used"`
	Free         uint64  `json:"free"`
	UsagePercent float64 `json:"usage_percent"`
	IORead       uint64  `json:"io_read"`
	IOWrite      uint64  `json:"io_write"`
}

// TempStats represents temperature information
type TempStats struct {
	CPU     float64 `json:"cpu"`
	GPU     float64 `json:"gpu"`
	Board   float64 `json:"board"`
	Ambient float64 `json:"ambient"`
}

// GPIOStats represents GPIO pin status
type GPIOStats struct {
	Pins map[string]GPIOState `json:"pins"`
}

// GPIOState represents the state of a GPIO pin
type GPIOState struct {
	Pin   string `json:"pin"`
	Value int    `json:"value"`
	Mode  string `json:"mode"` // "in" or "out"
}

// SystemMonitor handles system monitoring
type SystemMonitor struct {
	log *logrus.Logger
}

// NewSystemMonitor creates a new system monitor instance
func NewSystemMonitor(log *logrus.Logger) *SystemMonitor {
	return &SystemMonitor{
		log: log,
	}
}

// GetSystemStats collects all system statistics
func (sm *SystemMonitor) GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{
		Timestamp: time.Now(),
	}

	// Collect CPU stats
	if cpuStats, err := sm.getCPUStats(); err == nil {
		stats.CPU = *cpuStats
	} else {
		sm.log.Warnf("Failed to get CPU stats: %v", err)
	}

	// Collect memory stats
	if memStats, err := sm.getMemoryStats(); err == nil {
		stats.Memory = *memStats
	} else {
		sm.log.Warnf("Failed to get memory stats: %v", err)
	}

	// Collect disk stats
	if diskStats, err := sm.getDiskStats(); err == nil {
		stats.Disk = *diskStats
	} else {
		sm.log.Warnf("Failed to get disk stats: %v", err)
	}

	// Collect temperature stats
	if tempStats, err := sm.getTemperatureStats(); err == nil {
		stats.Temperature = *tempStats
	} else {
		sm.log.Warnf("Failed to get temperature stats: %v", err)
	}

	// Collect GPIO stats
	if gpioStats, err := sm.getGPIOStats(); err == nil {
		stats.GPIO = *gpioStats
	} else {
		sm.log.Warnf("Failed to get GPIO stats: %v", err)
	}

	return stats, nil
}

// getCPUStats collects CPU information
func (sm *SystemMonitor) getCPUStats() (*CPUStats, error) {
	stats := &CPUStats{}

	// Get CPU usage percentage
	if usage, err := cpu.Percent(0, false); err == nil && len(usage) > 0 {
		stats.UsagePercent = usage[0]
	}

	// Get load average from /proc/loadavg
	if loadAvg, err := sm.readLoadAverage(); err == nil {
		stats.LoadAverage = loadAvg
	}

	// Get CPU frequency
	if freq, err := sm.readCPUFrequency(); err == nil {
		stats.Frequency = freq
	}

	return stats, nil
}

// getMemoryStats collects memory information
func (sm *SystemMonitor) getMemoryStats() (*MemStats, error) {
	if vmstat, err := mem.VirtualMemory(); err != nil {
		return nil, err
	} else {
		return &MemStats{
			Total:        vmstat.Total,
			Used:         vmstat.Used,
			Free:         vmstat.Free,
			Available:    vmstat.Available,
			UsagePercent: vmstat.UsedPercent,
		}, nil
	}
}

// getDiskStats collects disk information
func (sm *SystemMonitor) getDiskStats() (*DiskStats, error) {
	// Get disk usage for root filesystem
	if usage, err := disk.Usage("/"); err != nil {
		return nil, err
	} else {
		stats := &DiskStats{
			Total:        usage.Total,
			Used:         usage.Used,
			Free:         usage.Free,
			UsagePercent: usage.UsedPercent,
		}

		// Get disk I/O stats from /proc/diskstats
		if ioStats, err := sm.readDiskIO(); err == nil {
			stats.IORead = ioStats.Read
			stats.IOWrite = ioStats.Write
		}

		return stats, nil
	}
}

// getTemperatureStats collects temperature information
func (sm *SystemMonitor) getTemperatureStats() (*TempStats, error) {
	stats := &TempStats{}

	// Common temperature sensor paths
	tempPaths := map[string]string{
		"cpu":     "/sys/class/thermal/thermal_zone0/temp",
		"gpu":     "/sys/class/thermal/thermal_zone1/temp",
		"board":   "/sys/class/thermal/thermal_zone2/temp",
		"ambient": "/sys/class/thermal/thermal_zone3/temp",
	}

	for sensor, path := range tempPaths {
		if temp, err := sm.readTemperature(path); err == nil {
			switch sensor {
			case "cpu":
				stats.CPU = temp
			case "gpu":
				stats.GPU = temp
			case "board":
				stats.Board = temp
			case "ambient":
				stats.Ambient = temp
			}
		}
	}

	return stats, nil
}

// getGPIOStats collects GPIO pin status
func (sm *SystemMonitor) getGPIOStats() (*GPIOStats, error) {
	stats := &GPIOStats{
		Pins: make(map[string]GPIOState),
	}

	// Check for GPIO sysfs interface
	gpioPath := "/sys/class/gpio"
	if _, err := os.Stat(gpioPath); os.IsNotExist(err) {
		return stats, nil // GPIO not available
	}

	// Read GPIO pins
	if files, err := ioutil.ReadDir(gpioPath); err == nil {
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "gpio") {
				pinName := file.Name()
				if value, mode, err := sm.readGPIOState(filepath.Join(gpioPath, pinName)); err == nil {
					stats.Pins[pinName] = GPIOState{
						Pin:   pinName,
						Value: value,
						Mode:  mode,
					}
				}
			}
		}
	}

	return stats, nil
}

// readLoadAverage reads load average from /proc/loadavg
func (sm *SystemMonitor) readLoadAverage() ([]float64, error) {
	data, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid loadavg format")
	}

	loads := make([]float64, 3)
	for i := 0; i < 3; i++ {
		if loads[i], err = strconv.ParseFloat(fields[i], 64); err != nil {
			return nil, err
		}
	}

	return loads, nil
}

// readCPUFrequency reads CPU frequency from /proc/cpuinfo
func (sm *SystemMonitor) readCPUFrequency() (float64, error) {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu MHz") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				freqStr := strings.TrimSpace(parts[1])
				if freq, err := strconv.ParseFloat(freqStr, 64); err == nil {
					return freq, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("cpu frequency not found")
}

// readTemperature reads temperature from a sensor file
func (sm *SystemMonitor) readTemperature(path string) (float64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	tempStr := strings.TrimSpace(string(data))
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}

	// Convert from millidegrees to degrees Celsius
	return temp / 1000.0, nil
}

// readGPIOState reads the state of a GPIO pin
func (sm *SystemMonitor) readGPIOState(gpioPath string) (int, string, error) {
	// Read direction
	directionPath := filepath.Join(gpioPath, "direction")
	direction, err := ioutil.ReadFile(directionPath)
	if err != nil {
		return 0, "", err
	}
	mode := strings.TrimSpace(string(direction))

	// Read value
	valuePath := filepath.Join(gpioPath, "value")
	valueData, err := ioutil.ReadFile(valuePath)
	if err != nil {
		return 0, mode, err
	}

	value, err := strconv.Atoi(strings.TrimSpace(string(valueData)))
	if err != nil {
		return 0, mode, err
	}

	return value, mode, nil
}

// DiskIOStats represents disk I/O statistics
type DiskIOStats struct {
	Read  uint64 `json:"read"`
	Write uint64 `json:"write"`
}

// readDiskIO reads disk I/O statistics from /proc/diskstats
func (sm *SystemMonitor) readDiskIO() (*DiskIOStats, error) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 14 && (fields[2] == "sda" || fields[2] == "mmcblk0") {
			// Fields: major minor name reads reads_merged reads_sectors reads_time writes writes_merged writes_sectors writes_time
			reads, _ := strconv.ParseUint(fields[3], 10, 64)
			writes, _ := strconv.ParseUint(fields[7], 10, 64)

			return &DiskIOStats{
				Read:  reads,
				Write: writes,
			}, nil
		}
	}

	return &DiskIOStats{}, nil
}
