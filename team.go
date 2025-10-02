package main

type Team int

const (
	Unassigned Team = -1
	Red        Team = 0
	Blue       Team = 1
)

func (t Team) GetOppositeTeam() Team {
	switch t {
	case Red:
		return Blue
	case Blue:
		return Red
	default:
		return Unassigned
	}
}
