#!/bin/sh

# Install dependencies
os_arch=$(uname -a)
echo "Detected OS and Architecture: $os_arch"

echo "Installing TailwindCSS CLI"
# Determine the download URL based on OS and Architecture
twpath="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss"
case "$os_arch" in
  *Linux*x86_64*)
    url="$twpath-linux-x64"
    ;;
  *Linux*arm64*)
    url="$twpath-linux-arm64"
    ;;
  *Darwin*arm64*)
    url="$twpath-macos-arm64"
    ;;
  *Darwin*x86_64*)
    url="$twpath-macos-x64"
    ;;
  *)
    echo "Unsupported OS or Architecture"
    exit 1
    ;;
esac
echo "Downloading from: $url"

# Download the tailwindcss CLI executable
curl -L $url -o ./bin/tailwindcss

# Make the tailwindcss CLI executable
chmod +x ./bin/tailwindcss

echo "Installing Air"
curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s;

echo "Installing Templ CLI"
go install github.com/a-h/templ/cmd/templ@latest;

echo "Installing go modules..."
go mod tidy
