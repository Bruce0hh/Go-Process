package chain

type Department interface {
	Execute(*Patient)
	SetNext(Department)
}
