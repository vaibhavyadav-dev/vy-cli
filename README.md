# vy-cli

A powerful command-line tool designed to automate repetitive tasks, simplify complex workflows, and boost productivity on Ubuntu.

## Features
> More features will be added in upcoming updates, stay tuned! 🚀

- **Ubuntu Backup Automation**  
  Easily back up system settings, configurations, and preferences to OneDrive with a single command:
  ```bash
  vy backup
  ```

- **Simplified Git Workflow**  
  Streamline your Git operations with a single command for `git add` and `git commit`:
  ```bash
  vy commit "Your commit message"
  ```

- **File Extraction**  
  Extract files of different types (e.g., `.zip`, `.tar`, `.7z`, `.rar`) using a single command:
  ```bash
  vy extract <FILE_NAME>
  ```

- **File and Folder Search**  
  Quickly search for files or folders with specific names, sizes, or types:
  ```bash
  vy find -n <NAME> -s <SIZE> --type <FILE_TYPE>
  ```

- **System Information**  
  View detailed system information like memory usage, storage, CPU details, and running processes in well-formatted tables:
  ```bash
  vy sysinfo
  ```

- **Weather Information**  
  Fetch current weather details, including AQI, sunrise, sunset, and more:
  ```bash
  vy weather
  ```

- **System Update and Upgrade**  
  Update and upgrade your system with `-y` already included:
  ```bash
  vy rfh
  ```

- **Install Essential Languages**  
  Install Go (v1.22.11), Python (v3.10.12), and Node.js (v22.13.1), skipping already installed versions:
  ```bash
  vy stlng
  ```

---

## Installation

### Prerequisites
- Ubuntu operating system
- Golang installed (version 1.20 or higher recommended)
- OneDrive account for backup functionality (support for Google Drive and more in future updates)

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/vaibhavyadav-dev/vy-cli.git
   ```

2. Navigate to the project directory:
   ```bash
   cd vy-cli
   ```

3. Execute the `install.sh` file:
   ```bash
   chmod +x install.sh
   ./install.sh
   ```

---

## Usage

```
Usage:

    vy <command> [arguments]

Commands:
    date              Show date and time
    backup            Back up system settings, configurations, and preferences to OneDrive
                      Usage:
                        vy backup [-v] [-f folder] [-d drive]
                      Arguments:
                        -v: Verbose mode
                        -f: Folder (absolute path) to back up
                        -d: Drive to back up to (currently supports folders only)

    commit            Stage and commit all project changes
                      Example:
                        vy commit "first commit" (must provide a message in double quotes)
    
    extract           Extract files of types `.zip`, `.tar`, `.7z`, `.rar`
                      Example:
                        vy extract <FILE_NAME>
    
    find              Find files or folders based on name, size, or type
                      Usage:
                        vy find -n <NAME> -s <SIZE> --type <FILE_TYPE>
                      Arguments:
                        -n: Name of the file or folder
                        -s: File size (default: 0 MB)
                        -h: Hard search (search all files and folders)
                        --type: Specify file type (currently limited support)
    
    sysinfo           Display system information (memory, storage, CPU, processes)
    
    weather           Fetch weather details (AQI, sunrise, sunset, etc.)
    
    rfh               Update and upgrade the system (`-y` is already included)
    
    stlng             Install Go (v1.22.11), Python (v3.10.12), Node.js (v22.13.1)
                      (skips already installed versions)
    
    help              Display help information

Arguments:
    -v                Verbose mode
```

---

## Author
Developed by [Vaibhav Yadav](https://www.linkedin.com/in/vaibhav-yadav-4397351b9/).

---
