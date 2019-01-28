package main

import (
	"fmt"
	"log"

	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongoDB"
	"github.com/KHs000/lxndr/pkg/rndtoken"

	"github.com/mongodb/mongo-go-driver/bson"
)

func main() {
	tkn, hash := rndtoken.SendToken("felipe.carbone@dito.com.br")

	fmt.Println("Token")
	fmt.Println(tkn)
	fmt.Println("Hash")
	fmt.Println(hash)

	identifier.IdentityCheck("felipe.carbone@dito.com.br", hash)

	conn := mongoDB.Connect()
	res := mongoDB.Find(conn, bson.M{})

	fmt.Println("Consulta de teste a base do Mongo:")

	for res.Next(conn.Ctx) {
		var result bson.M

		err := res.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result["test"])
	}
}
