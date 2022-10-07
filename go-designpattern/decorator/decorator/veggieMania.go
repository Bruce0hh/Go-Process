package decorator

type VeggieMania struct {
}

func (v *VeggieMania) GetPrice() int {
	return 15
}

var _ IPizza = (*VeggieMania)(nil)
