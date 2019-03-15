package main

import (
	"log"
	"net/http"

	"github.com/KHs000/lxndr/pkg/httphandler"
)

func main() {
	httphandler.StartDatabase()
	httphandler.ExportRoutes()

	log.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
