package adapter

import "fmt"

type Client struct {
}

func (c *Client) InsertLightningConnectorIntoComputer(com Computer) {
	fmt.Printf("Client inserts Lightning connector into computer.\n")
	com.InsertIntoLightningPort()
}
