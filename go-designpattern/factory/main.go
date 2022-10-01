package main

import (
	"fmt"
	"go-designpattern/factory/method"
)

func main() {
	ak47, _ := method.GetGun("ak47")
	musket, _ := method.GetGun("musket")

	printDetails(ak47)
	printDetails(musket)
}

func printDetails(g method.IGun) {
	fmt.Printf("Gun: %s\n", g.GetName())
	fmt.Printf("Power: %d\n", g.GetPower())
}
