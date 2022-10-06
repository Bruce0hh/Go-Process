package prototype

import "fmt"

type Folder struct {
	Children []Inode
	Name     string
}

func (f *Folder) Print(indentation string) {
	fmt.Printf("%s %s\n", indentation, f.Name)
	for _, child := range f.Children {
		child.Print(indentation + indentation)
	}
}

func (f *Folder) Clone() Inode {
	cloneFolder := &Folder{Name: f.Name + "_clone"}
	var tempChildren []Inode
	for _, child := range f.Children {
		copyChild := child.Clone()
		tempChildren = append(tempChildren, copyChild)
	}
	cloneFolder.Children = tempChildren
	return cloneFolder
}

var _ Inode = (*Folder)(nil)
