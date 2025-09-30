<template>
  <div class="layout">
    <div class="left-panel">
      <ConnectName v-if="!connected" />
      <GameLog v-if="connected" />
    </div>
    <div class="center-panel">
      <WordList />
    </div>
    <div class="right-panel">
      <PlayerPanel v-if="connected" />
    </div>
  </div>
</template>

<script setup lang="ts">
import ConnectName from '@/components/ConnectName.vue';
import GameLog from '@/components/GameLog.vue';
import PlayerPanel from '@/components/PlayerPanel.vue';
import WordList from '@/components/WordList.vue';
import { useGameStore } from '@/stores/game';
import { useSocketStore } from '@/stores/socketStore';

import { storeToRefs } from 'pinia';
import { GameState, MessageType, type GameStateChangedMessage, type MessageBase } from './types/messages';
import { useLogStore } from './stores/logStore';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const gameStore = useGameStore();
const { connected } = storeToRefs(gameStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message: MessageBase | null) => {
      if (!message) return;
      switch (message.type) {
        case MessageType.GameStateChangedMsg:
          handleGameStateChanged(message as GameStateChangedMessage);
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

</script>
