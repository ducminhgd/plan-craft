#!/bin/bash

# Quick start script for Plan Craft
# This script sets up everything you need to run the application

set -e

echo "🚀 Plan Craft - Quick Start"
echo "============================"
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "   Download from: https://go.dev/dl/"
    exit 1
fi
echo "✅ Go $(go version | awk '{print $3}')"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 16 or higher."
    echo "   Download from: https://nodejs.org/"
    exit 1
fi
echo "✅ Node.js $(node --version)"

# Check npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm."
    exit 1
fi
echo "✅ npm $(npm --version)"

# Check Wails
if ! command -v wails &> /dev/null; then
    echo "⚠️  Wails CLI is not installed."
    echo "   Installing Wails CLI..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    
    # Check if installation was successful
    if ! command -v wails &> /dev/null; then
        echo "❌ Failed to install Wails CLI."
        echo "   Please ensure \$GOPATH/bin is in your PATH."
        echo "   Run: export PATH=\$PATH:\$(go env GOPATH)/bin"
        exit 1
    fi
fi
echo "✅ Wails $(wails version 2>&1 | head -n 1 || echo 'installed')"

echo ""
echo "📦 Installing dependencies..."

# Install Go dependencies
echo "  → Installing Go dependencies..."
go mod download
go mod tidy

# Install frontend dependencies
echo "  → Installing frontend dependencies..."
cd frontend
npm install
cd ..

echo ""
echo "✅ Setup complete!"
echo ""
echo "🎯 You can now run the application with:"
echo ""
echo "   make wails-dev      # Development mode with hot reload (recommended)"
echo "   make wails-build    # Build production binary"
echo ""
echo "Or use Wails directly:"
echo ""
echo "   wails dev           # Development mode"
echo "   wails build         # Production build"
echo ""
echo "📚 For more information, see docs/GETTING_STARTED.md"
echo ""

