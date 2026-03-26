package api

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/emicklei/go-restful/v3"
	"github.com/rs/zerolog"
	"github.com/xmtlzzz/vMusic/apps/music/impl"
)

type MusicApiHandler struct {
	log zerolog.Logger
}

func NewMusicApiHandler() *MusicApiHandler {
	return &MusicApiHandler{
		log: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
}

func (m *MusicApiHandler) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/music/local").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Doc("Local music APIs")

	ws.Route(ws.GET("/{path}").To(m.GetFromLocalPath).
		Doc("get from local path").
		Param(ws.PathParameter("path", "local directory path").DataType("string")).
		Writes(impl.LocalMusic{}.MusicList).
		Returns(200, "OK", impl.LocalMusic{}.MusicList).
		Returns(404, "music not found", nil))

	ws.Route(ws.POST("/upload").To(m.SaveFileToLocal).
		Consumes("multipart/form-data").
		Doc("upload music files to local path").
		Param(ws.MultiPartFormParameter("files", "music files to upload").DataType("file").AllowMultiple(true)).
		Returns(200, "OK", map[string]interface{}{}).
		Returns(400, "bad request", nil).
		Returns(500, "internal server error", nil))

	return ws
}

func (m *MusicApiHandler) GetFromLocalPath(request *restful.Request, response *restful.Response) {
	log.Println("request local path to find music list")
	lp := request.PathParameter("path")
	if lp == "" {
		lp = "C:\\Users\\Administrator\\Desktop\\code\\Go\\vMusic\\test"
	}
	lm := impl.NewMusicObj().GetFromLocalPath(lp)
	if lm.MusicList == nil {
		log.Printf("could't find any music in local path: %v", lp)
		response.WriteErrorString(404, "music not found")
		return
	}
	err := response.WriteEntity(lm.MusicList)
	if err != nil {
		log.Printf("http response error: %v", err)
		response.WriteErrorString(500, "internal server error")
		return
	}
}

func (m *MusicApiHandler) SaveFileToLocal(request *restful.Request, response *restful.Response) {
	err := request.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("failed to parse multipart form: %v", err)
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
			log.Printf("failed to open file %s: %v", fileHeader.Filename, err)
			continue
		}

		dst, err := os.Create(filepath.Join(savePath, fileHeader.Filename))
		if err != nil {
			log.Printf("failed to create file %s: %v", fileHeader.Filename, err)
			file.Close()
			continue
		}
		_, err = io.Copy(dst, file)
		file.Close()
		dst.Close()

		if err != nil {
			log.Printf("failed to save file %s: %v", fileHeader.Filename, err)
			continue
		}
		savedCount++
	}
	if savedCount == 0 {
		response.WriteErrorString(500, "all files failed to save")
		return
	}
	response.WriteEntity(map[string]interface{}{
		"message": "save success",
		"saved":   savedCount,
		"total":   len(files),
	})
}
