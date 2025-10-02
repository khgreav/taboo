package main

type Round struct {
	Team        Team
	GuesserId   string
	HintGiverId string
	StartTime   int64
	Duration    int
	Words       []*TabooWord
}

func (r *Round) CreateStartRoundMessage() *StartRoundMessage {
	return &StartRoundMessage{
		TypeProperty: TypeProperty{Type: StartRoundMsg},
		Team:         r.Team,
		GuesserId:    r.GuesserId,
		HintGiverId:  r.HintGiverId,
		StartTime:    r.StartTime,
		Duration:     r.Duration,
		Words:        r.Words,
	}
}
