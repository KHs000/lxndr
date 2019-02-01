package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/KHs000/lxndr/pkg/mongo"
)

type configs struct {
	ConnStr string `json:"connectionString"`
}

func main() {
	startDatabase()
	exportRoutes()

	log.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func startDatabase() {
	configF, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configF.Close()

	bytes, err := ioutil.ReadAll(configF)
	if err != nil {
		log.Fatal(err)
	}

	var config configs
	json.Unmarshal(bytes, &config)
	mongo.Connect(config.ConnStr)
}
