package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/xmtlzzz/vMusic/apps/music/api"
)

func main() {
	ws := api.NewMusicApiHandler().WebService()
	restful.Add(ws)
	log.Println("Server starting on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
