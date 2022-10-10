package strategy

import "fmt"

type Lfu struct{}

func (l *Lfu) evict(c *Cache) {
	fmt.Printf("Evicting by lfu strategy\n")
}
