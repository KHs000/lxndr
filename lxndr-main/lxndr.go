package main

import (
	"fmt"
	
	"github.com/KHs000/lxndr/rndtoken"
)

func main() {
	hash := rndtoken.GenerateToken("felilpe.carbone@dito.com.br")

	fmt.Println("Primeiro hash")
	fmt.Println(hash)
}
