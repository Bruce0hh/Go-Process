package builder

type IglooBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func (i *IglooBuilder) setWindowType() {
	i.windowType = "Snow Window"
}

func (i *IglooBuilder) setDoorType() {
	i.doorType = "Snow Door"
}

func (i *IglooBuilder) setNumFloor() {
	i.floor = 1
}

func (i *IglooBuilder) getHouse() House {
	return House{
		windowType: i.windowType,
		DoorType:   i.doorType,
		floorNum:   i.floor,
	}
}

func newIglooBuilder() *IglooBuilder {
	return &IglooBuilder{}
}

var _ IBuilder = (*IglooBuilder)(nil)
