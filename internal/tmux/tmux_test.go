package tmux_test

import (
	"strings"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/tmux"
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

func countCalls(calls [][]string, keyword string) int {
	n := 0
	for _, c := range calls {
		for _, a := range c {
			if strings.Contains(a, keyword) {
				n++
				break
			}
		}
	}
	return n
}

func TestCreateSession(t *testing.T) {
	cmd := newRecCmd()
	err := tmux.CreateSession(cmd, "sf", "swarm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasCall(cmd.calls, "new-session") {
		t.Fatal("missing new-session call")
	}
}

func TestSplitGrid(t *testing.T) {
	cmd := newRecCmd()
	err := tmux.SplitGrid(cmd, "sf", "swarm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if countCalls(cmd.calls, "split-window") != 3 {
		t.Fatalf("expected 3 split-window calls, got %d", countCalls(cmd.calls, "split-window"))
	}
}

func TestSetPaneTitles(t *testing.T) {
	cmd := newRecCmd()
	err := tmux.SetPaneTitles(cmd, "sf", "swarm")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasCall(cmd.calls, "select-pane") {
		t.Fatal("missing select-pane call")
	}
}

func TestKillSession(t *testing.T) {
	cmd := newRecCmd()
	err := tmux.KillSession(cmd, "sf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasCall(cmd.calls, "kill-session") {
		t.Fatal("missing kill-session call")
	}
}

func TestLaunchAgent(t *testing.T) {
	cmd := newRecCmd()
	err := tmux.LaunchAgent(cmd, "sf", 0, "Architect", "/tmp/prompt.md", "/project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasCall(cmd.calls, "send-keys") {
		t.Fatal("missing send-keys call")
	}
	if !hasCall(cmd.calls, "SwarmForge Architect") {
		t.Fatal("missing agent name in command")
	}
	if !hasCall(cmd.calls, "--permission-mode acceptEdits") {
		t.Fatal("missing permission mode")
	}
}

func TestSendKeys(t *testing.T) {
	cmd := newRecCmd()
	err := tmux.SendKeys(cmd, "sf", "swarm", 3, "tail -f log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasCall(cmd.calls, "send-keys") {
		t.Fatal("missing send-keys call")
	}
	if !hasCall(cmd.calls, "tail -f log") {
		t.Fatal("missing command text")
	}
}
