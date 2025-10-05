import type { OtherPlayer, Team } from './player';
import type { Word } from './words';

export enum MessageType {
  // general messages
  ErrorResponseMsg = 'error_response',
  // player connections
  ConnectMsg = 'connect',
  ConnectAckMsg = 'connect_ack',
  PlayerJoinedMsg = 'player_joined',
  PlayerLeftMsg = 'player_left',
  // lobby state
  PlayerListMsg = 'player_list',
  ChangeNameMsg = 'change_name',
  NameChangedMsg = 'name_changed',
  ChangeTeamMsg = 'change_team',
  TeamChangedMsg = 'team_changed',
  PlayerReadyMsg = 'player_ready',
  GameStateChangedMsg = 'game_state_changed',
  // game rounds
  RoundSetupMsg = 'round_setup',
  StartRoundMsg = 'start_round',
  RoundStartedMsg = 'round_started',
  RoundEndedMsg = 'round_ended',
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
  InRound
}

export interface MessageBase {
  type: string;
}

export interface ConnectMessage extends MessageBase {
  type: MessageType.ConnectMsg;
  playerId: string | null;
  name: string;
}

export interface ConnectAckMessage extends MessageBase {
  type: MessageType.ConnectAckMsg;
  playerId: string;
  name: string;
}

export interface ChangeNameMessage extends MessageBase {
  type: MessageType.ChangeNameMsg;
  playerId: string;
  name: string;
}

export interface NameChangedMessage extends MessageBase {
  type: MessageType.NameChangedMsg;
  playerId: string;
  name: string;
}

export interface PlayerLeftMessage extends MessageBase {
  type: MessageType.PlayerLeftMsg;
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

export interface StartRoundMessage extends MessageBase {
  type: MessageType.StartRoundMsg;
  playerId: string;
}

export interface RoundStartedMessage extends MessageBase {
  type: MessageType.RoundStartedMsg;
  playerId: string;
}
