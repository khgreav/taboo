import type { Word } from '@/types/words';
import { defineStore } from 'pinia';
import { ref, type Ref } from 'vue';

export const useWordStore = defineStore('wordStore', () => {
  const words: Ref<Word[]> = ref([]);
  const currentWordIndex: Ref<number> = ref(0);

  function addWords(newWords: Word[]): void {
    words.value = words.value.concat(newWords);
  }

  function advanceWord(): void {
    if (currentWordIndex.value < words.value.length - 1) {
      currentWordIndex.value += 1;
    }
  }

  return {
    words,
    addWords,
    currentWordIndex,
    advanceWord,
  }
});
