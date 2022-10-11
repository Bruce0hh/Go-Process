package state

import "fmt"

type HasMoneyState struct {
	vendingMachine *VendingMachine
}

func (h *HasMoneyState) AddItem(count int) error {
	return fmt.Errorf("item dispense in progress")
}

func (h *HasMoneyState) RequestItem() error {
	return fmt.Errorf("item dispense in progress")
}

func (h *HasMoneyState) InsertMoney(money int) error {
	return fmt.Errorf("item out of stock")
}

func (h *HasMoneyState) DispenseItem() error {
	fmt.Println("Dispensing Item")
	h.vendingMachine.itemCount = h.vendingMachine.itemCount - 1
	if h.vendingMachine.itemCount == 0 {
		h.vendingMachine.setState(h.vendingMachine.noItem)
	} else {
		h.vendingMachine.setState(h.vendingMachine.hasItem)
	}
	return nil
}
