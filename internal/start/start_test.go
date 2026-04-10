package start_test

import (
	"strings"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/start"
)

type recCmd struct {
	calls    [][]string
	sessions map[string]bool
}

func newRecCmd() *recCmd {
	return &recCmd{sessions: make(map[string]bool)}
}

func (r *recCmd) Run(args ...string) error {
	r.calls = append(r.calls, args)
	return nil
}

func (r *recCmd) HasSession(name string) bool {
	return r.sessions[name]
}

func hasCall(calls [][]string, keyword string) bool {
	for _, c := range calls {
		for _, a := range c {
			if strings.Contains(a, keyword) {
				return true
			}
		}
	}
	return false
}

func TestRunKillsExistingSession(t *testing.T) {
	cmd := newRecCmd()
	cmd.sessions["swarmforge"] = true
	cfg := start.Config{
		Commander:   cmd,
		Session:     "swarmforge",
		ProjectRoot: "/project",
	}
	err := start.Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasCall(cmd.calls, "kill-session") {
		t.Fatal("should kill existing session")
	}
	if !hasCall(cmd.calls, "new-session") {
		t.Fatal("should create new session")
	}
}

func TestRunNoExistingSession(t *testing.T) {
	cmd := newRecCmd()
	cfg := start.Config{
		Commander:   cmd,
		Session:     "swarmforge",
		ProjectRoot: "/project",
	}
	err := start.Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasCall(cmd.calls, "kill-session") {
		t.Fatal("should not kill when no session exists")
	}
	if !hasCall(cmd.calls, "new-session") {
		t.Fatal("should create new session")
	}
}

func TestRunSplitsGrid(t *testing.T) {
	cmd := newRecCmd()
	cfg := start.Config{
		Commander:   cmd,
		Session:     "swarmforge",
		ProjectRoot: "/project",
	}
	err := start.Run(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	count := 0
	for _, c := range cmd.calls {
		for _, a := range c {
			if strings.Contains(a, "split-window") {
				count++
				break
			}
		}
	}
	if count != 3 {
		t.Fatalf("expected 3 splits, got %d", count)
	}
}
