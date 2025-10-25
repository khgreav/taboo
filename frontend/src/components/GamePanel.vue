<template>
  <div class="centered centered-column">
    <ConnectName
      v-if="!connected"
      @update-duration='adjustRemaining(duration)'
    />
    <RoleBanner v-if="gameState !== GameState.InLobby" />
    <GameStart v-if="gameState === GameState.InProgress" />
    <div v-if="[GameState.InRound, GameState.RoundPaused].includes(gameState)">
      <div>
        <h3>
          {{ `${$t('components.wordList.roundTime')}: ${remainingSeconds}` }}
        </h3>
      </div>
      <WordList v-if="player.id !== guesserId" />
      <div v-if="player.id === hintGiverId">
        <button
          :disabled="currentWordIndex === words.length - 1"
          @click="guessWord()"
          >
          {{ $t('components.controls.guess') }}
        </button>
        <button
          :disabled="currentWordIndex === words.length - 1"
          @click="skipWord()"
          >
          {{ $t('components.controls.skip') }}
        </button>
        <button
          v-if='gameState === GameState.RoundPaused'
          @click='resumeRound()'
        >
          {{ $t('components.controls.resume') }}
        </button>
      </div>
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
  type ResumeRoundMessage,
  type RoundPausedMessage,
  type RoundResumedMessage,
  type RoundSetupMessage,
  type RoundStartedMessage,
  type WordGuessedMessage,
  type WordListMessage,
  type WordSkippedMessage,
} from '@/types/messages';

const i18n = useI18n();
const gameStore = useGameStore();
const { gameState, guesserId, hintGiverId, duration } = storeToRefs(gameStore);
const playerStore = usePlayerStore();
const { player, connected } = storeToRefs(playerStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();
const wordStore = useWordStore();
const { currentWordIndex, words } = storeToRefs(wordStore);
const { startCountdown, stopCountdown, adjustRemaining, remainingSeconds } = useCountdown(60);

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
        case MessageType.RoundPausedMsg:
          handleRoundPaused(message as RoundPausedMessage);
          break;
        case MessageType.RoundResumedMsg:
          handleRoundResumed(message as RoundResumedMessage);
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
    i18n.t(
      'messages.round.guessed',
      { name: playerName },
    ),
  );
};

const skipWord = () => {
  clientSocket.sendMessage({
    type: MessageType.SkipWordMsg,
    playerId: player.value.id,
  });
}

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
  logStore.addLogRecord(
    i18n.t('messages.round.ended'),
  );
}

const handleRoundPaused = (message: RoundPausedMessage) => {
  gameStore.setGameState(GameState.RoundPaused);
  stopCountdown();
  adjustRemaining(message.remainingDuration);
  logStore.addLogRecord(
    i18n.t('messages.round.paused'),
  );
}

const resumeRound = () => {
  clientSocket.sendMessage({
    type: MessageType.ResumeRoundMsg,
    playerId: player.value.id,
  } as ResumeRoundMessage);
}

const handleRoundResumed = (message: RoundResumedMessage) => {
  gameStore.setGameState(GameState.InRound);
  startCountdown(remainingSeconds.value);
  const playerName = playerStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t(
      'messages.round.resumed',
      { name: playerName },
    ),
  );
}
</script>
