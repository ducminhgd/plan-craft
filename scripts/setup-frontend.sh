#!/bin/bash

# Setup and build frontend for Plan Craft

set -e

echo "🚀 Setting up frontend..."

# Check if node is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js first."
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm first."
    exit 1
fi

echo "📦 Installing frontend dependencies..."
cd frontend
npm install

echo "🔨 Building frontend..."
npm run build

echo "✅ Frontend setup complete!"
echo ""
echo "You can now run the application with:"
echo "  wails dev    (for development with hot reload)"
echo "  wails build  (for production build)"

