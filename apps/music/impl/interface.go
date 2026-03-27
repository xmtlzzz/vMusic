package impl

type MusicHandler interface {
	// 读取本地的music
	GetFromLocalPath(localPath string) *LocalMusic
	// 获取当前音乐库目录
	LibraryDir() string
	// 通过文件名查找音乐
	GetTrackByFileName(fileName string) (*Objector, error)
}
