<template>
  <div
    v-if="currentWord !== null"
    class="taboo-card"
  >
    <hr />
    <div class="guessed-word">{{ currentWord.word }}</div>
    <hr />
    <div class="taboo-words">{{ $t('components.tabooCard.tabooWords') }}</div>
    <ul>
      <li
        v-for="tabooWord in currentWord.taboo"
        :key="tabooWord"
        class="taboo-word"
      >
        {{ tabooWord }}
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { useWordStore } from '@/stores/wordStore';
import { storeToRefs } from 'pinia';

const wordStore = useWordStore();
const { currentWord } = storeToRefs(wordStore);

</script>
