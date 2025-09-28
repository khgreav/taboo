import type { LogRecord } from '@/types/logs';
import { defineStore } from 'pinia';
import { ref, type Ref } from 'vue';

export const useLogStore = defineStore('log', () => {
  const logs: Ref<LogRecord[]> = ref<LogRecord[]>([]);

  function addLogRecord(message: string) {
    logs.value.push({
      timestamp: Date.now(),
      message,
    });
  }

  function clearLogs() {
    logs.value = [];
  }

  return {
    logs,
    addLogRecord,
    clearLogs
  }
});
