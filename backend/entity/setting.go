package entity

type WebDavSetting struct {
	// 若需新增配置，只需在此处添加字段，并在 json tag 中定义 key 即可
	WebDavURL       string
	WebDavUserName  string
	WebDavPassword  string
	WebDavBaseDir   string
	LastCloudBackup string
}

type DataBackupSetting struct {
	WebDavSetting
	LocalBackupDirPath string
}
