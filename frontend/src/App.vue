<template>
  <h1>Taboo MP!</h1>
  <div class="layout">
    <div class="left-panel">
      <ConnectName v-if="!connected" />
      <NameChange v-else />
      <GameLog />
    </div>
    <div class="center-panel">
      <WordList />
    </div>
    <div class="right-panel">
      <PlayerList />
    </div>
  </div>
</template>

<script setup lang="ts">
import ConnectName from '@/components/ConnectName.vue';
import GameLog from '@/components/GameLog.vue';
import NameChange from '@/components/NameChange.vue';
import PlayerList from '@/components/PlayerList.vue';
import WordList from '@/components/WordList.vue';
import { useGameStore } from '@/stores/game';

import { storeToRefs } from 'pinia';

const gameStore = useGameStore();
const { connected } = storeToRefs(gameStore);

</script>
