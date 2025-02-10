#!/bin/sh

# Install dependencies
os_arch=$(uname -a)
echo "Detected OS and Architecture: $os_arch"

echo "Installing TailwindCSS CLI"
# Determine the download URL based on OS and Architecture
npm install

echo "Installing Air"
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s;

echo "Installing Templ CLI"
go install github.com/a-h/templ/cmd/templ@latest;

echo "Installing go modules..."
go mod tidy
