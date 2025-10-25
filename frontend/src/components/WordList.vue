<template>
  <div v-if="words.length > 0">
    <h2>{{ $t('components.wordList.currentWord') }}</h2>
    <hr />
    <b>{{ words[currentWordIndex].word }}</b>
    <hr />
    <h3>{{ $t('components.wordList.tabooWords') }}</h3>
    <ul>
      <li v-for="tabooWord in words[currentWordIndex].taboo" :key="tabooWord">
        {{ tabooWord }}
      </li>
    </ul>
  </div>

</template>

<script setup lang="ts">
import { useWordStore } from '@/stores/wordStore';
import { storeToRefs } from 'pinia';

const wordStore = useWordStore();
const { words, currentWordIndex } = storeToRefs(wordStore);

</script>
