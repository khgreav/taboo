import { ref, computed, onUnmounted } from 'vue';

export function useCountdown(defaultSeconds: number = 60) {
  const initialSeconds = ref(defaultSeconds);
  const timeLeft = ref(defaultSeconds);
  const isRunning = ref(false);
  let timer: ReturnType<typeof setInterval> | null = null;

  const startCountdown = (seconds?: number) => {
    if (typeof seconds === 'number') {
      initialSeconds.value = seconds;
      timeLeft.value = seconds;
    }
    stopCountdown(); // reset timer if running
    isRunning.value = true;
    timer = setInterval(() => {
      if (timeLeft.value > 0) {
        timeLeft.value -= 1;
      } else {
        stopCountdown();
      }
    }, 1000);
  };

  const resetCountdown = () => {
    stopCountdown();
    timeLeft.value = initialSeconds.value;
  };

  const stopCountdown = () => {
    isRunning.value = false;
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  };

  onUnmounted(() => {
    stopCountdown();
  });

  const minutes = computed(() => Math.floor(timeLeft.value / 60));
  const seconds = computed(() => timeLeft.value % 60);

  return {
    timeLeft,
    isRunning,
    startCountdown,
    stopCountdown,
    resetCountdown,
    minutes,
    seconds,
    initialSeconds,
  };
}
