package prototype

type Inode interface {
	Print(string)
	Clone() Inode
}
