export enum Team {
  Unassigned = -1,
  Red,
  Blue
}

interface PlayerBase {
  name: string;
  team: Team;
  isReady: boolean;
}

export interface Player extends PlayerBase {
  id: string | null;
  sessionToken: string | null;
}

export interface OtherPlayer extends PlayerBase {
  id: string;
  connected: boolean;
}
