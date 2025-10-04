import type { OtherPlayer } from './player';
import type { Word } from './words';

export enum MessageType {
  ConnectMsg = 'connect',
  ConnectAckMsg = 'connect_ack',
  ChangeNameMsg = 'change_name',
  NameChangedMsg = 'name_changed',
  PlayerJoinedMsg = 'player_joined',
  PlayerLeftMsg = 'player_left',
  PlayerListMsg = 'player_list',
  ChangeTeamMsg = 'change_team',
  TeamChangedMsg = 'team_changed',
  PlayerReadyMsg = 'player_ready',
  GameStateChangedMsg = 'game_state_changed',
  WordListMsg = 'word_list',
  SkipWordMsg = 'skip_word',
  GuessWordMsg = 'guess_word',
  WordGuessedMsg = 'word_guessed',
  WordSkippedMsg = 'word_skipped',
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
