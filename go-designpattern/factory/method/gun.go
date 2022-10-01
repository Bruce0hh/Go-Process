package method

type Gun struct {
	name  string
	power int
}

func (g *Gun) setName(name string) {
	g.name = name
}

func (g *Gun) setPower(power int) {
	g.power = power
}

func (g *Gun) GetName() string {
	return g.name
}

func (g *Gun) GetPower() int {
	return g.power
}

var _ IGun = (*Gun)(nil)
