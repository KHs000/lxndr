package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/KHs000/lxndr/pkg/mongo"
)

type configs struct {
	ConnStr string `json:"connectionString"`
}

func main() {
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

	createUser("daniel.silveira@dito.com.br")
}
