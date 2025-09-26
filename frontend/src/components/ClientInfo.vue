<template>
  <div>
    <p>
      Player ID: {{ playerId }}<br />
      Player name: {{ playerName }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/game';
import { storeToRefs } from 'pinia';

const gameStore = useGameStore();
const { playerId, playerName } = storeToRefs(gameStore);
</script>
