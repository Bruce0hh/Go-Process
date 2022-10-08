package main

import "go-designpattern/chain/chain"

func main() {

	cashier := &chain.Cashier{}

	medical := &chain.Medical{}
	medical.SetNext(cashier)

	doctor := &chain.Doctor{}
	doctor.SetNext(medical)

	reception := &chain.Reception{}
	reception.SetNext(doctor)

	patient := &chain.Patient{Name: "abc"}
	reception.Execute(patient)
}
