package main

import (
	"log"
	"log/slog"
	"os"

	"golang-mcp-testing/internal/utils"
	"golang-mcp-testing/tools/config"
	"golang-mcp-testing/tools/dropbox"
	"golang-mcp-testing/tools/terminal"

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

	s.Tool("dropbox_list_dropbox_folder", "List all dropbox folders within a given path.",
		dropbox.HandleListDropboxFolder)

	s.Tool("dropbox_files_download", "Download a file at a provided path.",
		dropbox.HandleFilesDownload)

	s.Tool("terminal_write_file", "Write a file to the filesystem.",
		terminal.HandleWriteFile)

	// for testing - using the new generic handler utility
	err := utils.CallHandlerDirectly(logger, "HandleWriteFile",
		terminal.WriteFileArgs{Path: "/Users/bittelc/Desktop/file.txt", Content: "this content"},
		terminal.HandleWriteFile)
	if err != nil {
		log.Fatalf("direct call to handler failed: %v", err)
	}

	if err := s.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
