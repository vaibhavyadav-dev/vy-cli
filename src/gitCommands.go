package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

func CommitAndStage(message string) string {
	gitInit()
	// check if all the changes staged or not!
	cmd := exec.Command("git", "add", ".")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("commit error: %w\nOutput: %s\nError: %s",
		err, stdout.String(), stderr.String())
	}
	
	// commit the changes
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("commit error: %w\nOutput: %s\nError: %s",
		err, stdout.String(), stderr.String())
	}
	return "all changes has been Staged and Committed :)"
}

// this will initialize the git repository
func gitInit() {
	_ = exec.Command("git", "init")
}