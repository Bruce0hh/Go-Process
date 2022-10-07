package adapter

import "fmt"

type WindowAdapter struct {
	WindowMachine *Windows
}

func (w *WindowAdapter) InsertIntoLightningPort() {
	fmt.Printf("Adapter converts Lightning signal to USB.\n")
	w.WindowMachine.insertIntoUSBPort()
}
