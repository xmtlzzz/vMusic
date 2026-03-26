package main

import (
	"log"
	"net/http"

	_ "github.com/xmtlzzz/vMusic/apps"
)

func main() {

	log.Println("Server starting on :8080")
	log.Println("Swagger UI: http://localhost:8080/swagger/")
	log.Println("OpenAPI JSON: http://localhost:8080/apidocs.json/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
