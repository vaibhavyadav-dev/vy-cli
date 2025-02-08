package feature

import (
	"bytes"
	"fmt"
	"os/exec"
)

func CommitAndStage(message string) error {
	gitInit()
	// check if all the changes staged or not!
	cmd := exec.Command("git", "add", ".")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("commit error: %w\nOutput: %s\nError: %s",
		err, stdout.String(), stderr.String())
	}
	
	// commit the changes
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("commit error: %w\nOutput: %s\nError: %s",
		err, stdout.String(), stderr.String())
	}
	return nil
}

// this will initialize the git repository
func gitInit() {
	_ = exec.Command("git", "init")
}