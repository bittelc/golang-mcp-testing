package main

import (
	"log"
	"log/slog"
	"os"

	"golang-mcp-testing/tools/config"
	"golang-mcp-testing/tools/dropbox"

	"github.com/localrivet/gomcp/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	s := server.NewServer("ColeMCPServer",
		server.WithLogger(logger),
	).AsStdio()

	s.Tool("get_config", "Get the complete server configuration as JSON.",
		config.HandleGetConfig)

	dropbox.HelperCallHandlerDirectly(logger) // for testing if needed
	// s.Tool("list_dropbox_folders", "List all dropbox folders within a given path.",
	// dropbox.HandleListDropboxFolders)

	if err := s.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
