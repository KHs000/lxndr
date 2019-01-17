package main

import (
	"fmt"

	"main/identifier"
	"main/rndtoken"
)

func main() {
	tkn, hash := rndtoken.SendToken("felipe.carbone@dito.com.br")

	fmt.Println("Token")
	fmt.Println(tkn)
	fmt.Println("Hash")
	fmt.Println(hash)

	identifier.IdentityCheck("felipe.carbone@dito.com.br", hash)

}
