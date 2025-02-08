package feature

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func FindFileFolder(fileType, fileName string, fileSize int, hardSearch bool) {
	fmt.Printf("\n\033[34m\033[1m" +
		"███████╗██╗██╗     ███████╗    ███████╗██╗███╗   ██╗██████╗ ███████╗██████╗ \n" +
		"██╔════╝██║██║     ██╔════╝    ██╔════╝██║████╗  ██║██╔══██╗██╔════╝██╔══██╗\n" +
		"█████╗  ██║██║     █████╗      █████╗  ██║██╔██╗ ██║██║  ██║█████╗  ██████╔╝\n" +
		"██╔══╝  ██║██║     ██╔══╝      ██╔══╝  ██║██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗\n" +
		"██║     ██║███████╗███████╗    ██║     ██ ██║ ╚████║██████╔╝███████╗██║  ██║\n" +
		"╚═╝     ╚═╝╚══════╝╚══════╝    ╚═╝     ╚═╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝\n\033[0m")

	if fileType != "" && fileName == "" {
		if _, err := searchByFileType(fileType, hardSearch); err != nil {
			log.Fatal("Error: ", err)
		}
		os.Exit(0)
	}

	err := filepath.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if !hardSearch {
			baseName := filepath.Base(path)
			for i := range excludeFilesFolder {
				if strings.HasSuffix(baseName, excludeFilesFolder[i]) {
					return filepath.SkipDir
				}
			}
		}

		if strings.Contains(strings.ToLower(filepath.Base(path)), strings.ToLower(fileName)) {
			fmt.Printf("\n\033[32m\033[1m\033[5mFOUND!\033[0m")
			fmt.Printf("\n\033[32mLocation: %s (%d MB)\033[0m \n\n", path, info.Size()/1024/1024)
			os.Exit(0)
		}

		if info.Size()/1024/1024 >= int64(fileSize) {
			if fileSize >= 5 {
				fmt.Printf("\033[33mSize Found: %s (%d bytes) (%d MB)\033[0m \n", path, info.Size(), info.Size()/1024/1024)
			}
		}

		return nil
	})
	if err == nil {
		fmt.Printf("\n\033[31m\033[1mNo File or Directory Found\033[0m\n")
		fmt.Println("Make sure you've typed the correct file or Folder Name")
	}
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
	}
}

func searchByFileType(fileType string, hardSearch bool) (bool, error) {
	fmt.Printf("\n\033[34m\033[1mSearching for files with extension: %s\033[0m click (ctrl + left mouse button) to open \n", fileType)
	err := filepath.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// Skip directories
		if !hardSearch && isExcluded(path) {
			return filepath.SkipDir
		}

		// fmt.Println("Checking: ", path)

		// Check if file extension matches
		IsExist, _ := isFileType(path, fileType)
		if IsExist {
			fmt.Printf("\033[32mLocation: file://%s (%d MB)\033[0m\n", path, info.Size()/1024/1024)
		}

		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func isFileType(filePath string, expectedType string) (bool, error) {
	magicBytes := map[string][]byte{
		"jpeg": []byte("\xFF\xD8\xFF"),
		"png":  []byte("\x89PNG"),
		"pdf":  []byte("%PDF"),
		"gif":  []byte("GIF87a"),
		"bmp":  []byte("BM"),
		"tiff": []byte("II*\x00"),
		"zip":  []byte("PK\x03\x04"),
		"rar":  []byte("Rar!\x1a\x07"),
		"7z":   []byte("7z\xBC\xAF\x27\x1C"),
		"exe":  []byte("MZ"),
		"elf":  []byte("\x7FELF"),
		"mp3":  []byte("ID3"),
		"mp4":  []byte("\x00\x00\x00\x20ftyp"),
		"avi":  []byte("RIFF"),
		"wav":  []byte("RIFF"),
		"ogg":  []byte("OggS"),
		"webm": []byte("\x1A\x45\xDF\xA3"),
		"mov":  []byte("\x00\x00\x00\x14ftyp"),
		"doc":  []byte("\xD0\xCF\x11\xE0\xA1\xB1\x1A\xE1"),
		"docx": []byte("PK\x03\x04"),
		"xls":  []byte("\xD0\xCF\x11\xE0\xA1\xB1\x1A\xE1"),
		"xlsx": []byte("PK\x03\x04"),
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return false, err
	}

	// Skip devices
	if info.Mode()&os.ModeDevice != 0 {
		return false, nil
	}

	// Read the first 512 bytes
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		return false, err
	}

	// Check for matching magic bytes
	if signature, exists := magicBytes[expectedType]; exists {
		return bytes.HasPrefix(buffer, signature), nil
	}
	return false, nil
}

// Skip directories
var excludeFilesFolder = []string{"node_modules", "tmp", "proc", "sys", "dev", "run",
	".snap", ".deb", ".vscode", ".idea", 
	".nvm", ".cache", ".npm", ".yarn",
	".config"}

// exclude the files or directories from the search
func isExcluded(path string) bool {
	for _, exclude := range excludeFilesFolder {
		if strings.Contains(path, exclude) {
			return true
		}
	}
	return false
}
