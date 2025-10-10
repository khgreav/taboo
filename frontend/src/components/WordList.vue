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
  <div v-if="player.id === hintGiverId">
    <button
      :disabled="currentWordIndex === words.length - 1"
      @click="guessWord()"
    >
      {{ $t('components.wordList.actions.guess') }}
    </button>
    <button
      :disabled="currentWordIndex === words.length - 1"
      @click="skipWord()"
    >
      {{ $t('components.wordList.actions.skip') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/gameStore';
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { useWordStore } from '@/stores/wordStore';
import { MessageType } from '@/types/messages';
import { storeToRefs } from 'pinia';

const clientSocket = useSocketStore();
const gameStore = useGameStore();
const { hintGiverId } = storeToRefs(gameStore);
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
const wordStore = useWordStore();
const { words, currentWordIndex } = storeToRefs(wordStore);

const guessWord = () => {
  clientSocket.sendMessage({
    type: MessageType.GuessWordMsg,
    playerId: player.value.id,
  });
}

const skipWord = () => {
  clientSocket.sendMessage({
    type: MessageType.SkipWordMsg,
    playerId: player.value.id,
  });
}
</script>
