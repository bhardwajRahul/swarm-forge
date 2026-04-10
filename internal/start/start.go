package start

import (
	"github.com/swarm-forge/swarm-forge/internal/tmux"
)

// Config holds everything needed for the start sequence.
type Config struct {
	Commander   tmux.Commander
	Session     string
	ProjectRoot string
}

const window = "swarm"

// Run performs the full startup sequence.
func Run(cfg Config) error {
	if cfg.Commander.HasSession(cfg.Session) {
		if err := tmux.KillSession(cfg.Commander, cfg.Session); err != nil {
			return err
		}
	}
	if err := tmux.CreateSession(cfg.Commander, cfg.Session, window); err != nil {
		return err
	}
	if err := tmux.SplitGrid(cfg.Commander, cfg.Session, window); err != nil {
		return err
	}
	return tmux.SetPaneTitles(cfg.Commander, cfg.Session, window)
}
