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
- Linux or macOS operating system.
- Sudo privileges for master user setup.
- [Specify dependencies, e.g., Python 3.8+, Node.js, or specific libraries if applicable].
- Network access for hosting the filesystem.
- Terminal access for configuration and user management.

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/Archivus.git
   cd Archivus
   ```

2. Build and release:
   ```bash
   make build    # Build both backend and frontend
   make package  # Create release archives in dist/packages
   ```

## Versioning
You can check the current version of the Archivus backend by running:
```bash
./archivus-v2/archivus-v2 -v
```

## Release Structure
The `make package` command generates archives for different platforms in the `dist/packages/` directory:
- `archivus-v2-<version>-linux-amd64.tar.gz`
- `archivus-v2-<version>-darwin-amd64.tar.gz`
- `archivus-v2-<version>-darwin-arm64.tar.gz`