import { ref, type Ref } from 'vue';
import { defineStore } from 'pinia';

export const usePlayerStore = defineStore('player', () => {
  const playerId: Ref<string | null> = ref(null);
  const name: Ref<string> = ref('Player');

  return { playerId, name };
});
