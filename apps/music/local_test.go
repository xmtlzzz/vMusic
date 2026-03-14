package music_test

import (
	"testing"
	"vmusic/apps/music"
)

func TestGetFromLocalPath(t *testing.T) {
	res := music.GetFromLocalPath("C:\\Users\\Administrator\\Desktop\\code\\Go\\vMusic\\test")
	if res == nil {
		t.Fatal("未搜到对应的music")
	}
	t.Log(res)
}
