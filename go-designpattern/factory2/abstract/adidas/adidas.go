package adidas

import (
	"go-designpattern/factory2/abstract/sportfactory"
)

type Adidas struct {
}

func (a *Adidas) MakeShoe() sportfactory.IShoe {
	return &AdidasShoe{sportfactory.Shoe{
		Logo: "adidas",
		Size: 14,
	}}
}

func (a *Adidas) MakeShirt() sportfactory.IShirt {
	return &AdidasShirt{sportfactory.Shirt{
		Logo: "adidas",
		Size: 14,
	}}
}
