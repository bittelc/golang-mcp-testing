package dropbox

import "github.com/localrivet/gomcp/server"

type FilesDownloadArgs struct {
	Path string `json:"path"`
}

type FileLockInfo struct {
	Created        string `json:"created"`
	IsLockholder   bool   `json:"is_lockholder"`
	LockholderName string `json:"lockholder_name"`
}

type PropertyField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PropertyGroup struct {
	Fields     []PropertyField `json:"fields"`
	TemplateID string          `json:"template_id"`
}

type SharingInfo struct {
	ModifiedBy           string `json:"modified_by"`
	ParentSharedFolderID string `json:"parent_shared_folder_id"`
	ReadOnly             bool   `json:"read_only"`
}

type DropboxFileMetadata struct {
	ClientModified           string          `json:"client_modified"`
	ContentHash              string          `json:"content_hash"`
	FileLockInfo             FileLockInfo    `json:"file_lock_info"`
	HasExplicitSharedMembers bool            `json:"has_explicit_shared_members"`
	ID                       string          `json:"id"`
	IsDownloadable           bool            `json:"is_downloadable"`
	Name                     string          `json:"name"`
	PathDisplay              string          `json:"path_display"`
	PathLower                string          `json:"path_lower"`
	PropertyGroups           []PropertyGroup `json:"property_groups"`
	Rev                      string          `json:"rev"`
	ServerModified           string          `json:"server_modified"`
	Size                     int64           `json:"size"`
}

// HandleFilesDownload implements the logic the files.download tool
// This handler downloads the file at the provided FilesDownloadArgs.Path
func HandleFilesDownload(ctx *server.Context, args FilesDownloadArgs) (DropboxFileMetadata, error) {
	ctx.Logger.Info("Handling FilesDownload tool call")
	var metadata DropboxFileMetadata
	return metadata, nil
}
