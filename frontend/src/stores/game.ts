import type { Player } from '@/types/player';
import { defineStore } from 'pinia';
import { ref, type Ref } from 'vue';

export const useGameStore = defineStore('game', () => {
  const connected: Ref<boolean> = ref(false);
  const running: Ref<boolean> = ref(false);
  const playerId: Ref<string | null> = ref(null);
  const playerName: Ref<string> = ref('Player');
  const playerList: Ref<Player[]> = ref([]);

  function setPlayerId(id: string) {
    playerId.value = id;
  }

  function setPlayerName(name: string) {
    playerName.value = name;
  }

  function setConnected(status: boolean) {
    connected.value = status;
  }

  function setRunning(isRunning: boolean) {
    running.value = isRunning;
  }

  function setPlayerList(players: Player[]) {
    playerList.value = players;
    sortPlayerList();
  }

  function addPlayer(player: Player) {
    playerList.value.push(player);
    sortPlayerList();
  }

  function removePlayer(playerId: string) {
    const idx = playerList.value.findIndex(player => player.id === playerId);
    if (idx === -1) {
      return;
    }
    playerList.value.splice(idx, 1);
  }

  function sortPlayerList() {
    playerList.value.sort((a: Player, b: Player) => {
      if (a.id === playerId.value) {
        return -1;
      }
      if (b.id === playerId.value) {
        return 1;
      }
      return a.name.localeCompare(b.name);
    });
  }

  return {
    playerId,
    setPlayerId,
    playerName,
    setPlayerName,
    connected,
    setConnected,
    running,
    setRunning,
    playerList,
    setPlayerList,
    addPlayer,
    removePlayer
  };
});
