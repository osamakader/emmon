package monitor

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func newTestMonitor() *SystemMonitor {
	log := logrus.New()
	log.SetOutput(nil) // Silence output
	return NewSystemMonitor(log)
}

func TestGetCPUStats(t *testing.T) {
	m := newTestMonitor()
	stats, err := m.getCPUStats()
	if err != nil {
		t.Fatalf("getCPUStats error: %v", err)
	}
	if stats.UsagePercent < 0 || stats.UsagePercent > 100 {
		t.Errorf("CPU usage out of range: %v", stats.UsagePercent)
	}
}

func TestGetMemoryStats(t *testing.T) {
	m := newTestMonitor()
	stats, err := m.getMemoryStats()
	if err != nil {
		t.Fatalf("getMemoryStats error: %v", err)
	}
	if stats.Total == 0 {
		t.Error("Total memory should not be zero")
	}
	if stats.Used > stats.Total {
		t.Error("Used memory greater than total")
	}
}

func TestGetDiskStats(t *testing.T) {
	m := newTestMonitor()
	stats, err := m.getDiskStats()
	if err != nil {
		t.Fatalf("getDiskStats error: %v", err)
	}
	if stats.Total == 0 {
		t.Error("Total disk should not be zero")
	}
	if stats.Used > stats.Total {
		t.Error("Used disk greater than total")
	}
}

func TestGetTemperatureStats(t *testing.T) {
	m := newTestMonitor()
	_, err := m.getTemperatureStats()
	if err != nil {
		t.Logf("Temperature stats not available: %v (this is OK on some systems)", err)
	}
}
