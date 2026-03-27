package impl

// 通用的类型信息
type Objector struct {
	// 音乐唯一标识
	ID string `json:"id"`
	// 名称
	MusicName string `json:"music_name" toml:"music_name"`
	// 作者
	Author string `json:"author" toml:"author"`
	// 音乐的来源
	From string `json:"from" toml:"from"`
	// 是否为喜欢的音乐
	IsLike bool `json:"is_like" toml:"is_like"`
	// 是否为收藏
	IsCollection bool `json:"is_collection" toml:"is_collection"`
	// 获取到该音乐的路径
	Path string `json:"path" toml:"path"`
	// 原始文件名
	FileName string `json:"file_name"`
	// 音频格式
	Format string `json:"format"`
	// 前端可直接播放的音频地址
	AudioURL string `json:"audio_url"`
}

type LocalMusic struct {
	LibraryPath string     `json:"library_path"`
	MusicList   []Objector `json:"music_list"`
}

type UploadResult struct {
	Message   string     `json:"message"`
	Saved     int        `json:"saved"`
	Total     int        `json:"total"`
	MusicList []Objector `json:"music_list"`
}

type RemoteMusic struct {
	// 该音乐列表是通过哪个远端获取的（网易云、qq等）
	MusicFrom string `json:"music_from" toml:"music_from"`
	MusicList []Objector
}

func NewLocalMusic() *LocalMusic {
	return &LocalMusic{}
}

func NewRemoteMusic() *RemoteMusic {
	return &RemoteMusic{}
}
