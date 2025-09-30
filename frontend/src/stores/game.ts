import { GameState } from '@/types/messages';
import { Team, type OtherPlayer, type Player } from '@/types/player';
import { defineStore } from 'pinia';
import { ref, type Ref } from 'vue';

export const useGameStore = defineStore('game', () => {
  const connected: Ref<boolean> = ref(false);
  const gameState: Ref<GameState> = ref(GameState.InLobby);
  const player: Ref<Player> = ref({
    id: null,
    name: '',
    team: Team.Unassigned,
    isReady: false,
  })
  const playerList: Ref<OtherPlayer[]> = ref([]);

  function setPlayerId(id: string): void {
    player.value.id = id;
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

  function setConnected(status: boolean): void {
    connected.value = status;
  }

  function setGameState(state: GameState): void {
    gameState.value = state;
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
    }
    const item = playerList.value.find(player => player.id === playerId);
    if (!item) {
      return;
    }
    item.team = team;
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

  function sortPlayerList(): void {
    playerList.value.sort((a: Player, b: Player) => {
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
    getPlayerName,
    setPlayerName,
    setPlayerTeam,
    setPlayerReadyState,
    connected,
    setConnected,
    gameState,
    setGameState,
    playerList,
    setPlayerList,
    addPlayer,
    removePlayer,
    setPlayerReady,
  };
});
