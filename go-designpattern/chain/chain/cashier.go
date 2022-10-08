package chain

import "fmt"

type Cashier struct {
	next Department
}

func (c *Cashier) Execute(p *Patient) {
	if p.paymentDone {
		fmt.Printf("Payment Done\n")
	}
	fmt.Printf("Cashier getting money from patient patient\n")
}

func (c *Cashier) SetNext(next Department) {
	c.next = next
}
