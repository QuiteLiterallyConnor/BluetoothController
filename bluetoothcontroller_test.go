package bluetoothmanager

import (
	"testing"
)

func TestNewBluetoothController(t *testing.T) {
	listener := func(e Event) {
		// This is a dummy listener for testing.
	}

	controller, err := NewBluetoothController(listener)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if controller == nil {
		t.Fatalf("Expected non-nil controller, got nil")
	}
}
