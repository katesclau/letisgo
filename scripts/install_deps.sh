#!/bin/sh

# Install dependencies
os_arch=$(uname -a)
echo "Detected OS and Architecture: $os_arch"

echo "Installing TailwindCSS CLI"
# Determine the download URL based on OS and Architecture
npm install

echo "Installing Air"
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s;

echo "Installing Terraform"
case "$os_arch" in
  *"Darwin"*)
    echo "Detected macOS, using brew to install Terraform and tflocal"
    brew tap hashicorp/tap
    brew install hashicorp/tap/terraform
    brew install python3 pipx
    pipx install terraform-local
    ;;
  *"Ubuntu"*)
    echo "Detected Ubuntu, using apt to install Terraform and tflocal"
    wget -O - https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
    sudo apt update && sudo apt install -y terraform
    sudo apt-get install -y python3-pip
    pipx install terraform-local
    ;;
  *)
    echo "Non-macOS/Ubuntu system detected, skipping installation for Terraform and tflocal"
    ;;
esac
export PATH=$PATH:$HOME/.local/bin

echo "Installing Templ CLI"
go install github.com/a-h/templ/cmd/templ@latest;

echo "Installing go modules..."
go mod tidy
