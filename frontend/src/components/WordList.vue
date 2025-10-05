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

<script setup lang="ts">
import { useGameStore } from '@/stores/gameStore';
import { useLogStore } from '@/stores/logStore';
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { useWordStore } from '@/stores/wordStore';
import { MessageType, type WordGuessedMessage, type WordListMessage, type WordSkippedMessage } from '@/types/messages';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
const gameStore = useGameStore();
const wordStore = useWordStore();
const { words, currentWordIndex } = storeToRefs(wordStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();

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
        case MessageType.WordGuessedMsg:
          handleWordGuessed(message as WordGuessedMessage);
          break;
      }
    });
  }
});

const handleWordList = (message: WordListMessage) => {
  wordStore.addWords(message.words);
};

const guessWord = () => {
  clientSocket.sendMessage({
    type: MessageType.GuessWordMsg,
    playerId: player.value.id,
  });
}

const handleWordGuessed = (message: WordGuessedMessage) => {
  gameStore.setRedScore(message.redScore);
  gameStore.setBlueScore(message.blueScore);
  wordStore.advanceWord();
  const playerName = playerStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t('messages.guessed', { name: playerName }),
  );
};

const skipWord = () => {
  clientSocket.sendMessage({
    type: MessageType.SkipWordMsg,
    playerId: player.value.id,
  });
}

const handleWordSkipped = (message: WordSkippedMessage) => {
  void message; // TODO use message if needed
  wordStore.advanceWord();
  const playerName = playerStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t('messages.skipped', { name: playerName }),
  );
};
</script>
