package nike

import (
	"go-designpattern/factory2/abstract/sportfactory"
)

type Nike struct {
}

func (n *Nike) MakeShoe() sportfactory.IShoe {
	return &NikeShoe{sportfactory.Shoe{
		Logo: "nike",
		Size: 16,
	}}
}

func (n *Nike) MakeShirt() sportfactory.IShirt {
	return &NikeShirt{sportfactory.Shirt{
		Logo: "nike",
		Size: 16,
	}}
}
