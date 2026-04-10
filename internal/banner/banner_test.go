package banner_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/banner"
)

func TestPrintContainsSwarmForge(t *testing.T) {
	var buf bytes.Buffer
	banner.Print(&buf)
	if !strings.Contains(buf.String(), "SwarmForge") {
		t.Fatalf("banner missing SwarmForge: %s", buf.String())
	}
}

func TestPrintContainsMotto(t *testing.T) {
	var buf bytes.Buffer
	banner.Print(&buf)
	if !strings.Contains(buf.String(), "Disciplined agents build better software") {
		t.Fatalf("banner missing motto: %s", buf.String())
	}
}
