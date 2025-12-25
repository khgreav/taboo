package main

type Team int

const MaxTeamMembers = 2

const (
	Unassigned Team = -1
	Red        Team = 0
	Blue       Team = 1
)

type TeamState struct {
	MemberCount int
	Score       int
}

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
