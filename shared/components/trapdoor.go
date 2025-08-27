package components

type TrapDoorEnum int

type TrapDoor struct {
	IsOpenCallbackID TrapDoorEnum
	Open             bool
	LastChangedTick  int `json:"-"`
}
