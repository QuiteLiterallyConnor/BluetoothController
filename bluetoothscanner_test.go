package bluetoothmanager

import (
	"testing"
)

func TestNewBluetoothScanner(t *testing.T) {
	listener := func(d Device) {
		// This is a dummy listener for testing.
	}

	scanner, err := NewBluetoothScanner(listener)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if scanner == nil {
		t.Fatalf("Expected non-nil scanner, got nil")
	}
}
