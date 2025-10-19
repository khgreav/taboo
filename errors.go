package main

type ErrorCode int

const (
	ErrGameFull = iota
	ErrPlayerNotFound
	ErrSessionTokenMissing
	ErrGameNotInLobby
	ErrTeamFull
	ErrPlayerNotInTeam
	ErrGameNotStarted
	ErrNotHintGiver
	ErrRoundNotActive
)

func GetErrMessage(code ErrorCode) string {
	switch code {
	case ErrGameFull:
		return "Game is full."
	case ErrPlayerNotFound:
		return "Player ID does not exist."
	case ErrSessionTokenMissing:
		return "Session token is missing."
	case ErrGameNotInLobby:
		return "Game not in lobby state."
	case ErrTeamFull:
		return "Team is full."
	case ErrPlayerNotInTeam:
		return "Player has not selected team yet."
	case ErrGameNotStarted:
		return "Game has not started yet."
	case ErrNotHintGiver:
		return "Only hint giver can start a round."
	case ErrRoundNotActive:
		return "Round is not active."
	default:
		return "Unknown error."
	}
}

func CreateErrorMessage(messageType MessageType, code ErrorCode) *ErrorResponseMessage {
	return &ErrorResponseMessage{
		TypeProperty: TypeProperty{
			Type: ErrorResponseMsg,
		},
		FailedType: messageType,
		Error:      GetErrMessage(code),
		ErrorCode:  code,
	}
}
