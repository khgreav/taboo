package main

import "time"

type Round struct {
	Team        Team
	GuesserId   string
	HintGiverId string
	Duration    int
	Words       []*TabooWord
	StartTime   int64
}

func (r *Round) SetDuration(duration int) {
	r.Duration = duration
}

func (r Round) CreateRoundSetupMessage() *RoundSetupMessage {
	return &RoundSetupMessage{
		TypeProperty: TypeProperty{Type: RoundSetupMsg},
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

func (r Round) CreateRoundPausedMessage() *RoundPausedMessage {
	return &RoundPausedMessage{
		TypeProperty:      TypeProperty{Type: RoundPausedMsg},
		RemainingDuration: r.Duration,
	}
}

func (r Round) CreateRoundResumedMessage() *RoundResumedMessage {
	return &RoundResumedMessage{
		TypeProperty:     TypeProperty{Type: RoundResumedMsg},
		PlayerIdProperty: PlayerIdProperty{PlayerId: r.HintGiverId},
	}
}

func (r Round) CalculateRoundPausedDuration() int {
	endTime := r.StartTime + int64(r.Duration)*1000
	remaining := endTime - time.Now().UnixMilli()
	if remaining < 0 {
		return 0
	}
	return int(remaining / 1000)
}
