package prompt_test

import (
	"strings"
	"testing"

	"github.com/swarm-forge/swarm-forge/internal/prompt"
)

func TestArchitectInstructionsNotEmpty(t *testing.T) {
	if prompt.ArchitectInstructions == "" {
		t.Fatal("ArchitectInstructions is empty")
	}
}

func TestCoderInstructionsNotEmpty(t *testing.T) {
	if prompt.CoderInstructions == "" {
		t.Fatal("CoderInstructions is empty")
	}
}

func TestE2EInterpreterInstructionsNotEmpty(t *testing.T) {
	if prompt.E2EInterpreterInstructions == "" {
		t.Fatal("E2EInterpreterInstructions is empty")
	}
}

func TestBuildContainsRole(t *testing.T) {
	cfg := prompt.AgentConfig{
		Role:         "Architect",
		Instructions: prompt.ArchitectInstructions,
		Session:      "swarmforge",
		ProjectRoot:  "/project",
	}
	result := prompt.Build(cfg, "Constitution text")
	if !strings.Contains(result, "You are the Architect agent") {
		t.Fatalf("prompt missing role header: %s", result)
	}
}

func TestBuildContainsConstitution(t *testing.T) {
	cfg := prompt.AgentConfig{
		Role:         "Coder",
		Instructions: prompt.CoderInstructions,
		Session:      "swarmforge",
		ProjectRoot:  "/project",
	}
	result := prompt.Build(cfg, "Rule 1: TDD")
	if !strings.Contains(result, "Rule 1: TDD") {
		t.Fatalf("prompt missing constitution")
	}
}

func TestBuildContainsCoordination(t *testing.T) {
	cfg := prompt.AgentConfig{
		Role:         "Coder",
		Instructions: prompt.CoderInstructions,
		Session:      "swarmforge",
		ProjectRoot:  "/project",
	}
	result := prompt.Build(cfg, "")
	if !strings.Contains(result, "notify-agent.sh") {
		t.Fatal("prompt missing notify-agent.sh")
	}
	if !strings.Contains(result, "swarm-log.sh") {
		t.Fatal("prompt missing swarm-log.sh")
	}
	if !strings.Contains(result, "agent_context/") {
		t.Fatal("prompt missing agent_context/")
	}
	if !strings.Contains(result, "Pane 0 = Architect") {
		t.Fatal("prompt missing pane layout")
	}
}
