### 
```bash
sudo apt-get update
sudo apt-get install git
git clone https://github.com/Samar2170/archivus-v2.git
sudo apt install curl make

curl -LO https://dl.google.com
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.26.0.linux-amd64.tar.gz
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile

curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash
source .bashrc
nvm install node

sudo apt update && sudo apt install libatomic1
