package swarmlog_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/swarmlog"
)

func TestWriteFormatsRoleAndMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := swarmlog.New(&buf)
	err := logger.Write("Architect", "task started")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[Architect] task started") {
		t.Fatalf("missing formatted entry: %s", buf.String())
	}
}

func TestWriteMultipleWriters(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	logger := swarmlog.New(&buf1, &buf2)
	err := logger.Write("Coder", "done")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf1.String(), "[Coder] done") {
		t.Fatalf("writer 1 missing entry: %s", buf1.String())
	}
	if !strings.Contains(buf2.String(), "[Coder] done") {
		t.Fatalf("writer 2 missing entry: %s", buf2.String())
	}
}
