# Archivus

Archivus is a tool to host your filesystem on a network, allowing secure access and management of files for multiple users on Linux and macOS systems. It supports user creation, access control (read/write permissions), and file operations like upload, download, and move.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Creating a Master User](#creating-a-master-user)
  - [Creating a New User](#creating-a-new-user)
  - [Managing User Access](#managing-user-access)
  - [File Operations](#file-operations)
- [Contributing](#contributing)
- [License](#license)

## Features
- Host your filesystem on a network for remote access (Linux/macOS).
- Create a master user with sudo privileges for administration.
- Add new users via terminal with customizable access.
- Grant users access to either the full filesystem or a specific directory.
- Toggle read/write permissions for users.
- Support for file operations: upload, download, and move files.

## Prerequisites
- **Linux or macOS** operating system.
- **Go 1.18+** (for building the backend).
- **Node.js 16+ and npm** (for building and running the frontend).
- **Make** (standard on macOS, install via `apt install make` on Linux).
- **Sudo privileges** for installing services.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/Archivus.git
   cd Archivus
   ```

2. **Build the project**:
   This will compile the Go backend and build the Next.js frontend into the `dist` directory.
   ```bash
   make build
   ```

3. **Run the installer**:
   The installer will copy the binaries and frontend files to `~/archivus-v2` and set up system services (systemd on Linux, launchd on macOS).
   ```bash
   sudo ./install.sh
   ```

## Usage

### Starting the Server
The server is automatically started as a service after installation. You can access the web interface at:
`http://localhost:3000`

### Creating a Master User
To manage the system, you first need to create a master user via the terminal:
```bash
~/archivus-v2/bin/archivus-v2 new-user
```
Follow the prompts to enter a username, password, 6-digit PIN, and email. Select `y` when asked if this is a master user.

**Note**: You will be prompted for your sudo password to verify permissions.

### Managing User Settings
You can toggle user permissions (like directory locking or write access) using the following command:
```bash
~/archivus-v2/bin/archivus-v2 user-settings
```
You will need the master PIN to perform these operations.

### File Operations
- **Upload**: Securely upload files via the web interface.
- **Download**: Access and download your files from any device on the network.
- **Move**: Organize your filesystem remotely.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

## License
MIT License