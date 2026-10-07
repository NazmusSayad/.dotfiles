package gmail

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func runGog(result any, args ...string) error {
	cmd := exec.Command("gog", append(args, "--json")...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("gog failed: %w", err)
	}
	if err := json.Unmarshal(output, result); err != nil {
		return fmt.Errorf("failed to parse gog output: %w", err)
	}
	return nil
}
