import { GameState } from '@/types/messages';
import { defineStore } from 'pinia';
import { ref, type Ref } from 'vue';

export const useGameStore = defineStore('game', () => {
  const gameState: Ref<GameState> = ref(GameState.InLobby);
  const redScore: Ref<number> = ref(0);
  const blueScore: Ref<number> = ref(0);
  const guesserId: Ref<string | null> = ref(null);
  const hintGiverId: Ref<string | null> = ref(null);
  const duration: Ref<number> = ref(60);

  function setGameState(state: GameState): void {
    gameState.value = state;
  }

  function setRedScore(score: number): void {
    redScore.value = score;
  }

  function setBlueScore(score: number): void {
    blueScore.value = score;
  }

  function setGuesserId(id: string | null): void {
    guesserId.value = id;
  }

  function setHintGiverId(id: string | null): void {
    hintGiverId.value = id;
  }

  function setDuration(seconds: number): void {
    duration.value = seconds;
  }

  return {
    gameState,
    setGameState,
    redScore,
    setRedScore,
    blueScore,
    setBlueScore,
    guesserId,
    setGuesserId,
    hintGiverId,
    setHintGiverId,
    duration,
    setDuration,
  };
});
