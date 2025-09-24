import { ref, type Ref } from 'vue';
import { defineStore } from 'pinia';

export const usePlayerStore = defineStore('player', () => {
  const playerId: Ref<string | null> = ref(null);
  const playerName: Ref<string> = ref('Player');

  function setPlayerId(id: string) {
    playerId.value = id;
  }

  function setPlayerName(name: string) {
    playerName.value = name;
  }

  return {
    playerId,
    playerName,
    setPlayerId,
    setPlayerName,
  };
});
