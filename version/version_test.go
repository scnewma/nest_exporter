package version

import (
	"strings"
	"testing"
)

func TestVersionInfo(t *testing.T) {
	setup()

	expected := "(version=0.0.1)"

	info := Info()
	if info != expected {
		t.Errorf("expected: %s, got: %s\n", expected, info)
	}
}

func TestPrint(t *testing.T) {
	setup()

	result := Print()

	if !strings.Contains(result, Program) {
		t.Error("expected version string to contain program name")
	}

	if !strings.Contains(result, Version) {
		t.Error("expected version string to contain version")
	}

	if !strings.Contains(result, GoVersion) {
		t.Error("expected version string to contain go version")
	}
}

func setup() {
	Program = "nest_exporter"
	GoVersion = "1.11.0"
	Version = "0.0.1"
}
