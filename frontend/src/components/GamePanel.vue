<template>
  <ConnectName v-if="!connected" />
  <GameStart v-if="gameState === GameState.InProgress" />
  <WordList v-if="gameState === GameState.InRound && player.id !== guesserId"/>
</template>

<script lang="ts" setup>
import { useGameStore } from '@/stores/gameStore';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';
import ConnectName from './ConnectName.vue';
import { useSocketStore } from '@/stores/socketStore';
import { useLogStore } from '@/stores/logStore';
import { GameState, MessageType, type GameStateChangedMessage, type MessageBase, type RoundSetupMessage } from '@/types/messages';
import { useWordStore } from '@/stores/wordStore';
import { usePlayerStore } from '@/stores/playerStore';
import WordList from './WordList.vue';
import GameStart from './GameStart.vue';

const i18n = useI18n();
const gameStore = useGameStore();
const { gameState, guesserId } = storeToRefs(gameStore);
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
const { connected } = storeToRefs(playerStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();
const wordStore = useWordStore();

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message: MessageBase | null) => {
      if (!message) return;
      switch (message.type) {
        case MessageType.GameStateChangedMsg:
          handleGameStateChanged(message as GameStateChangedMessage);
          break;
        case MessageType.RoundSetupMsg:
          handleRoundSetup(message as RoundSetupMessage);
          break;
      }
    });
  }
});

const handleGameStateChanged = (message: GameStateChangedMessage) => {
  gameStore.setGameState(message.state);

  let logMsg: string;
  switch (message.state) {
    case GameState.InLobby:
      logMsg = 'messages.gameInLobby';
      break;
    case GameState.InProgress:
      logMsg = 'messages.gameInProgress';
      break;
    case GameState.InRound:
      logMsg = 'messages.gameInRound';
      break;
  }

  logStore.addLogRecord(i18n.t(logMsg));
}

const handleRoundSetup = (message: RoundSetupMessage) => {
  gameStore.setGameState(GameState.InProgress);
  gameStore.setDuration(message.duration);
  logStore.addLogRecord(
    i18n.t('messages.roundDuration', { duration: message.duration }),
  );
  gameStore.setGuesserId(message.guesserId);
  logStore.addLogRecord(
    i18n.t('messages.guesserSelected', { name: playerStore.getPlayerName(message.guesserId) }),
  );
  gameStore.setHintGiverId(message.hintGiverId);
  logStore.addLogRecord(
    i18n.t('messages.hintGiverSelected', { name: playerStore.getPlayerName(message.hintGiverId) }),
  );
  wordStore.addWords(message.words);
}
</script>
