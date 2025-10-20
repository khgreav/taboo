<template>
  <div class="centered centered-column">
    <ConnectName v-if="!connected" />
    <RoleBanner v-if="gameState !== GameState.InLobby" />
    <GameStart v-if="gameState === GameState.InProgress" />
    <div v-if="gameState === GameState.InRound">
      <div>
        <h3>
          {{ `${$t('components.wordList.roundTime')}: ${remainingSeconds}` }}
        </h3>
      </div>
      <WordList v-if="player.id !== guesserId" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';

import ConnectName from '@/components/ConnectName.vue';
import WordList from '@/components/WordList.vue';
import GameStart from '@/components/GameStart.vue';
import RoleBanner from '@/components/RoleBanner.vue';
import { useCountdown } from '@/composables/useCountdown';
import { useGameStore } from '@/stores/gameStore';
import { useLogStore } from '@/stores/logStore';
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { useWordStore } from '@/stores/wordStore';
import {
  GameState,
  MessageType,
  type GameStateChangedMessage,
  type MessageBase,
  type RoundSetupMessage,
  type RoundStartedMessage,
  type WordGuessedMessage,
  type WordListMessage,
  type WordSkippedMessage,
} from '@/types/messages';

const i18n = useI18n();
const gameStore = useGameStore();
const { gameState, guesserId, duration } = storeToRefs(gameStore);
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
const { connected } = storeToRefs(playerStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();
const wordStore = useWordStore();
const { startCountdown, stopCountdown, remainingSeconds } = useCountdown(60);

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
        case MessageType.WordListMsg:
          handleWordList(message as WordListMessage);
          break;
        case MessageType.RoundStartedMsg:
          handleRoundStarted(message as RoundStartedMessage);
          break;
        case MessageType.WordGuessedMsg:
          handleWordGuessed(message as WordGuessedMessage);
          break;
        case MessageType.WordSkippedMsg:
          handleWordSkipped(message as WordSkippedMessage);
          break;
        case MessageType.RoundEndedMsg:
          handleRoundEnded();
          break;
      }
    });
  }
});

const handleGameStateChanged = (message: GameStateChangedMessage) => {
  gameStore.setGameState(message.state);

  switch (message.state) {
    case GameState.InLobby:
      logStore.addLogRecord(i18n.t('messages.gameState.inLobby'));
      break;
    case GameState.InProgress:
      logStore.addLogRecord(i18n.t('messages.gameState.inProgress'));
      break;
    case GameState.InRound:
      logStore.addLogRecord(i18n.t('messages.gameState.inRound'));
      break;
    case GameState.Ended:
      logStore.addLogRecord(i18n.t('messages.gameState.ended'));
      break;
  }
}

const handleRoundSetup = (message: RoundSetupMessage) => {
  gameStore.setGameState(GameState.InProgress);
  gameStore.setDuration(message.duration);
  logStore.addLogRecord(
    i18n.t(
      'messages.round.duration',
      { duration: message.duration },
    ),
  );
  gameStore.setCurrentTeam(message.team)
  gameStore.setHintGiverId(message.hintGiverId);
  logStore.addLogRecord(
    i18n.t(
      'messages.round.hintGiver',
      { name: playerStore.getPlayerName(message.hintGiverId) },
    ),
  );
  gameStore.setGuesserId(message.guesserId);
  logStore.addLogRecord(
    i18n.t(
      'messages.round.guesser',
      { name: playerStore.getPlayerName(message.guesserId) },
    ),
  );
  wordStore.addWords(message.words);
}

const handleRoundStarted = (message: RoundStartedMessage) => {
  gameStore.setGameState(GameState.InRound);
  logStore.addLogRecord(
    i18n.t(
      'messages.round.started',
      { name: playerStore.getPlayerName(message.playerId) },
    ),
  );
  startCountdown(duration.value);
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
    i18n.t(
      'messages.round.guessed',
      { name: playerName },
    ),
  );
};

const handleWordSkipped = (message: WordSkippedMessage) => {
  wordStore.advanceWord();
  const playerName = playerStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t(
      'messages.round.skipped',
      { name: playerName },
    ),
  );
};

const handleWordList = (message: WordListMessage) => {
  wordStore.addWords(message.words);
};

const handleRoundEnded = () => {
  gameStore.setGameState(GameState.InProgress);
  stopCountdown();
}
</script>
