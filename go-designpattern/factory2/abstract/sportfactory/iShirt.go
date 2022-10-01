package sportfactory

type IShirt interface {
	setLogo(logo string)
	setSize(size int)
	GetLogo() string
	GetSize() int
}

type Shirt struct {
	Logo string
	Size int
}

func (s *Shirt) setLogo(logo string) {
	s.Logo = logo
}

func (s *Shirt) setSize(size int) {
	s.Size = size
}

func (s *Shirt) GetLogo() string {
	return s.Logo
}

func (s *Shirt) GetSize() int {
	return s.Size
}

var _ IShirt = (*Shirt)(nil)
