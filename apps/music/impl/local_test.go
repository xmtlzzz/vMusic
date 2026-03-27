package impl_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/xmtlzzz/vMusic/apps/music/api"
	"github.com/xmtlzzz/vMusic/apps/music/impl"
)

func TestGetFromLocalPath(t *testing.T) {
	libraryDir := t.TempDir()
	mustWriteFile(t, filepath.Join(libraryDir, "冬天的秘密-周传雄.mp3"), []byte("fake mp3"))
	mustWriteFile(t, filepath.Join(libraryDir, "纯音乐.flac"), []byte("fake flac"))
	mustWriteFile(t, filepath.Join(libraryDir, "README.md"), []byte("ignore me"))

	res := impl.NewMusicObj().GetFromLocalPath(libraryDir)
	if res == nil {
		t.Fatal("expected local music result")
	}
	if len(res.MusicList) != 2 {
		t.Fatalf("expected 2 tracks, got %d", len(res.MusicList))
	}
	if res.MusicList[0].AudioURL == "" {
		t.Fatal("expected audio url to be generated")
	}
}

func TestSaveFileToLocal(t *testing.T) {
	libraryDir := t.TempDir()
	previousDir := os.Getenv("VMUSIC_LIBRARY_DIR")
	t.Cleanup(func() {
		if previousDir == "" {
			os.Unsetenv("VMUSIC_LIBRARY_DIR")
			return
		}
		os.Setenv("VMUSIC_LIBRARY_DIR", previousDir)
	})
	os.Setenv("VMUSIC_LIBRARY_DIR", libraryDir)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "测试歌曲-测试作者.mp3")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = part.Write([]byte("fake mp3 content")); err != nil {
		t.Fatal(err)
	}
	if err = writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/music/local/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()
	container := restful.NewContainer()
	container.Add(api.NewMusicApiHandler().WebService())
	container.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", recorder.Code, recorder.Body.String())
	}

	var response impl.UploadResult
	if err = json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Saved != 1 {
		t.Fatalf("expected 1 saved file, got %d", response.Saved)
	}

	files, err := os.ReadDir(libraryDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file in library, got %d", len(files))
	}

	listReq := httptest.NewRequest(http.MethodGet, "/music/local/list", nil)
	listRecorder := httptest.NewRecorder()
	container.ServeHTTP(listRecorder, listReq)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d, body=%s", listRecorder.Code, listRecorder.Body.String())
	}

	var listResponse impl.LocalMusic
	if err = json.Unmarshal(listRecorder.Body.Bytes(), &listResponse); err != nil {
		t.Fatal(err)
	}
	if len(listResponse.MusicList) != 1 {
		t.Fatalf("expected 1 track in list, got %d", len(listResponse.MusicList))
	}
}

func mustWriteFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
}
