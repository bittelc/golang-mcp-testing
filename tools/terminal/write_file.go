package terminal

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/localrivet/gomcp/server"
)

// WriteFileArgs defines the arguments for the write_file tool.
type WriteFileArgs struct {
	Path    string `json:"path" description:"The path of the file to write to." required:"true"`
	Content string `json:"content" description:"The content to write to the file." required:"true"`
}

// HandleWriteFile implements the write_file tool using the new API
func HandleWriteFile(ctx *server.Context, args WriteFileArgs) (string, error) {
	ctx.Logger.Info("Handling write_file tool call")

	// Expand the path to handle ~ and relative paths
	expandedPath, err := expandPath(args.Path)
	if err != nil {
		ctx.Logger.Info("Error expanding path", "path", args.Path, "error", err)
		return "Error expanding path", err
	}

	// Check if file exists
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		ctx.Logger.Info("File does not exist, will create it", "path", expandedPath)
	} else {
		ctx.Logger.Info("File exists, will overwrite it", "path", expandedPath)
	}

	// Write the content to the file. 0644 is a common permission for files.
	err = os.WriteFile(expandedPath, []byte(args.Content), 0644)
	if err != nil {
		ctx.Logger.Info("Error writing file", "path", expandedPath, "error", err)
		return "Error writing file", err
	}

	return "File written successfully.", nil
}

// expandPath expands ~ to home directory and converts to absolute path
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr.HomeDir, path[2:])
	}
	return filepath.Abs(path)
}
