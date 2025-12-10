import type { OtherPlayer, Team } from './player';
import type { Word } from './words';

export enum MessageType {
  // general messages
  ErrorResponseMsg = 'error_response',
  // player connections
  ConnectMsg = 'connect',
  ConnectAckMsg = 'connect_ack',
  ReconnectMsg = 'reconnect',
  ReconnectAckMsg = 'reconnect_ack',
  PlayerJoinedMsg = 'player_joined',
  PlayerLeftMsg = 'player_left',
  PlayerDisconnectedMsg = 'player_disconnected',
  PlayerReconnectedMsg = 'player_reconnected',
  // lobby state
  PlayerListMsg = 'player_list',
  ChangeTeamMsg = 'change_team',
  TeamChangedMsg = 'team_changed',
  PlayerReadyMsg = 'player_ready',
  GameStateChangedMsg = 'game_state_changed',
  // game rounds
  RoundSetupMsg = 'round_setup',
  StartRoundMsg = 'start_round',
  RoundStartedMsg = 'round_started',
  RoundEndedMsg = 'round_ended',
  RoundPausedMsg = 'round_paused',
  ResumeRoundMsg = 'resume_round',
  RoundResumedMsg = 'round_resumed',
  GameEndedMsg = 'game_ended',
  ResetGameMsg = 'reset_game',
  GameResetMsg = 'game_reset',
  // round actions
  SkipWordMsg = 'skip_word',
  WordSkippedMsg = 'word_skipped',
  GuessWordMsg = 'guess_word',
  WordGuessedMsg = 'word_guessed',
  WordListMsg = 'word_list',
}

export enum GameState {
  InLobby = 0,
  InProgress,
  InRound,
  RoundPaused,
  Ended,
}

export interface MessageBase {
  type: string;
}

export interface ErrorResponseMessage extends MessageBase {
  type: MessageType.ErrorResponseMsg;
  failedType: MessageType;
  error: string;
  errorCode: number;
}

export interface ConnectMessage extends MessageBase {
  type: MessageType.ConnectMsg;
  name: string;
}

export interface ConnectAckMessage extends MessageBase {
  type: MessageType.ConnectAckMsg;
  playerId: string;
  sessionToken: string;
  name: string;
}

export interface ReconnectMessage extends MessageBase {
  type: MessageType.ReconnectMsg;
  playerId: string;
  sessionToken: string;
}

export interface ReconnectAckMessage extends Omit<ConnectAckMessage, 'type'> {
  type: MessageType.ReconnectAckMsg;
  team: Team.Red | Team.Blue;
  state: GameState.InProgress | GameState.InRound | GameState.RoundPaused;
  remainingDuration: number;
  currentTeam: Team.Red | Team.Blue;
  guesserId: string;
  hintGiverId: string;
  redScore: number;
  blueScore: number;
  words: Word[]
}

export interface PlayerLeftMessage extends MessageBase {
  type: MessageType.PlayerLeftMsg;
  playerId: string;
}

export interface PlayerDisconnectedMessage extends MessageBase {
  type: MessageType.PlayerDisconnectedMsg;
  playerId: string;
}

export interface PlayerReconnectedMessage extends MessageBase {
  type: MessageType.PlayerReconnectedMsg;
  playerId: string;
}

export interface PlayerJoinedMessage extends MessageBase {
  type: MessageType.PlayerJoinedMsg;
  playerId: string;
  name: string;
}

export interface PlayerListMessage extends MessageBase {
  type: MessageType.PlayerListMsg;
  players: OtherPlayer[];
}

export interface ChangeTeamMessage extends MessageBase {
  type: MessageType.ChangeTeamMsg;
  playerId: string;
  team: number;
}

export interface TeamChangedMessage extends MessageBase {
  type: MessageType.TeamChangedMsg;
  playerId: string;
  team: number;
}

export interface PlayerReadyMessage extends MessageBase {
  type: MessageType.PlayerReadyMsg;
  playerId: string;
  isReady: boolean;
}

export interface GameStateChangedMessage extends MessageBase {
  type: MessageType.GameStateChangedMsg;
  state: GameState;
}

export interface WordListMessage extends MessageBase {
  type: MessageType.WordListMsg;
  words: Word[];
}

export interface SkipWordMessage extends MessageBase {
  type: MessageType.SkipWordMsg;
  playerId: string;
}

export interface WordSkippedMessage extends MessageBase {
  type: MessageType.WordSkippedMsg;
  playerId: string;
}

export interface GuessWordMessage extends MessageBase {
  type: MessageType.GuessWordMsg;
  playerId: string;
}

export interface WordGuessedMessage extends MessageBase {
  type: MessageType.WordGuessedMsg;
  playerId: string;
  redScore: number;
  blueScore: number;
}

export interface RoundSetupMessage extends MessageBase {
  type: MessageType.RoundSetupMsg;
  team: Team;
  guesserId: string;
  hintGiverId: string;
  duration: number;
  words: Word[];
}

export interface RoundStartedMessage extends MessageBase {
  type: MessageType.RoundStartedMsg;
  playerId: string;
}

export interface StartRoundMessage extends MessageBase {
  type: MessageType.StartRoundMsg;
  playerId: string;
}

export interface RoundPausedMessage extends MessageBase {
  type: MessageType.RoundPausedMsg;
  remainingDuration: number;
}

export interface ResumeRoundMessage extends MessageBase {
  type: MessageType.ResumeRoundMsg;
  playerId: string;
}

export interface RoundResumedMessage extends MessageBase {
  type: MessageType.RoundResumedMsg;
  playerId: string;
}

export interface GameEndedMsg extends MessageBase {
  type: MessageType.GameEndedMsg;
  redScore: number;
  blueScore: number;
}

export interface ResetGameMessage extends MessageBase {
  type: MessageType.ResetGameMsg;
  playerId: string;
}

export interface GameResetMessage extends MessageBase {
  type: MessageType.GameResetMsg;
}
