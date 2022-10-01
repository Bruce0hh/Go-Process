package abstract

type Nike struct {
}

func (n *Nike) MakeShoe() IShoe {
	return &NikeShoe{Shoe{
		logo: "nike",
		size: 16,
	}}
}

func (n *Nike) MakeShirt() IShirt {
	return &NikeShirt{Shirt{
		logo: "nike",
		size: 16,
	}}
}

var _ ISportsFactory = (*Nike)(nil)
