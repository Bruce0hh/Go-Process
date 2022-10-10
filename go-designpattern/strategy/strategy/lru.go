package strategy

import "fmt"

type Lru struct{}

func (l *Lru) evict(c *Cache) {
	fmt.Printf("Evicting by lru strategy\n")
}
