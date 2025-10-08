<template>
  <div>
    <button
      :disabled="currentWordIndex === words.length - 1"
      @click="guessWord()"
    >
      {{ $t('actions.guess') }}
    </button>
    <button
      :disabled="currentWordIndex === words.length - 1"
      @click="skipWord()"
    >
      {{ $t('actions.skip') }}
    </button>
  </div>
</template>

<script lang="ts" setup>
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { useWordStore } from '@/stores/wordStore';
import { MessageType } from '@/types/messages';
import { storeToRefs } from 'pinia';

const clientSocket = useSocketStore();
const wordStore = useWordStore();
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
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
