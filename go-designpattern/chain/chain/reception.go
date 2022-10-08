package chain

import "fmt"

type Reception struct {
	next Department
}

func (r *Reception) Execute(p *Patient) {
	if p.registrationDone {
		fmt.Printf("Patient registration already done\n")
		r.next.Execute(p)
		return
	}
	fmt.Printf("Reception registering patient\n")
	p.registrationDone = true
	r.next.Execute(p)
}

func (r *Reception) SetNext(next Department) {
	r.next = next
}
