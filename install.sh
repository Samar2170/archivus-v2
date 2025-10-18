#!/bin/bash

set -e

OS=$(uname -s)
ARCH=$(uname -m)

USER=$(whoami)
GROUP=$(id -gn)

PROJECT_DIR=$(pwd)
PROJECT_NAME=archivus-v2
INSTALL_DIR="$HOME/$PROJECT_NAME"
BIN_DIR="$INSTALL_DIR/bin"
SERVICE_DIR="/etc/systemd/system"
LAUNCHD_DIR="$HOME/Library/LaunchAgents"

sudo rm "$BIN_DIR/$PROJECT_NAME"
sudo rm -rf "$INSTALL_DIR/frontend"

sudo mkdir -p "$INSTALL_DIR"
sudo mkdir -p "$INSTALL_DIR/frontend"
sudo mkdir -p "$BIN_DIR"

echo "Installing $PROJECT_NAME..."

# Install dependencies
if [ "$OS" = "Linux" ]; then
    # Install Node.js (example for Ubuntu/Debian)
    if ! command -v node >/dev/null 2>&1; then
        echo "Installing Node.js..."
        sudo apt-get update
        sudo apt-get install -y nodejs npm
    fi
elif [ "$OS" = "Darwin" ]; then
    # Install Node.js (example for macOS using Homebrew)
    if ! command -v node >/dev/null 2>&1; then
        echo "Installing Node.js..."
        brew install node
    fi
fi


if [ "$OS" = "Linux" ]; then
    sudo cp "dist/bin/linux_amd64/$PROJECT_NAME" "$BIN_DIR/$PROJECT_NAME"
    sudo cp dist/bin/linux_amd64/config.prod.yaml "$BIN_DIR/config.prod.yaml"
elif [ "$OS" = "Darwin" ]; then
    sudo cp "dist/bin/macos_arm/$PROJECT_NAME" "$BIN_DIR/$PROJECT_NAME" 
fi

sudo chmod +x "$BIN_DIR/$PROJECT_NAME"


# Copy frontend files
cp -r dist/frontend/* "$INSTALL_DIR/frontend/"
cp -r dist/frontend/.next "$INSTALL_DIR/frontend/.next"
cd "$INSTALL_DIR/frontend"
npm install --production

cd "$HOME"
sudo chmod -R 777 "$INSTALL_DIR"

cd "$PROJECT_DIR"




sudo bash -c "cat > $SERVICE_DIR/archivus_v2.service" <<EOF
[Unit]
Description=archivus-v2 Go Backend
After=network.target

[Service]
ExecStart=$HOME/archivus-v2/bin/archivus-v2 server -m prod
WorkingDirectory=$HOME/archivus-v2
Restart=always
User=$USER
Group=$GROUP

[Install]
WantedBy=multi-user.target
EOF


sudo bash -c "cat > $SERVICE_DIR/archivus_client.service" <<EOF
[Unit]
Description=archivus-client Next.js Frontend
After=network.target

[Service]
WorkingDirectory=$HOME/archivus-v2/frontend
ExecStart=/usr/bin/npm start
Restart=always
User=$USER
Group=$GROUP

[Install]
WantedBy=multi-user.target
EOF





# Set up systemd services (Linux)
if [ "$OS" = "Linux" ]; then
    sudo systemctl daemon-reload
    sudo systemctl enable archivus_v2.service
    sudo systemctl enable archivus_client.service
    sudo systemctl start archivus_v2.service
    sudo systemctl start archivus_client.service
    echo "Systemd services installed and started."
fi


# Set up launchd (macOS)
if [ "$OS" = "Darwin" ]; then
    cp launchd/com.myapp.backend.plist "$LAUNCHD_DIR/"
    launchctl load "$LAUNCHD_DIR/com.myapp.backend.plist"
    echo "Launchd service installed and started."
    # Note: Next.js can be run manually or via another launchd service
fi

echo "Installation complete! Access your app at http://localhost:3000"