package adapter

import "fmt"

type Mac struct {
}

func (m *Mac) InsertIntoLightningPort() {
	fmt.Printf("Lightning connector is plugged into mac machine.\n")
}
