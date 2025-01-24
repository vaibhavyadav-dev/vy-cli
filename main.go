package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"

	"github.com/vaibhavyadav-dev/vy-cli/src"
	"github.com/vaibhavyadav-dev/vy-cli/src/sysconfig"

	// Load the .env file
	"github.com/joho/godotenv"
)

//go:embed src/cmd.txt
var cmdFile string
//go:embed .env
var embeddedEnv string

func main() {
	if len(os.Args) < 2 {
		if cmdFile == "" {
			fmt.Println("Seems like package is not successfully installed :(") 
			fmt.Println("Please install it with default configuration")
			return
		}
		cmd.PrintRainbowGlowLargeText("Vaibhav Yadav")
		fmt.Println("Command line made For and By VAIBHAV YADAV")
		fmt.Println(cmdFile)
		return
	}
	
	// Write the embedded .env content to a temporary file
	tmpFile, err := os.CreateTemp("", "embedded-env-*.env")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(embeddedEnv)); err != nil {
		fmt.Println("Error writing to temp file:", err)
		return
	}
	tmpFile.Close()

	// Load the .env file from the temporary location
	if err := godotenv.Load(tmpFile.Name()); err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	
	command := os.Args[1]

	switch command {	
	case "date":
		fmt.Println(cmd.Date())
	case "backup":
		drive := "gdrive:"
		verbose := false
		folder := ""

		for i := 0; i < len(os.Args); i++ {

			// Check if the user has provided the folder to backup
			if os.Args[i] == "-f" && i+1 < len(os.Args) {
				// get absolute path of the folder
				currentDir, err := os.Getwd()
				if err != nil {
					fmt.Printf("Error getting current directory: %v\n", err)
					return
				}
				folder = fmt.Sprintf("%s/%s", currentDir, os.Args[i+1])
				continue
			}

			if os.Args[i] == "-v" {
				verbose = true
				continue;
			}

			if i+1 < len(os.Args) && os.Args[i] == "-d" {
				drive = fmt.Sprintf("%s:", os.Args[i+1])
				continue;	
			}
		}
		
		fmt.Println("Selected Drive: ", drive)
		fmt.Println("Selected Drive: ", folder)
		cmd.HandleBackup(verbose, folder, drive)
	case "commit":
		if len(os.Args) < 3{
			fmt.Println("Please provide commit message")
			os.Exit(0)
		}
		// Get the commit message
		msg := cmd.CommitAndStage(os.Args[2])
		fmt.Println(msg)
	case "stlng":
		if len(os.Args) == 2 {
			sysconfig.SetupGoNodePython()
		}else{
			fmt.Println("Invalid usage. Use 'setlang' to setup Go, Node and Python")
			os.Exit(0)
		}
	case "rfh":
		sysconfig.Refresh()
	case "weather":
		lat, _ := strconv.ParseFloat(os.Getenv("LATITUDE_S63_H149"), 64)
		long, _ := strconv.ParseFloat(os.Getenv("LONGITUDE_S63_H149"), 64)
		fmt.Printf("Location: %s\n", os.Getenv("S63_H149"))
		cmd.GetWeatherData(lat, long)
	case "help":
		fmt.Println(cmdFile)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println(cmdFile)
		os.Exit(1)
	}
}