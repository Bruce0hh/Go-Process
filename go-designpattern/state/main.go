package main

import (
	"fmt"
	"go-designpattern/state/state"
	"log"
)

func main() {
	vendingMachine := state.NewVendingMachine(1, 10)

	if err := vendingMachine.RequestItem(); err != nil {
		log.Fatalf(err.Error())
	}

	if err := vendingMachine.InsertMoney(10); err != nil {
		log.Fatalf(err.Error())
	}

	if err := vendingMachine.DispenseItem(); err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println()

	if err := vendingMachine.AddItem(2); err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println()

	if err := vendingMachine.RequestItem(); err != nil {
		log.Fatalf(err.Error())
	}

	if err := vendingMachine.InsertMoney(10); err != nil {
		log.Fatalf(err.Error())
	}

	if err := vendingMachine.DispenseItem(); err != nil {
		log.Fatalf(err.Error())
	}
}
