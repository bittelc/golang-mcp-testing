package dropbox

type ListDropboxFileArgs struct {
	Path string `json:"path"`
}

type DropboxFiles []DropboxFile
type DropboxFile struct {
	Created     string `json:"created"`
	Description string `json:"description"`
	Destination string `json:"destination"`
	FileCount   int    `json:"file_count"`
	ID          string `json:"id"`
	IsOpen      bool   `json:"is_open"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

// func HandleListDropboxFiles(ctx *server.Context, args ListDropboxFileArgs) (DropboxFolders, error) {
// 	ctx.Logger.Info("Handling ListDropboxFiles tool call")
// }
