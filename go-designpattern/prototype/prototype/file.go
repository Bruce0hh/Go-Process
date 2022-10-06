package prototype

import "fmt"

type File struct {
	Name string
}

func (f *File) Print(indentation string) {
	fmt.Printf("%s %s\n", indentation, f.Name)
}

func (f *File) Clone() Inode {
	return &File{Name: f.Name + "_clone"}
}

var _ Inode = (*File)(nil)
