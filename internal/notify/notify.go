package notify

import (
	"fmt"

	"github.com/swarm-forge/swarm-forge/internal/swarmlog"
	"github.com/swarm-forge/swarm-forge/internal/tmux"
)

// Notify logs a message and sends it to a tmux pane.
func Notify(cmd tmux.Commander, logger *swarmlog.Logger, session string, pane int, message string) error {
	role := fmt.Sprintf("pane %d", pane)
	if err := logger.Write(role, message); err != nil {
		return err
	}
	return tmux.SendKeys(cmd, session, "swarm", pane, message)
}
