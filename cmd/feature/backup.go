package feature

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func HandleBackup(verbose bool, absLocalFilePath, drive string) {

	// Check if the localFilePath is a file
	// If it is, upload the whole folder and it's content and return
	isFolder := len(strings.TrimSpace(absLocalFilePath)) > 0
	if isFolder {
		uploadFolder(absLocalFilePath, drive, verbose)
		return
	}

	// upload the config file
	if err := uploadConfigFile(verbose, drive); err != nil {
		log.Fatal(err)
	}

}

// Upload the configuration files to the specified drive using rclone
func uploadConfigFile(verbose bool, drive string) error {

	// Check if rclone is installed and configured
	err := checkRcloneInstallation(drive)
	if err != nil {
		return fmt.Errorf("rclone is not installed or configured")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory")
	}

	filesToBackup := make(map[string]string)
	possibleFiles := map[string]string{
		// Terminal and shell settings
		".bashrc":        filepath.Join(homeDir, ".bashrc"),
		".profile":       filepath.Join(homeDir, ".profile"),
		".zshrc":         filepath.Join(homeDir, ".zshrc"),
		".bash_aliases":  filepath.Join(homeDir, ".bash_aliases"),
		".zsh_aliases":   filepath.Join(homeDir, ".zsh_aliases"),
		"custom-scripts": filepath.Join(homeDir, "bin"), // Custom scripts directory

		// Desktop environment settings
		"dconf-settings": filepath.Join(homeDir, ".config/dconf"),
		"gnome-terminal": filepath.Join(homeDir, ".config/gnome-terminal"),
		"gtk-3.0":        filepath.Join(homeDir, ".config/gtk-3.0"),
		"gtk-4.0":        filepath.Join(homeDir, ".config/gtk-4.0"),
		"nautilus":       filepath.Join(homeDir, ".config/nautilus"),
		"autostart":      filepath.Join(homeDir, ".config/autostart"),

		// Theme and appearance
		"themes":      filepath.Join(homeDir, ".themes"),
		"icons":       filepath.Join(homeDir, ".icons"),
		"backgrounds": filepath.Join(homeDir, ".local/share/backgrounds"),
		"fonts":       filepath.Join(homeDir, ".fonts"),
		"local-fonts": filepath.Join(homeDir, ".local/share/fonts"),

		// Application preferences
		"mimeapps":         filepath.Join(homeDir, ".config/mimeapps.list"),
		"user-dirs":        filepath.Join(homeDir, ".config/user-dirs.dirs"),
		"preferences":      filepath.Join(homeDir, ".local/share/preferences"),
		"custom-launchers": filepath.Join(homeDir, ".local/share/applications"),
		"plank":            filepath.Join(homeDir, ".config/plank"),

		// SSH keys and configuration
		"ssh": filepath.Join(homeDir, ".ssh"),

		// Environment and system settings
		"pam-environment": filepath.Join(homeDir, ".pam_environment"),
		"xprofile":        filepath.Join(homeDir, ".xprofile"),
		"xinitrc":         filepath.Join(homeDir, ".xinitrc"),

		// Developer tools and editors
		"vimrc":     filepath.Join(homeDir, ".vimrc"),
		"vim":       filepath.Join(homeDir, ".vim"),
		"nvim":      filepath.Join(homeDir, ".config/nvim"),
		"tmux":      filepath.Join(homeDir, ".tmux.conf"),
		"gitconfig": filepath.Join(homeDir, ".gitconfig"),

		// Systemd user services
		"systemd-user": filepath.Join(homeDir, ".config/systemd/user"),
	}

	for name, path := range possibleFiles {
		if _, err := os.Stat(path); err == nil {
			filesToBackup[name] = path
		}
	}

	// Define backup directory on OneDrive
	backupDir := "Backups/ubuntu-settings"

	// Track backup status
	successCount := 0
	totalFiles := len(filesToBackup)

	fmt.Printf("Please wait.... I'm Uploading files to %s.....\nThis Will Take Time Depending Upon Speed of Internet and Size of Folder :) ...", drive)
	if verbose {
		fmt.Println("Starting Ubuntu settings backup...")
		fmt.Printf("Total configurations to backup: %d\n\n", totalFiles)
	}

	// Process each file individually
	for name, path := range filesToBackup {
		if _, err := os.Stat(path); err != nil {
			fmt.Printf("  ❌ Skipping: configuration not found\n\n")
			continue
		}

		if verbose {
			fmt.Printf("📤 Uploading %s to %s... ", name, drive)
		}

		err = rclone(path, filepath.Join(backupDir, name), drive)
		if err != nil {
			fmt.Printf("❌ Failed\n  Error: %v\n\n", err)
			continue
		}

		if verbose {
			fmt.Printf("✅ Success\n\n")
		}

		successCount++
	}

	fmt.Printf("Backup completed! Successfully backed up %d of %d configurations\n", successCount, totalFiles)

	return nil
}

// Upload a folder to the specified drive using rclone
func uploadFolder(folder, drive string, verbose bool) {
	if verbose {
		fmt.Printf("📤 Uploading folder (Local) %s to %s...", folder, drive)
	}

	folderName := filepath.Base(folder)
	remoteDir := filepath.Join("Backups", folderName)

	err := rclone(folder, remoteDir, drive)
	if err != nil {
		fmt.Printf("❌ Failed to upload folder (Local): %v\n", err)
		return
	}

	if verbose {
		fmt.Println("✅ successfull")
	}
}

// Upload a file to the specified drive using rclone
// localpath:  //remotePath		//driveName (oneDrive, googleDrive e.t.c)
func rclone(localFilePath, remotePath, drive string) error {
	// Check if the destination directory exists, if not create it
	checkCmd := exec.Command("rclone", "lsf", drive+filepath.Dir(remotePath))
	if err := checkCmd.Run(); err != nil {
		mkdirCmd := exec.Command("rclone", "mkdir", drive+filepath.Dir(remotePath))
		if err := mkdirCmd.Run(); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	cmd := exec.Command("rclone", "copy", localFilePath, drive+remotePath, "--progress")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("rclone error: %w\nOutput: %s\nError: %s",
			err, stdout.String(), stderr.String())
	}
	return nil
}

// Check if rclone is installed and configured with the specified drive
func checkRcloneInstallation(drive string) error {
	// Check if rclone is installed
	_, err := exec.LookPath("rclone")
	if err != nil {
		return fmt.Errorf("rclone is not installed. Please install it using:\ncurl https://rclone.org/install.sh | sudo bash\nThen authenticate with: rclone config")
	}

	// Check if rclone is configured with drive or not
	cmd := exec.Command("rclone", "listremotes")
	output, err := cmd.Output()
	if err != nil || !bytes.Contains(output, []byte(drive)) {
		return fmt.Errorf("rclone is not configured with %s.\nPlease run 'rclone config' and set up your OneDrive connection", drive)
	}
	return nil
}
