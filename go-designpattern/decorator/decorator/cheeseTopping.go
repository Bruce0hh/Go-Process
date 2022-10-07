package decorator

type CheeseTopping struct {
	Pizza IPizza
}

func (c *CheeseTopping) GetPrice() int {
	price := c.Pizza.GetPrice()
	return price + 10
}
