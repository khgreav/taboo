<template>
  <div
    id="overlay"
    class="overlay"
  >
    <div class="overlay-text">{{ $t('components.overlay.reconnecting') }}</div>
  </div>
</template>

<script setup lang="ts">
import { useSocketStore } from '@/stores/socketStore';
import { storeToRefs } from 'pinia';
import { watch } from 'vue';

const socketStore = useSocketStore();
const { connected: isConnected } = storeToRefs(socketStore);

watch(isConnected, (newVal) => {
  document.getElementById('overlay')?.classList.toggle('active', !newVal);
});
</script>
