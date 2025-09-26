import type { Player } from './player';

export enum MessageType {
  ConnectMsg = 'connect',
  ConnectAckMsg = 'connect_ack',
  ChangeNameMsg = 'change_name',
  NameChangedMsg = 'name_changed',
  PlayerJoinedMsg = 'player_joined',
  PlayerLeftMsg = 'player_left',
  PlayerListMsg = 'player_list',
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
  players: Player[];
}
