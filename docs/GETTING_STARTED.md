# Getting Started with Plan Craft

This guide will help you set up and run Plan Craft on your local machine.

## Prerequisites

Before you begin, ensure you have the following installed:

1. **Go 1.21 or higher**
   ```bash
   go version
   ```

2. **Node.js 16 or higher** (for frontend)
   ```bash
   node --version
   npm --version
   ```

3. **Wails CLI** (for desktop app development)
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

4. **Platform-specific dependencies**:
   - **Linux**: `gcc`, `pkg-config`, `libgtk-3-dev`, `libwebkit2gtk-4.0-dev`
     ```bash
     # Ubuntu/Debian
     sudo apt install gcc pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev
     
     # Fedora
     sudo dnf install gcc pkg-config gtk3-devel webkit2gtk3-devel
     ```
   
   - **macOS**: Xcode Command Line Tools
     ```bash
     xcode-select --install
     ```
   
   - **Windows**: No additional dependencies required

## Quick Start

### Option 1: Development Mode (Recommended)

This is the easiest way to get started. It runs both the backend and frontend with hot reload:

```bash
# 1. Clone the repository
git clone git@github.com:ducminhgd/plan-craft.git
cd plan-craft

# 2. Install Go dependencies
make deps

# 3. Install frontend dependencies
make frontend-install

# 4. Run in development mode
make wails-dev
```

The application will open automatically with hot reload enabled. Any changes to Go or frontend code will trigger a rebuild.

### Option 2: Manual Setup

If you prefer more control:

```bash
# 1. Install Go dependencies
go mod download

# 2. Install frontend dependencies
cd frontend
npm install

# 3. Build frontend
npm run build
cd ..

# 4. Run the application
wails dev
```

## Available Make Commands

```bash
# Frontend commands
make frontend-install    # Install frontend dependencies
make frontend-build      # Build frontend for production
make frontend-dev        # Run frontend dev server only

# Wails commands
make wails-dev          # Run in development mode (hot reload)
make wails-build        # Build production binary

# Go commands
make deps               # Install Go dependencies
make test               # Run all tests
make fmt                # Format Go code
make lint               # Run linter (requires golangci-lint)

# Database commands
make db-clean           # Remove database files

# Cleanup
make clean              # Remove all build artifacts
```

## Project Structure

```
plan-craft/
├── frontend/           # React frontend
│   ├── src/           # React source code
│   ├── dist/          # Built frontend (generated)
│   └── package.json   # Frontend dependencies
├── internal/          # Internal Go packages
│   ├── entities/      # Domain entities
│   ├── repositories/  # Data access layer
│   └── usecases/      # Business logic
├── cmd/               # Application entry points
├── config/            # Configuration
├── migrations/        # Database migrations
├── main.go            # Wails application entry
└── app.go             # Application logic
```

## Development Workflow

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test -v ./internal/entities/...
```

### Building for Production

```bash
# Build for your current platform
make wails-build

# The binary will be in build/bin/
```

### Database Management

The application uses SQLite by default. The database file is created automatically at `data/plancraft.db`.

```bash
# Remove database (fresh start)
make db-clean

# Database migrations run automatically on startup
```

## Troubleshooting

### Error: "index.html: file does not exist"

This means the frontend hasn't been built yet. Run:

```bash
make frontend-build
```

Or use development mode which handles this automatically:

```bash
make wails-dev
```

### Error: "wails: command not found"

Install the Wails CLI:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Make sure `$GOPATH/bin` is in your PATH.

### Frontend dependencies not installing

Make sure you have Node.js and npm installed:

```bash
node --version  # Should be 16+
npm --version
```

### Build errors on Linux

Install the required GTK and WebKit dependencies:

```bash
# Ubuntu/Debian
sudo apt install gcc pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev
```

## Next Steps

1. **Explore the code**: Start with `main.go` and `app.go`
2. **Read the documentation**: Check `docs/` for detailed guides
3. **Run tests**: `make test` to ensure everything works
4. **Make changes**: Edit code and see hot reload in action
5. **Build**: `make wails-build` when ready to create a binary

## Additional Resources

- [Wails Documentation](https://wails.io/docs/introduction)
- [React Documentation](https://react.dev/)
- [GORM Documentation](https://gorm.io/docs/)
- [Project README](../README.md)

## Getting Help

- **Issues**: https://github.com/ducminhgd/plan-craft/issues
- **Discussions**: https://github.com/ducminhgd/plan-craft/discussions

