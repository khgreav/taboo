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
  const playerList: Ref<OtherPlayer[]> = ref([]);

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
    const item = playerList.value.find(player => player.id === id);
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

  function setPlayerList(players: OtherPlayer[]): void {
    playerList.value = players.filter(p => p.id !== player.value.id);
    sortPlayerList();
  }

  function addPlayer(player: OtherPlayer): void {
    playerList.value.push(player);
    sortPlayerList();
  }

  function removePlayer(playerId: string): void {
    const idx = playerList.value.findIndex(player => player.id === playerId);
    if (idx === -1) {
      return;
    }
    playerList.value.splice(idx, 1);
  }

  function setPlayerTeam(playerId: string, team: Team): void {
    if (playerId === player.value.id) {
      player.value.team = team;
      player.value.isReady = false;
    }
    const item = playerList.value.find(player => player.id === playerId);
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
    const item = playerList.value.find(player => player.id === playerId);
    if (!item) {
      return;
    }
    item.isReady = isReady;
  }

  function setPlayerConnected(playerId: string, connected: boolean): void {
    const item = playerList.value.find(player => player.id === playerId);
    if (!item) {
      return;
    }
    item.connected = connected;
  }

  function resetPlayerTeams(): void {
    player.value.team = Team.Unassigned;
    player.value.isReady = false;
    playerList.value.forEach((p: OtherPlayer) => {
      p.team = Team.Unassigned;
      p.isReady = false;
    });
  }

  function sortPlayerList(): void {
    playerList.value.sort((a: OtherPlayer, b: OtherPlayer) => {
      if (a.id === player.value.id) {
        return -1;
      }
      if (b.id === player.value.id) {
        return 1;
      }
      return a.name.localeCompare(b.name);
    });
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
    playerList,
    setPlayerList,
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
