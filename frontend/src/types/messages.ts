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
  PlayerReadyMsg = 'player_ready',
  WordListMsg = 'word_list',
  SkipWordMsg = 'skip_word',
  WordSkippedMsg = 'word_skipped',
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

export interface PlayerReadyMessage extends MessageBase {
  type: MessageType.PlayerReadyMsg;
  playerId: string;
  isReady: boolean;
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
