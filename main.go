package main

import (
	"log"
	"log/slog"
	"os"

	"golang-mcp-testing/tools/config"
	"golang-mcp-testing/tools/dropbox"

	"github.com/davecgh/go-spew/spew"
	"github.com/localrivet/gomcp/server"
)

func main() {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Create a new server
	// s := server.NewServer("GoCreate",
	s := server.NewServer("ColeMCPServer",
		server.WithLogger(logger),
	).AsStdio()
	serverContext := &server.Context{
		Logger: logger,
	}

	// Register tools using the API
	// Configuration tools
	s.Tool("get_config", "Get the complete server configuration as JSON.",
		config.HandleGetConfig)
	folders, err := dropbox.HandleListDropboxFolders(serverContext, dropbox.ListDropboxFoldersArgs{Path: ""})
	if err != nil {
		log.Fatalf("Failed to list dropbox folders", "error", err)
	}
	spew.Dump(folders)
	// s.Tool("list_dropbox_folders", "List all dropbox folders within a given path.",
	// 	dropbox.HandleListDropboxFolders)

	if err := s.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
