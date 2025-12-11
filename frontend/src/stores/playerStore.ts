import { type Player, Team, type OtherPlayer } from '@/types/player';
import { defineStore } from 'pinia';
import { type Ref, ref } from 'vue';

export const usePlayerStore = defineStore('playerStore', () => {
  const connected: Ref<boolean> = ref(false);
  const player: Ref<Player> = ref({
    id: null,
    sessionToken: null,
    name: '',
    team: Team.Unassigned,
    isReady: false,
  })
  const playerMap: Ref<Map<string, OtherPlayer>> = ref(new Map());

  function setConnected(status: boolean): void {
    connected.value = status;
  }

  function setPlayerId(id: string): void {
    player.value.id = id;
  }

  function setPlayerSessionToken(sessionToken: string): void {
    player.value.sessionToken = sessionToken;
  }

  function clearPlayerId(): void {
    player.value.id = null;
  }

  function clearPlayerSessionToken(): void {
    player.value.sessionToken = null;
  }

  function clearPlayerName(): void {
    player.value.name = '';
  }

  function clearSessionData(): void {
    clearPlayerId();
    clearPlayerSessionToken();
    clearPlayerName();
  }

  function getPlayerName(id: string): string | null {
    if (player.value.id === id) {
      return player.value.name;
    }
    const item = playerMap.value.get(id);
    if (!item) {
      return null;
    }
    return item.name;
  }

  function setPlayerName(name: string): void {
    player.value.name = name;
  }

  function setPlayerReadyState(isReady: boolean): void {
    player.value.isReady = isReady;
  }

  function setPlayers(players: OtherPlayer[]): void {
    const map = new Map<string, OtherPlayer>();
    players.forEach(item => {
      if (item.id === player.value.id) {
        return;
      }
      map.set(item.id, item);
    });
    playerMap.value = map;
  }

  function addPlayer(player: OtherPlayer): void {
    playerMap.value.set(player.id, player);
  }

  function removePlayer(playerId: string): void {
    playerMap.value.delete(playerId);
  }

  function setPlayerTeam(playerId: string, team: Team): void {
    if (playerId === player.value.id) {
      player.value.team = team;
      player.value.isReady = false;
    }
    const item = playerMap.value.get(playerId);
    if (!item) {
      return;
    }
    item.team = team;
    item.isReady = false;
  }

  function setPlayerReady(playerId: string, isReady: boolean): void {
    if (playerId === player.value.id) {
      player.value.isReady = isReady;
    }
    const item = playerMap.value.get(playerId);
    if (!item) {
      return;
    }
    item.isReady = isReady;
  }

  function setPlayerConnected(playerId: string, connected: boolean): void {
    const item = playerMap.value.get(playerId);
    if (!item) {
      return;
    }
    item.connected = connected;
  }

  function resetPlayerTeams(remainingPlayers: Set<string>): void {
    player.value.team = Team.Unassigned;
    player.value.isReady = false;
    for (const id of playerMap.value.keys()) {
      if (remainingPlayers.has(id)) {
        const item = playerMap.value.get(id)!;
        item.team = Team.Unassigned;
        item.isReady = false;
      } else {
        playerMap.value.delete(id);
      }
    }
  }

  return {
    player,
    setPlayerId,
    clearPlayerId,
    setPlayerSessionToken,
    clearPlayerSessionToken,
    clearSessionData,
    getPlayerName,
    setPlayerName,
    setPlayerTeam,
    setPlayerReadyState,
    connected,
    setConnected,
    setPlayerConnected,
    playerMap,
    setPlayers,
    addPlayer,
    removePlayer,
    setPlayerReady,
    resetPlayerTeams,
  };
}, {
  persist: [
    {
      pick: ["player.id", "player.sessionToken", "player.name"],
    },
  ],
});
