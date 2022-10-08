package command

import "fmt"

type TV struct {
	isRunning bool
}

func (t *TV) on() {
	t.isRunning = true
	fmt.Printf("Turning tv on\n")
}

func (t *TV) off() {
	t.isRunning = false
	fmt.Printf("Turning tv on\n")
}

var _ Device = (*TV)(nil)
