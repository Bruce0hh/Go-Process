package main

import "go-designpattern/observer/observer"

func main() {

	shirtItem := observer.NewItem("Nike Shirt")

	observerFirst := &observer.Customer{Id: "abc@gmail.com"}
	observerSecond := &observer.Customer{Id: "xyz@gmail.com"}

	shirtItem.Register(observerFirst)
	shirtItem.Register(observerSecond)

	shirtItem.UpdateAvailability()
}
