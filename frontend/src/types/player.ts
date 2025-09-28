interface PlayerBase {
  name: string;
  isReady: boolean;
}

export interface Player extends PlayerBase {
  id: string | null;
}

export interface OtherPlayer extends PlayerBase {
  id: string;
}
