package sysCmd

import (
	"fmt"
	"os"
	"os/exec"
)

// install Go, Node and Python
func SetupGoNodePython(){
	Refresh() // update and upgrade Linux

	setupGo() // v1.22.11
	setupPython() // v3.10.12
	setupNode() // v22.13.1
}

func setupGo(){
	
	// Check if Go is already installed
	cmd := exec.Command("go", "version")
	if err := cmd.Run(); err == nil {
		fmt.Println("Go is already installed.... Skipping installations...")
		return
	}
	
	
	fmt.Println("\n\n--------- Downloading Golang 1.22.11 ---------------")
	// // Download Go source
	cmd = exec.Command("wget", "https://go.dev/dl/go1.22.11.src.tar.gz")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error downloading Go source: %v\n", err)
		return
	}

	// Extract the archive
	fmt.Println("\n\n--------- Exracting Golang ---------------")
	cmd = exec.Command("tar", "-xzf", "go1.22.11.src.tar.gz")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error extracting Go source: %v\n", err)
		return
	}

	// Build and install Go
	fmt.Println("\n\n--------- Building Golang ---------------")
	cmd = exec.Command("./make.bash")
	cmd.Dir = "go/src"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error building Go: %v\n", err)
		return
	}

	
	fmt.Println("\n\n--------- Setting Golang Path ---------------")
	// Remove existing Go installation if it exists
	cmd = exec.Command("sudo", "rm", "-rf", "/usr/local/go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error removing existing Go installation: %v\n", err)
		return
	}

	// Move Go installation to /usr/local
	cmd = exec.Command("sudo", "mv", "go", "/usr/local/")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error moving Go installation: %v\n", err)
		return
	}
	// Clean up downloaded file
	os.Remove("go1.22.11.src.tar.gz")
}

func setupPython(){

	// Check if Python is already installed
	cmd := exec.Command("python3", "--version")
	if err := cmd.Run(); err == nil {
		fmt.Println("python is already installed.... Skipping installations...")
		return
	}


	fmt.Println("\n\n--------- Installing Python3 ---------------")
	// Install Python
	cmd = exec.Command("sudo", "apt", "install", "python3", "python3-pip", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error during Python installation: %v\n", err)
		return
	}
}

func setupNode(){

	// Check if Nodejs is already installed
	cmd := exec.Command("node", "-v")
	if err := cmd.Run(); err == nil {
		fmt.Println("node is already installed.... Skipping installations...")
		return
	}

	fmt.Println("\n--------- Installing NVM 22 ---------------")
	// Install NVM using nodesource
	cmd = exec.Command("bash", "-c", "curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error setting up Node.js repository: %v\n", err)
		return
	}

	fmt.Println("\n--------- Installing Nodejs 22 ---------------")
	// Install Node.js 22 (LTS)
	cmd = exec.Command("bash", "-c", "source ~/.nvm/nvm.sh && nvm install 22 && nvm use 22")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error installing Node.js: %v\n", err)
		return
	}
}