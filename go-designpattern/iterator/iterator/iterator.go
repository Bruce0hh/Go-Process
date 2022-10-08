package iterator

type Iterator interface {
	HasNext() bool
	GetNext() *User
}
