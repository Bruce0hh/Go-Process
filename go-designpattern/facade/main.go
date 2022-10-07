package main

import (
	"fmt"
	"go-designpattern/facade/facade"
	"log"
)

func main() {
	fmt.Println()
	walletFacade := facade.NewWalletFacade("abc", 1234)
	fmt.Println()

	if err := walletFacade.AddMoneyToWallet("abc", 1234, 10); err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
	fmt.Println()
	if err := walletFacade.DeductMoneyFromWallet("abc", 1234, 5); err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
}
