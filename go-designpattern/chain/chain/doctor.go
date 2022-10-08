package chain

import "fmt"

type Doctor struct {
	next Department
}

func (d *Doctor) Execute(p *Patient) {
	if p.doctorCheckUpDone {
		fmt.Printf("Doctor checkup already done\n")
		d.next.Execute(p)
		return
	}
	fmt.Printf("Doctor checking patient\n")
	p.doctorCheckUpDone = true
	d.next.Execute(p)
}

func (d *Doctor) SetNext(next Department) {
	d.next = next
}
