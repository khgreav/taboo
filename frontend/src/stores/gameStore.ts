import { GameState } from '@/types/messages';
import { Team } from '@/types/player';
import { defineStore } from 'pinia';
import { computed, ref, type Ref } from 'vue';

export const useGameStore = defineStore('game', () => {
  const gameState: Ref<GameState> = ref(GameState.InLobby);
  const redScore: Ref<number> = ref(0);
  const blueScore: Ref<number> = ref(0);
  const currentTeam: Ref<Team | null> = ref(null);
  const guesserId: Ref<string | null> = ref(null);
  const hintGiverId: Ref<string | null> = ref(null);
  const duration: Ref<number> = ref(60);
  const winner = computed(() => {
    if (redScore.value > blueScore.value) {
      return Team.Red;
    } else if (blueScore.value > redScore.value) {
      return Team.Blue;
    } else {
      return null;
    }
  });

  function setGameState(state: GameState): void {
    gameState.value = state;
  }

  function setRedScore(score: number): void {
    redScore.value = score;
  }

  function setBlueScore(score: number): void {
    blueScore.value = score;
  }

  function setCurrentTeam(team: Team | null): void {
    currentTeam.value = team;
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

  function resetGame(): void {
    gameState.value = GameState.InLobby;
    redScore.value = 0;
    blueScore.value = 0;
    currentTeam.value = null;
    guesserId.value = null;
    hintGiverId.value = null;
    duration.value = 60;
  }

  return {
    gameState,
    setGameState,
    redScore,
    setRedScore,
    blueScore,
    setBlueScore,
    currentTeam,
    setCurrentTeam,
    guesserId,
    setGuesserId,
    hintGiverId,
    setHintGiverId,
    duration,
    setDuration,
    winner,
    resetGame,
  };
});
