<template>
  <div v-if="words.length > 0">
    <h2>Word to guess</h2>
    <hr />
    <b>{{ words[currentWordIndex].word }}</b>
    <hr />
    <h3>Taboo words</h3>
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
