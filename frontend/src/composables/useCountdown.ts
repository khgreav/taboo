import { computed, ref } from 'vue';

export function useCountdown(seconds: number = 60) {
  const updateMs = 250;
  const remainingMs = ref(seconds * 1000);
  const remainingSeconds = computed(() => Math.ceil(remainingMs.value / 1000));
  let timer: number | null = null;
  let endTime: number | null = null;
  let pausedTime = seconds * 1000;

  const startCountdown = (seconds?: number) => {
    if (timer) {
      // already running
      return;
    }
    if (seconds !== undefined) {
      remainingMs.value = seconds * 1000;
      pausedTime = seconds * 1000;
    }
    endTime = Date.now() + pausedTime;
    updateTime();
    timer = setInterval(updateTime, updateMs);
  }

  const updateTime = () => {
    if (!endTime) {
      return;
    }
    const now = Date.now();
    remainingMs.value = Math.max(0, endTime - now);
    if (remainingMs.value <= 0) {
      stopCountdown();
    }
  }

  const stopCountdown = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
    if (endTime) {
      pausedTime = Math.max(0, endTime - Date.now());
    }
    endTime = null;
  }

  const resetCountdown = () => {
    stopCountdown();
    remainingMs.value = seconds * 1000;
    pausedTime = seconds * 1000;
  }

  const adjustRemaining = (seconds: number) => {
    remainingMs.value = seconds * 1000;
    pausedTime = seconds * 1000;
  }

  return {
    remainingSeconds,
    startCountdown,
    stopCountdown,
    resetCountdown,
    adjustRemaining,
  }
}
