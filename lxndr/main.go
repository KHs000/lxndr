package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongo"
	"github.com/KHs000/lxndr/pkg/rndtoken"

	"github.com/mongodb/mongo-go-driver/bson"
)

type configs struct {
	ConnStr string `json:"connectionString"`
}

func main() {
	tkn, hash := rndtoken.SendToken("felipe.carbone@dito.com.br")

	fmt.Println("Token")
	fmt.Println(tkn)
	fmt.Println("Hash")
	fmt.Println(hash)

	identifier.IdentityCheck("felipe.carbone@dito.com.br", hash)

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

	coll := mongo.Collection{Database: "lxndr", CollName: "lxndr-quest"}
	res := mongo.Find(mongo.Conn, coll, bson.M{})

	fmt.Println("Consulta de teste a base do Mongo:")

	for res.Next(mongo.Conn.Ctx) {
		var result bson.M

		err := res.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result["test"])
	}
}
