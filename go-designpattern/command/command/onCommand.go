package command

type OnCommand struct {
	Device Device
}

func (o *OnCommand) execute() {
	o.Device.on()
}
