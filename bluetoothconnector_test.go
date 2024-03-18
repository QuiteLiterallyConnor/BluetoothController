package bluetoothmanager

import (
	"testing"
)

func TestNewBluetoothConnector(t *testing.T) {
	connector, err := NewBluetoothConnector()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if connector == nil {
		t.Fatalf("Expected non-nil connector, got nil")
	}
}
