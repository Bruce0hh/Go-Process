package chain

import "fmt"

type Medical struct {
	next Department
}

func (m *Medical) Execute(p *Patient) {
	if p.medicineDone {
		fmt.Printf("Medicine already given to patient\n")
		m.next.Execute(p)
		return
	}
	fmt.Printf("Medical giving medicine to patient\n")
	p.medicineDone = true
	m.next.Execute(p)
}

func (m *Medical) SetNext(next Department) {
	m.next = next
}
