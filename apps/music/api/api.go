package api

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/emicklei/go-restful/v3"
	"github.com/xmtlzzz/vMusic/apps/music/impl"
)

func WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/local").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{path}").To(GetFromLocalPath).
		Doc("get from local path").
		Writes(impl.LocalMusic{}.MusicList).
		Returns(200, "OK", impl.LocalMusic{}.MusicList).
		Returns(404, "music not found", nil))

	ws.Route(ws.POST("/upload").To(SaveFileToLocal).
		Consumes("multipart/form-data").
		Doc("upload music files to local path").
		Returns(200, "OK", map[string]interface{}{}).
		Returns(400, "bad request", nil).
		Returns(500, "internal server error", nil))

	return ws
}

func GetFromLocalPath(request *restful.Request, response *restful.Response) {
	log.Println("request local path to find music list")
	lp := request.PathParameter("path")
	if lp == "" {
		lp = "C:\\Users\\Administrator\\Desktop\\code\\Go\\vMusic\\test"
	}
	lm := impl.GetFromLocalPath(lp)
	if lm.MusicList == nil {
		log.Printf("could't find any music in local path: %v", lp)
		response.WriteErrorString(404, "music not found")
		return
	}
	err := response.WriteEntity(lm.MusicList)
	if err != nil {
		log.Printf("http返回错误: %v", err)
		response.WriteErrorString(500, "internal server error")
		return
	}
}

func SaveFileToLocal(request *restful.Request, response *restful.Response) {
	// 解析表单，最大32MB
	err := request.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("解析表单失败: %v", err)
		response.WriteErrorString(400, "failed to parse form")
		return
	}

	files := request.Request.MultipartForm.File["files"]
	if len(files) == 0 {
		response.WriteErrorString(400, "no files uploaded")
		return
	}

	savePath := "C:\\Users\\Administrator\\Desktop\\code\\Go\\vMusic\\test"
	savedCount := 0

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("打开文件失败 %s: %v", fileHeader.Filename, err)
			continue
		}

		dst, err := os.Create(filepath.Join(savePath, fileHeader.Filename))
		if err != nil {
			log.Printf("创建文件失败 %s: %v", fileHeader.Filename, err)
			file.Close()
			continue
		}
		_, err = io.Copy(dst, file)
		file.Close()
		dst.Close()

		if err != nil {
			log.Printf("保存文件失败 %s: %v", fileHeader.Filename, err)
			continue
		}
		savedCount++
	}
	if savedCount == 0 {
		response.WriteErrorString(500, "all files failed to save")
		return
	}
	response.WriteEntity(map[string]interface{}{
		"message": "保存成功",
		"saved":   savedCount,
		"total":   len(files),
	})
}
