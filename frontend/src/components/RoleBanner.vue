<template>
  <div>
    <div
      v-if="currentTeam === Team.Red"
      class="turn-base team-turn"
    >
      {{ $t('components.roleBanner.team.red') }}<br>
    </div>
    <div
      v-else-if="currentTeam === Team.Blue"
      class="turn-base team-turn"
    >
      {{ $t('components.roleBanner.team.blue') }}<br>
    </div>
    <div v-if="player.id === hintGiverId" class="turn-base player-turn">
      {{ $t('components.roleBanner.hintGiver.player') }}<br>
    </div>
    <div
      v-else
      class="turn-base player-turn"
    >
      {{ $t('components.roleBanner.hintGiver.other', { name: playerStore.getPlayerName(hintGiverId!) }) }}<br>
    </div>
    <div v-if="player.id === guesserId" class="turn-base player-turn">
      {{ $t('components.roleBanner.guesser.player') }}<br>
    </div>
    <div v-else class="turn-base player-turn">
      {{ $t('components.roleBanner.guesser.other', { name: playerStore.getPlayerName(guesserId!) }) }}<br>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/gameStore';
import { usePlayerStore } from '@/stores/playerStore';
import { Team } from '@/types/player';
import { storeToRefs } from 'pinia';

const gameStore = useGameStore();
const { currentTeam, guesserId, hintGiverId } = storeToRefs(gameStore);
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);

</script>
