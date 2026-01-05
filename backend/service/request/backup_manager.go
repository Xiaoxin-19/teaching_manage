package requestx

type SetBackupLocalPathRequest struct {
	LocalPath string `json:"local_path" validate:"required,max=2048"`
}

type OpenSelectPathRequest struct {
	Title       string         `json:"title" validate:"max=256"`
	DefaultPath string         `json:"default_path" validate:"max=2048"`
	Filters     []FileFilteDTO `json:"filters"`
	IsPath      bool           `json:"is_path"`
}

type FileFilteDTO struct {
	DisplayName string `json:"display_name"` // Filter information EG: "Image Files (*.jpg, *.png)"
	Pattern     string `json:"pattern"`      // semicolon separated list of extensions, EG: ["*.jpg", "*.png"]
}

type SetWebDavConfigRequest struct {
	WebDavURL      string `json:"url" validate:"required,url,max=1024"`
	WebDavUserName string `json:"username" validate:"required,max=256"`
	WebDavPassword string `json:"password" validate:"required,max=256"`
}

type RestoreBackupRequest struct {
	Type       string `json:"type" validate:"required,oneof=local webdav"`
	BackupPath string `json:"backup_path" validate:"required,max=2048"`
}

type CreateBackupRequest struct {
	Type string `json:"type" validate:"required,oneof=local webdav"`
}
