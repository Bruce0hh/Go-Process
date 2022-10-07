package main

import (
	"fmt"
	"go-designpattern/adapter/adapter"
)

func main() {
	fmt.Println()
	client := &adapter.Client{}
	mac := &adapter.Mac{}

	client.InsertLightningConnectorIntoComputer(mac)
	windowsMachine := &adapter.Windows{}
	windowsMachineAdapter := &adapter.WindowAdapter{
		WindowMachine: windowsMachine,
	}
	client.InsertLightningConnectorIntoComputer(windowsMachineAdapter)
}
