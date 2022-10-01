package sportfactory

type ISportsFactory interface {
	MakeShoe() IShoe
	MakeShirt() IShirt
}
