package responsex

type WebDavConfigResponse struct {
	WebDavURL      string `json:"url"`
	WebDavUserName string `json:"username"`
	WebDavPassword string `json:"password"`
}

type WebDavBackupItem struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`
	Path    string `json:"path"`
}
