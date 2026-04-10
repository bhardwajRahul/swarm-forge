package preflight_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/preflight"
)

func TestCheckAllFound(t *testing.T) {
	lookPath := func(name string) (string, error) {
		return "/usr/bin/" + name, nil
	}
	err := preflight.Check(lookPath, "tmux", "claude")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheckFirstMissing(t *testing.T) {
	lookPath := func(name string) (string, error) {
		return "", errors.New(name + ": not found")
	}
	err := preflight.Check(lookPath, "tmux", "claude")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "tmux") {
		t.Fatalf("error should mention tmux: %v", err)
	}
}

func TestCheckSecondMissing(t *testing.T) {
	lookPath := func(name string) (string, error) {
		if name == "claude" {
			return "", errors.New("claude: not found")
		}
		return "/usr/bin/" + name, nil
	}
	err := preflight.Check(lookPath, "tmux", "claude")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "claude") {
		t.Fatalf("error should mention claude: %v", err)
	}
}

func TestCheckNoDeps(t *testing.T) {
	lookPath := func(name string) (string, error) {
		return "", errors.New("unreachable")
	}
	err := preflight.Check(lookPath)
	if err != nil {
		t.Fatalf("expected nil for no deps, got %v", err)
	}
}
