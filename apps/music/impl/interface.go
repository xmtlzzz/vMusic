package impl

type MusicHandler interface {
	// 读取本地的music
	GetFromLocalPath(localPath string) *LocalMusic
}
