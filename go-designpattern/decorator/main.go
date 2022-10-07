package main

import (
	"fmt"
	"go-designpattern/decorator/decorator"
)

func main() {
	pizza := &decorator.VeggieMania{}

	pizzaWithCheese := &decorator.CheeseTopping{Pizza: pizza}

	pizzaWithCheeseAndTomato := &decorator.TomatoTopping{Pizza: pizzaWithCheese}
	fmt.Printf("Price of veggeMania with tomato and cheese topping is %d\n", pizzaWithCheeseAndTomato.GetPrice())
}
