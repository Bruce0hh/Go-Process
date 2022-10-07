package decorator

type TomatoTopping struct {
	Pizza IPizza
}

func (t *TomatoTopping) GetPrice() int {
	price := t.Pizza.GetPrice()
	return price + 7
}

var _ IPizza = (*TomatoTopping)(nil)
