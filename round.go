package main

type Round struct {
	Team        Team
	GuesserId   string
	HintGiverId string
	Duration    int
	Words       []*TabooWord
}

func (r Round) CreateRoundSetupMessage() *RoundSetupMessage {
	return &RoundSetupMessage{
		TypeProperty: TypeProperty{Type: StartRoundMsg},
		Team:         r.Team,
		GuesserId:    r.GuesserId,
		HintGiverId:  r.HintGiverId,
		Duration:     r.Duration,
		Words:        r.Words,
	}
}

func (r Round) CreateRoundStartedMessage() *RoundStartedMessage {
	return &RoundStartedMessage{
		TypeProperty:     TypeProperty{Type: RoundStartedMsg},
		PlayerIdProperty: PlayerIdProperty{PlayerId: r.HintGiverId},
	}
}

func (r Round) CreateRoundEndedMessage() *RoundEndedMessage {
	return &RoundEndedMessage{
		TypeProperty: TypeProperty{Type: RoundEndedMsg},
	}
}
