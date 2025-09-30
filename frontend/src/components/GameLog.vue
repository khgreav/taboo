<template>
  <label>{{ $t('sections.log') }}</label>
  <div class="log">
    <div
      v-for="(log, index) in logs"
      :key="index"
    >
      <span>{{ `${new Date(log.timestamp).toLocaleTimeString()}: ${log.message}` }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLogStore } from '@/stores/logStore';
import { storeToRefs } from 'pinia';

const logStore = useLogStore();
const { logs } = storeToRefs(logStore);
</script>
