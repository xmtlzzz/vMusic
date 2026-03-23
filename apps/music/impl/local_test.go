package impl_test

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/xmtlzzz/vMusic/apps/music/api"
	"github.com/xmtlzzz/vMusic/apps/music/impl"
)

func TestGetFromLocalPath(t *testing.T) {
	res := impl.NewMusicObj().GetFromLocalPath("C:\\Users\\Administrator\\Desktop\\code\\Go\\vMusic\\test")
	if res == nil {
		t.Fatal("未搜到对应的music")
	}
	t.Log(res)
}

func TestSaveFileToLocal(t *testing.T) {
	// 创建测试目录
	testDir := "C:\\Users\\Administrator\\Desktop\\code\\Go\\vMusic\\test"
	os.MkdirAll(testDir, 0755)

	// 创建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加测试文件1
	part1, err := writer.CreateFormFile("files", "测试歌曲-测试作者.mp3")
	if err != nil {
		t.Fatal(err)
	}
	part1.Write([]byte("fake mp3 content 1"))

	// 添加测试文件2
	part2, err := writer.CreateFormFile("files", "歌曲2-作者2.mp3")
	if err != nil {
		t.Fatal(err)
	}
	part2.Write([]byte("fake mp3 content 2"))

	writer.Close()

	// 创建 HTTP 请求
	req := httptest.NewRequest("POST", "/local/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 创建响应记录器
	recorder := httptest.NewRecorder()

	// 创建 restful 容器并添加路由
	container := restful.NewContainer()
	container.Add(api.NewMusicApiHandler().WebService())
	container.ServeHTTP(recorder, req)

	// 检查响应状态
	if recorder.Code != 200 {
		t.Errorf("期望状态码 200, 得到 %d, 响应: %s", recorder.Code, recorder.Body.String())
	}

	// 验证文件是否保存
	file1Path := filepath.Join(testDir, "测试歌曲-测试作者.mp3")
	if _, err := os.Stat(file1Path); os.IsNotExist(err) {
		t.Error("文件1未保存")
	} else {
		// 清理测试文件
		os.Remove(file1Path)
	}

	file2Path := filepath.Join(testDir, "歌曲2-作者2.mp3")
	if _, err := os.Stat(file2Path); os.IsNotExist(err) {
		t.Error("文件2未保存")
	} else {
		// 清理测试文件
		os.Remove(file2Path)
	}

	t.Logf("响应: %s", recorder.Body.String())
}
