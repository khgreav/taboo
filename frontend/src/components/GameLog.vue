<template>
  <div>
    <label>{{ $t('labels.log') }}</label>
    <div v-for="(log, index) in logs" :key="index" class="log-entry">
      <span class="timestamp">{{ new Date(log.timestamp).toLocaleTimeString() }}:</span>
      <span class="message">{{ log.message }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLogStore } from '@/stores/logStore';
import { storeToRefs } from 'pinia';

const logStore = useLogStore();
const { logs } = storeToRefs(logStore);
</script>
