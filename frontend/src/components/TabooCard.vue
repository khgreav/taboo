<template>
  <div v-if="currentWord !== null">
    <h2>{{ $t('components.tabooCard.currentWord') }}</h2>
    <hr />
    <b>{{ currentWord.word }}</b>
    <hr />
    <h3>{{ $t('components.tabooCard.tabooWords') }}</h3>
    <ul>
      <li v-for="tabooWord in currentWord.taboo" :key="tabooWord">
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
