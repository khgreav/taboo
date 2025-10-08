<template>
  <header>
    <h1 class="text-center">{{ $t('title') }}</h1>
  </header>
  <div class="layout">
    <div class="left-panel">
      <GameLog v-if="connected" />
    </div>
    <div class="center-panel">
      <GamePanel />
    </div>
    <div class="right-panel">
      <PlayerPanel v-if="connected" />
    </div>
  </div>
</template>

<script setup lang="ts">
import GameLog from '@/components/GameLog.vue';
import GamePanel from '@/components/GamePanel.vue';
import PlayerPanel from '@/components/PlayerPanel.vue';
import { usePlayerStore } from '@/stores/playerStore';

import { storeToRefs } from 'pinia';

const playerStore = usePlayerStore();
const { connected } = storeToRefs(playerStore);


</script>
