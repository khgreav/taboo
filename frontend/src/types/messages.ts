import type { Player } from './player';

export enum MessageType {
  ConnectMsg = 'connect',
  ConnectAckMsg = 'connect_ack',
  ChangeNameMsg = 'change_name',
  NameChangedMsg = 'name_changed',
  PlayerListMsg = 'player_list',
}

export interface MessageBase {
  type: string;
}

export interface ConnectMessage extends MessageBase {
  type: MessageType.ConnectMsg;
  playerId: string | null;
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

export interface PlayerListMessage extends MessageBase {
  type: MessageType.PlayerListMsg;
  data: Player[];
}
