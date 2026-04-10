package notify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/notify"
	"github.com/swarm-forge/swarm-forge/internal/swarmlog"
)

type recCmd struct {
	calls    [][]string
	sessions map[string]bool
}

func newRecCmd() *recCmd {
	return &recCmd{sessions: map[string]bool{"sf": true}}
}

func (r *recCmd) Run(args ...string) error {
	r.calls = append(r.calls, args)
	return nil
}

func (r *recCmd) HasSession(name string) bool {
	return r.sessions[name]
}

func TestNotifyLogsAndSends(t *testing.T) {
	var buf bytes.Buffer
	logger := swarmlog.New(&buf)
	cmd := newRecCmd()
	err := notify.Notify(cmd, logger, "sf", 0, "hello architect")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[pane 0] hello architect") {
		t.Fatalf("missing log entry: %s", buf.String())
	}
	found := false
	for _, c := range cmd.calls {
		for _, a := range c {
			if strings.Contains(a, "send-keys") {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("missing send-keys call")
	}
}
