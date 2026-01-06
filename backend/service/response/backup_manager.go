package responsex

type WebDavConfigResponse struct {
	WebDavURL        string `json:"url"`
	WebDavUserName   string `json:"username"`
	WebDavPassword   string `json:"password"`
	LastCloudBackup  int64  `json:"last_cloud_backup"`
	EnableAutoBackup bool   `json:"enable_auto_backup"`
}

type WebDavBackupItem struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`
	Path    string `json:"path"`
}
