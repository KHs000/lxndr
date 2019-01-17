package main

import (
	"fmt"

	"github.com/KHs000/lxndr/identifier"
	"github.com/KHs000/lxndr/rndtoken"
)

func main() {
	tkn, hash := rndtoken.SendToken("felipe.carbone@dito.com.br")

	fmt.Println("Token")
	fmt.Println(tkn)
	fmt.Println("Hash")
	fmt.Println(hash)

	identifier.IdentityCheck("felipe.carbone@dito.com.br", hash)

}
