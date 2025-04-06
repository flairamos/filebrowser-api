package filebrowserapi

type FileBrowserConfig struct {
	BaseURL  string // File Browser 服务器地址   结尾不需要 /
	Username string // 用户名
	Password string // 密码
	FilePath string // 目标路径   以/ 开头绝对路径相当于 /
	Override bool   // 是否覆盖
}

func NewConfig(url, username, password, filepath string) FileBrowserConfig {
	return FileBrowserConfig{
		BaseURL:  url,
		Username: username,
		Password: password,
		FilePath: filepath,
	}
}

// DownloadOptions 表示下载选项
type DownloadOptions struct {
	BaseURL    string // FileBrowser服务器地址
	FilePath   string // 要下载的文件路径
	OutputPath string // 输出路径
	AuthToken  string // 认证令牌
	Format     string // 压缩格式 (zip, tar, targz, tarbz2, tarxz, tarlz4, tarsz)
	Inline     bool   // 是否内联显示
}

// FileInfo 表示文件信息
type FileInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	IsDir     bool   `json:"isDir"`
	ModTime   string `json:"modTime"`
	Extension string `json:"extension"`
	Type      string `json:"type"`
}
