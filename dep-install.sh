#!/usr/bin/env bash
set -e

OS="$(uname -s)"
echo "Sift — dependency installer"
echo "OS: $OS"
echo ""

# ── macOS ──
if [ "$OS" = "Darwin" ]; then
    # Homebrew
    if ! command -v brew &>/dev/null; then
        echo "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    fi

    echo "Installing system packages..."
    brew install go node python3

    # pip packages
    pip3 install --upgrade pymupdf

    # Wails3
    go install github.com/wailsapp/wails/v3/cmd/wails3@latest

    echo ""
    echo "✓ macOS dependencies installed."
    echo "  Add ~/go/bin to your PATH if not already:"
    echo "  export PATH=\"\$HOME/go/bin:\$PATH\""
    echo ""
    echo "  Install Hermes: https://github.com/nousresearch/hermes-agent"
    echo "  Get SerpAPI key: https://serpapi.com"

# ── Linux ──
elif [ "$OS" = "Linux" ]; then
    # Detect package manager
    if command -v apt-get &>/dev/null; then
        PKG="apt-get"
        sudo apt-get update
        sudo apt-get install -y golang-go nodejs npm python3 python3-pip
    elif command -v dnf &>/dev/null; then
        PKG="dnf"
        sudo dnf install -y golang nodejs npm python3 python3-pip
    elif command -v pacman &>/dev/null; then
        PKG="pacman"
        sudo pacman -S --noconfirm go nodejs npm python python-pip
    elif command -v apk &>/dev/null; then
        PKG="apk"
        sudo apk add go nodejs npm python3 py3-pip
    else
        echo "Unsupported package manager. Install manually: go, nodejs, npm, python3, pip"
        exit 1
    fi

    # pip packages
    pip3 install --upgrade pymupdf

    # Wails3
    go install github.com/wailsapp/wails/v3/cmd/wails3@latest

    # Linux-specific Wails deps
    if command -v apt-get &>/dev/null; then
        sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.1-dev
    elif command -v dnf &>/dev/null; then
        sudo dnf install -y gtk3-devel webkit2gtk4.1-devel
    elif command -v pacman &>/dev/null; then
        sudo pacman -S --noconfirm gtk3 webkit2gtk-4.1
    fi

    echo ""
    echo "✓ Linux dependencies installed."
    echo "  Add ~/go/bin to your PATH if not already:"
    echo "  export PATH=\"\$HOME/go/bin:\$PATH\""
    echo ""
    echo "  Install Hermes: https://github.com/nousresearch/hermes-agent"
    echo "  Get SerpAPI key: https://serpapi.com"

else
    echo "Unsupported OS: $OS"
    echo "Install manually:"
    echo "  Go 1.25+     : https://go.dev/dl/"
    echo "  Node.js 22+  : https://nodejs.org/"
    echo "  Python 3     : https://python.org/"
    echo "  pymupdf       : pip install pymupdf"
    echo "  Wails3 CLI   : go install github.com/wailsapp/wails/v3/cmd/wails3@latest"
    echo "  Hermes       : https://github.com/nousresearch/hermes-agent"
    exit 1
fi
