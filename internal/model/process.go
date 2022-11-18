package model

type ProcessAction uint8

const (
	ProcessActionCreate ProcessAction = iota + 1
	ProcessActionUpdate
)

func (m *Model) Process(action ProcessAction) {

}
