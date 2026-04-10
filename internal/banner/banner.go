package banner

import (
	"fmt"
	"io"
)

// Print writes the SwarmForge startup banner to w.
func Print(w io.Writer) {
	fmt.Fprintln(w, "  ╔═══════════════════════════════════════════════╗")
	fmt.Fprintln(w, "  ║           SwarmForge v1.0 Starting            ║")
	fmt.Fprintln(w, "  ║   Disciplined agents build better software    ║")
	fmt.Fprintln(w, "  ╚═══════════════════════════════════════════════╝")
}
