package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/vaibhavyadav-dev/vy-cli/cmd/feature"
	"github.com/vaibhavyadav-dev/vy-cli/cmd/config"

	// Load the .env file
	"github.com/joho/godotenv"
)

//go:embed cmd/cmd.txt
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
		feature.PrintRainbowGlowLargeText("Vaibhav Yadav")
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
		feature.ShowDate()
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
		if folder == "" {
			fmt.Println("Backing up the System Files...")
		}else{
			fmt.Printf("Backing up the folder %s \n", folder)
		}
		feature.HandleBackup(verbose, folder, drive)
	case "commit":
		if len(os.Args) < 3{
			fmt.Println("Please provide commit message")
			os.Exit(0)
		}
		// Get the commit message
		if err := feature.CommitAndStage(os.Args[2]); err != nil {
			log.Fatal(err)
		}
		fmt.Println("All changes has been Staged and Committed :)")
	case "stlng":
		if len(os.Args) == 2 {
			sysCmd.SetupGoNodePython()
		}else{
			fmt.Println("Invalid usage. Use 'setlang' to setup Go, Node and Python")
			os.Exit(0)
		}
	case "rfh":
		sysCmd.Refresh()
	case "weather":
		lat, _ := strconv.ParseFloat(os.Getenv("LATITUDE_S63_H149"), 64)
		long, _ := strconv.ParseFloat(os.Getenv("LONGITUDE_S63_H149"), 64)
		fmt.Printf("Location: %s\n", os.Getenv("S63_H149"))
		feature.GetWeatherData(lat, long)
	case "help":
		fmt.Println(cmdFile)
	case "extract":
		if len(os.Args) < 3 {
			fmt.Println("Please provide the file to extract")
			os.Exit(0)
		}
		filePath := os.Args[2]
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("File does not exist: %s\n", filePath)
			os.Exit(1)
		}
		feature.Extracter(filePath)
	case "sysinfo":
		feature.GetSysStatus()
	case "find":
		// vy find -n <filename> -s <filesize>
		fileName := ""
		fileSize := 0
		hardSearch := false
		fileType := ""
		
		for i:= 2; i < len(os.Args); i++ {
			switch expression := os.Args[i]; expression {	
			case "--type":
				if i+1 < len(os.Args) {
					fileType = os.Args[i+1]
					i++;
				}
				feature.FindFileFolder(fileType, fileName, fileSize, hardSearch)
			case "-n":
				if i+1 < len(os.Args) {
					fileName = os.Args[i+1]
					i++;
				}
			case "-s":
				if i+1 < len(os.Args) {
					fileSize, _ = strconv.Atoi(os.Args[i+1])
					i++;
				}
			case "-h":
				hardSearch = true
			default:	
				log.Fatalf("Invalid flag: %s", expression)
			}
		}
		feature.FindFileFolder("", fileName, fileSize, hardSearch)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println(cmdFile)
		os.Exit(1)
	}
}