<template>
  <div class="default-border">
    <div class="side-panel-title">
      {{ $t('sections.log') }}
    </div>
    <div id="logContainer" class="log-entries">
      <div
        v-for="(log, index) in logs"
        :key="index"
      >
        <span>{{ `${new Date(log.timestamp).toLocaleTimeString()}: ${log.message}` }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLogStore } from '@/stores/logStore';
import { storeToRefs } from 'pinia';
import { nextTick, watch } from 'vue';

const logStore = useLogStore();
const { logs } = storeToRefs(logStore);

watch(() => logs.value.length, async () => {
  const logDiv = document.getElementById('logContainer');
  if (!logDiv) {
    return;
  }
  if (logDiv.scrollHeight > logDiv.clientHeight) {
    await nextTick();
    logDiv.scrollTop = logDiv.scrollHeight;
  }
});
</script>
