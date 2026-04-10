package preflight

import "fmt"

// LookPathFunc matches the signature of exec.LookPath.
type LookPathFunc func(name string) (string, error)

// Check verifies that all named dependencies are available.
func Check(lookPath LookPathFunc, deps ...string) error {
	for _, dep := range deps {
		if _, err := lookPath(dep); err != nil {
			return fmt.Errorf("%s is required but not installed", dep)
		}
	}
	return nil
}
