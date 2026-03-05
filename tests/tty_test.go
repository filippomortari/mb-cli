package tests

import (
	"testing"

	"github.com/andreagrandi/mb-cli/internal/cli"
)

func TestIsTTYReturnsBool(t *testing.T) {
	result := cli.IsTTY()
	// In test environments, stdout is typically not a TTY (piped)
	if result {
		t.Log("IsTTY() returned true (running in a terminal)")
	} else {
		t.Log("IsTTY() returned false (stdout is piped)")
	}
}

func TestIsTTYOverrideable(t *testing.T) {
	original := cli.IsTTY
	defer func() { cli.IsTTY = original }()

	cli.IsTTY = func() bool { return true }
	if !cli.IsTTY() {
		t.Error("expected IsTTY to return true after override")
	}

	cli.IsTTY = func() bool { return false }
	if cli.IsTTY() {
		t.Error("expected IsTTY to return false after override")
	}
}
