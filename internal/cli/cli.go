package cli

import "fmt"

// Handler is called when a subcommand is matched.
type Handler func(args []string) error

// Config holds the registered subcommand handlers.
type Config struct {
	Start  Handler
	Notify Handler
	Log    Handler
}

// Dispatch routes CLI arguments to the appropriate handler.
func Dispatch(args []string, cfg Config) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: swarmforge <start|notify|log> [args...]")
	}
	rest := args[1:]
	switch args[0] {
	case "start":
		return cfg.Start(rest)
	case "notify":
		return cfg.Notify(rest)
	case "log":
		return cfg.Log(rest)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}
