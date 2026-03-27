package api

import (
	"log"
	"net/http"
	"os"

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

	ws.Route(ws.GET("/list").To(m.ListLocalMusic).
		Doc("get local music list from configured library").
		Writes(impl.LocalMusic{}).
		Returns(200, "OK", impl.LocalMusic{}).
		Returns(500, "internal server error", nil))

	ws.Route(ws.POST("/upload").To(m.SaveFileToLocal).
		Consumes("multipart/form-data").
		Doc("upload music files to local path").
		Param(ws.MultiPartFormParameter("files", "music files to upload").DataType("file").AllowMultiple(true)).
		Returns(200, "OK", impl.UploadResult{}).
		Returns(400, "bad request", nil).
		Returns(500, "internal server error", nil))

	ws.Route(ws.GET("/file/{filename}").To(m.StreamLocalFile).
		Consumes("*/*").
		Produces("audio/mpeg", "audio/wav", "audio/flac", "audio/mp4", "audio/ogg").
		Doc("stream a local music file by filename").
		Param(ws.PathParameter("filename", "music file name").DataType("string")).
		Returns(200, "OK", nil).
		Returns(404, "music not found", nil))

	return ws
}

func (m *MusicApiHandler) ListLocalMusic(request *restful.Request, response *restful.Response) {
	log.Println("request configured local music library")

	lm := impl.NewMusicObj().GetFromLocalPath("")
	if err := response.WriteEntity(lm); err != nil {
		log.Printf("http response error: %v", err)
		response.WriteErrorString(500, "internal server error")
		return
	}
}

func (m *MusicApiHandler) StreamLocalFile(request *restful.Request, response *restful.Response) {
	fileName := request.PathParameter("filename")
	fullPath, err := impl.NewMusicObj().ResolveTrackPath(fileName)
	if err != nil {
		log.Printf("music file not found: %v", err)
		response.WriteErrorString(404, "music not found")
		return
	}

	http.ServeFile(response.ResponseWriter, request.Request, fullPath)
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

	result, err := impl.NewMusicObj().SaveUploadedFiles(files)
	if err != nil {
		log.Printf("save upload files failed: %v", err)
		response.WriteErrorString(500, "all files failed to save")
		return
	}

	if err = response.WriteEntity(result); err != nil {
		log.Printf("write upload response failed: %v", err)
		response.WriteErrorString(500, "internal server error")
	}
}
