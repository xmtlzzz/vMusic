package impl

import (
	"log"
	"os"
	"strings"
)

func GetFromLocalPath(localPath string) *LocalMusic {
	lm := NewLocalMusic()
	pE, err := os.ReadDir(localPath)
	if err != nil {
		log.Printf("读取本地路径失败，err: %v", err)
	}
	for _, entry := range pE {
		if entry.Name() == "README.md" {
			continue
		}
		// 文件名切分
		infoList := strings.Split(entry.Name(), "-")
		infoList[1] = strings.Trim(infoList[1], ".mp3")
		lm.MusicList = append(lm.MusicList, Objector{MusicName: infoList[0],
			Author: infoList[1],
			From:   "local",
			Path:   localPath})
	}
	return lm
}
