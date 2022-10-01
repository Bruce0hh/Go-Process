package abstract

type Adidas struct {
}

func (a *Adidas) MakeShoe() IShoe {
	return &AdidasShoe{Shoe{
		logo: "adidas",
		size: 14,
	}}
}

func (a *Adidas) MakeShirt() IShirt {
	return &AdidasShirt{Shirt{
		logo: "adidas",
		size: 14,
	}}
}

var _ ISportsFactory = (*Adidas)(nil)
