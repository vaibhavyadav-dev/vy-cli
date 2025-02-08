package feature

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func Extracter(AbsFilePath string) {
	// Check if the required tools are installed
	installExtracterTools()

	// Get the file extension
	ext := filepath.Ext(AbsFilePath)
	switch ext {
	case ".zip":
		if err := executeCommand("unzip", "-o", AbsFilePath); err != nil {
			log.Fatal(err)
		}
	case ".tar":
		if err := executeCommand("tar", "xvf", AbsFilePath); err != nil {
			log.Fatal(err)
		}
	case ".tar.gz":
		if err := executeCommand("tar", "xvf", AbsFilePath); err != nil {
			log.Fatal(err)
		}
	case ".tar.xz":
		if err := executeCommand("tar", "xvf", AbsFilePath); err != nil {
			log.Fatal(err)
		}
	case ".7z":
		if err := executeCommand("7z", "x", AbsFilePath); err != nil {
			log.Fatal(err)
		}
	case ".rar":
		if err := executeCommand("unrar", "x", AbsFilePath); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Println("Unsupported file format")
	}
}

func installExtracterTools() {
	cmds := []string{"unzip", "tar", "unrar", "7z"}
	for _, cmd := range cmds {
		if !isCommandAvailable(cmd) {
			fmt.Printf("%s is not installed....\nInstalling... Please wait...\n", cmd)
			// var stdout, stderr []b	yte
			if err := exec.Command("sudo", "apt", "install", "-y", cmd).Run(); err != nil {
				fmt.Printf("Error installing %s: %v\n", cmd, err)
			}
			fmt.Printf("Installed %s\n", cmd)
		}
	}
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Execute the command to extract the file
func executeCommand(tool, args, AbsFilePath string) error {
	cmd := exec.Command(tool, args, AbsFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("file is already extracted or Error extracting zip file: %v", err)
	}
	return nil
}