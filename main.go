package main

import (
	"log"
	"log/slog"
	"os"

	"golang-mcp-testing/dxt"
	"golang-mcp-testing/tools/config"
	"golang-mcp-testing/tools/dropbox"
	"golang-mcp-testing/tools/terminal"

	"github.com/localrivet/gomcp/server"
)

func main() {
	// Parse command line arguments
	if len(os.Args) > 1 && os.Args[1] == "--create-dxt" {
		log.Println("Creating DXT manifest...")
		dxt.CreateDxt()
		log.Println("DXT manifest created")
		return
	}

	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create a new server
	// s := server.NewServer("GoCreate",
	s := server.NewServer("ColeMCPServer",
		server.WithLogger(logger),
	).AsStdio()

	// Register tools using the API
	// Configuration tools
	s.Tool("get_config", "Get the complete server configuration as JSON.",
		config.HandleGetConfig)

	// Terminal tools
	s.Tool("execute_command", "Execute a terminal command with timeout.",
		terminal.HandleExecuteCommand)

	s.Tool("read_output", "Read new output from a running terminal session.",
		terminal.HandleReadOutput)

	s.Tool("force_terminate", "Force terminate a running terminal session.",
		terminal.HandleForceTerminate)

	s.Tool("list_sessions", "List all active terminal sessions.",
		terminal.HandleListSessions)

	s.Tool("execute_in_terminal", "Execute a command in the terminal (client-side execution).",
		terminal.HandleExecuteInTerminal)

	s.Tool("list_dropbox_folders", "List all dropbox folders within a given path.",
		dropbox.HandleListDropboxFolders)

	if err := s.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
