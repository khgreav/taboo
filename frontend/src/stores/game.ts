import type { Player } from '@/types/player';
import { defineStore } from 'pinia';
import { ref, type Ref } from 'vue';

export const useGameStore = defineStore('game', () => {
  const running: Ref<boolean> = ref(false);
  const playerList: Ref<Player[]> = ref([]);

  function setRunning(isRunning: boolean) {
    running.value = isRunning;
  }

  function setPlayerList(players: Player[]) {
    playerList.value = players;
  }

  return {
    running,
    playerList,
    setRunning,
    setPlayerList,
  };
});