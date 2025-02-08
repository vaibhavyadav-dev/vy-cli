// here you'll find methods that will be executed at OS level
package sysCmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// update and upgrade Linux 
func Refresh() {
	if runtime.GOOS == "windows" {
		fmt.Println("This command is not supported on Windows")
		fmt.Println("you can run create issue and make pull request on github :)")
	} else {
		cmd := exec.Command("sudo", "apt", "update", "-y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error during update: %v\n", err)
			return
		}

		cmd = exec.Command("sudo", "apt", "upgrade", "-y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error during upgrade: %v\n", err)
			return
		}

		// Remove unused packages
		cmd = exec.Command("sudo", "apt", "autoremove", "-y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error during autoremove: %v\n", err)
			return
		}
	}
}