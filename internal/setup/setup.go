package setup

// FS abstracts filesystem operations for testability.
type FS interface {
	MkdirAll(path string, perm uint32) error
	WriteFile(path string, data []byte, perm uint32) error
	ReadFile(path string) ([]byte, error)
	Stat(path string) (bool, error) // returns (exists, err)
}

// EnsureDirs creates the required project directories.
func EnsureDirs(fs FS, root string) error {
	dirs := []string{"features", "logs", "agent_context"}
	for _, d := range dirs {
		if err := fs.MkdirAll(root+"/"+d, 0o755); err != nil {
			return err
		}
	}
	return nil
}

const notifyScript = `#!/bin/bash
# Usage: ./notify-agent.sh <target-pane-index> "message"
# Panes: 0=Architect, 1=E2E Interpreter, 2=Coder, 3=Metrics
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
echo "[$TIMESTAMP] [pane $1] $2" >> logs/agent_messages.log
tmux send-keys -t swarmforge:swarm.$1 "$2" Enter
`

const logScript = `#!/bin/bash
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
echo "[$TIMESTAMP] [$1] $2" >> logs/agent_messages.log
echo "[$1] $2"
`

// WriteHelperScripts generates backward-compatible shell scripts.
func WriteHelperScripts(fs FS, root string) error {
	scripts := map[string]string{
		"notify-agent.sh": notifyScript,
		"swarm-log.sh":    logScript,
	}
	for name, content := range scripts {
		path := root + "/" + name
		if err := fs.WriteFile(path, []byte(content), 0o755); err != nil {
			return err
		}
	}
	return nil
}
