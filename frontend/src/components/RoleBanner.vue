<template>
  <div>
    <span v-if="currentTeam === Team.Red">
      {{ $t('components.roleBanner.team.red') }}<br>
    </span>
    <span v-else-if="currentTeam === Team.Blue">
      {{ $t('components.roleBanner.team.blue') }}<br>
    </span>
    <span v-if="player.id === hintGiverId">
      {{ $t('components.roleBanner.hintGiver.player') }}<br>
    </span>
    <span v-else>
      {{ $t('components.roleBanner.hintGiver.other', { name: playerStore.getPlayerName(hintGiverId!) }) }}<br>
    </span>
    <span v-if="player.id === guesserId">
      {{ $t('components.roleBanner.guesser.player') }}<br>
    </span>
    <span v-else>
      {{ $t('components.roleBanner.guesser.other', { name: playerStore.getPlayerName(guesserId!) }) }}<br>
    </span>
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
