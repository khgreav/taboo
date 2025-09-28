import type { OtherPlayer, Player } from '@/types/player';
import { defineStore } from 'pinia';
import { ref, type Ref } from 'vue';

export const useGameStore = defineStore('game', () => {
  const connected: Ref<boolean> = ref(false);
  const running: Ref<boolean> = ref(false);
  const player: Ref<Player> = ref({
    id: null,
    name: '',
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

  function setRunning(isRunning: boolean): void {
    running.value = isRunning;
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
    setPlayerReadyState,
    connected,
    setConnected,
    running,
    setRunning,
    playerList,
    setPlayerList,
    addPlayer,
    removePlayer,
    setPlayerReady,
  };
});
