# vy-cli

A powerful command-line tool designed to automate repetitive tasks, simplify complex workflows, and boost productivity on Ubuntu.

## Features
> More features will be added in upcoming updates, stay tuned for it :)

- **Ubuntu Backup Automation**:
    Easily back up system settings, configurations, and preferences to OneDrive with a single command:
    ```bash
    vy backup
    ```

- **Simplified Git Workflow**:  
    Streamline your Git operations with a single command for git add and git commit:
    ```bash
    vy commit "Your commit message"
    ```
   > **MORE FEATURE WILL BE ADDED**

- **Productivity Enhancements**:
    Add shortcuts for tedious or repetitive commands to make your workflow more efficient.

## Installation

### Prerequisites
- Ubuntu operating system
- Golang installed (version 1.20 or higher recommended)
- OneDrive account for backup functionality (as of now onedrive, will add google drive and more in future updates)

### Steps
1. Clone the repository:
     ```bash
     git clone https://github.com/<your-username>/vy-cli.git
     ```

2. Navigate to the project directory:
     ```bash
     cd vy-cli
     ```

3. Execute the install.sh file:
     ```bash
     chmod +x install.sh
     ./install.sh
     ```

## Usage

```
Usage:

    vy <commands> [arguments]

The commands are:
    date              show date and time
    backup            backup all the settings, config, preferances to OneDrive
                      
                      vy-cli backup [-v] [-f folder] [-d drive]
                      [-v]: Verbose mode
                      [-f]: Folder, location absolute path, to backup
                      [-d]: Drive to backup to
                      
                      This will take name of folder, currently only folders are supported!
    
    commit            stage and commit ALL the changes of project, 
                      example:
                        vy commit "first commit"
                        (must add message with double inverted comma!)

    weather           fetch all the weather data, like AQI, sunrise, sunset etc

    rfh               update and upgrade the system (-y is already included in command)
    stlng             install Go(v1.22.11), Python(v3.10.12), Node(v22.13.1), skip if already installed
    help              displays help profile

arguments:
    -v                verbose mode
```

## Author
Developed by [Vaibhav Yadav](https://www.linkedin.com/in/vaibhav-yadav-4397351b9/).
