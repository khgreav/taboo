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
    <hr />
    <button :disabled="currentWordIndex === words.length - 1">Correct guess</button>
    <button
      :disabled="currentWordIndex === words.length - 1"
      @click="skipWord()"
    >
      {{ $t('actions.skip') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/game';
import { useLogStore } from '@/stores/logStore';
import { useSocketStore } from '@/stores/socketStore';
import { MessageType, type WordListMessage, type WordSkippedMessage } from '@/types/messages';
import type { Word } from '@/types/words';
import { storeToRefs } from 'pinia';
import { ref, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const gameStore = useGameStore();
const { player } = storeToRefs(gameStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();
const words: Ref<Word[]> = ref([]);
const currentWordIndex: Ref<number> = ref(0);

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message) => {
      if (!message) return;
      switch (message.type) {
        case MessageType.WordListMsg:
          handleWordList(message as WordListMessage);
          break;
        case MessageType.WordSkippedMsg:
          handleWordSkipped(message as WordSkippedMessage);
          break;
      }
    });
  }
});

const handleWordList = (message: WordListMessage) => {
  words.value = words.value.concat(message.words);
};

const skipWord = () => {
  clientSocket.sendMessage({
    type: MessageType.SkipWordMsg,
    playerId: player.value.id, // TODO: get player ID from store
  });
}

const handleWordSkipped = (message: WordSkippedMessage) => {
  void message; // TODO use message if needed
  currentWordIndex.value++;
  const playerName = gameStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t('messages.skipped', { name: playerName }),
  );
};
</script>
