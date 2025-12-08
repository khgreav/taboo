import type { Word } from '@/types/words';
import Denque from 'denque';
import { defineStore } from 'pinia';
import { computed, ref, type ComputedRef, type Ref } from 'vue';

export const useWordStore = defineStore('wordStore', () => {
  const words: Ref<Denque<Word>> = ref(new Denque<Word>());
  const currentWord: ComputedRef<Word | null> = computed(() => {
    return words.value.peekFront() || null;
  });

  function addWords(newWords: Word[]): void {
    for (const word of newWords) {
      words.value.push(word);
    }
  }

  function advanceWord(): void {
    words.value.shift()
  }

  function clearWords(): void {
    words.value.clear();
  }

  return {
    words,
    currentWord,
    addWords,
    advanceWord,
    clearWords,
  }
});
