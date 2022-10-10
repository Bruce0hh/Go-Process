package strategy

import "fmt"

type Fifo struct {
}

func (f *Fifo) evict(c *Cache) {
	fmt.Printf("Evicting by fifo strategy\n")
}
