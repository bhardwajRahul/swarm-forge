package tmux

import "fmt"

// Commander abstracts tmux shell commands for testability.
type Commander interface {
	Run(args ...string) error
	HasSession(name string) bool
}

// CreateSession creates a new tmux session with the given name and window.
func CreateSession(cmd Commander, session, window string) error {
	return cmd.Run("new-session", "-d", "-s", session, "-n", window)
}

// SplitGrid splits a window into a 2x2 grid of panes.
func SplitGrid(cmd Commander, session, window string) error {
	target := session + ":" + window
	err := cmd.Run("split-window", "-t", target+".0", "-h", "-p", "50")
	if err != nil {
		return err
	}
	err = cmd.Run("split-window", "-t", target+".0", "-v", "-p", "50")
	if err != nil {
		return err
	}
	return cmd.Run("split-window", "-t", target+".2", "-v", "-p", "50")
}

// SetPaneTitles sets the title for each pane and enables border display.
func SetPaneTitles(cmd Commander, session, window string) error {
	target := session + ":" + window
	titles := []string{"Architect", "E2E Interpreter", "Coder", "Metrics"}
	for i, title := range titles {
		pane := fmt.Sprintf("%s.%d", target, i)
		if err := cmd.Run("select-pane", "-t", pane, "-T", title); err != nil {
			return err
		}
	}
	if err := cmd.Run("set-option", "-t", session, "pane-border-status", "top"); err != nil {
		return err
	}
	if err := cmd.Run("set-option", "-t", session, "pane-border-format", " #{pane_title} "); err != nil {
		return err
	}
	return cmd.Run("set-window-option", "-t", target, "allow-rename", "off")
}

// KillSession kills an existing tmux session.
func KillSession(cmd Commander, session string) error {
	return cmd.Run("kill-session", "-t", session)
}

// LaunchAgent sends a claude command to the given pane.
func LaunchAgent(cmd Commander, session string, pane int, name, promptFile, projectRoot string) error {
	target := fmt.Sprintf("%s:swarm.%d", session, pane)
	command := fmt.Sprintf(
		"cd '%s' && claude --append-system-prompt-file '%s' --permission-mode acceptEdits -n 'SwarmForge %s'",
		projectRoot, promptFile, name,
	)
	return cmd.Run("send-keys", "-t", target, command, "Enter")
}

// SendKeys sends keystrokes to a tmux pane.
func SendKeys(cmd Commander, session, window string, pane int, keys string) error {
	target := fmt.Sprintf("%s:%s.%d", session, window, pane)
	return cmd.Run("send-keys", "-t", target, keys, "Enter")
}
