<template>
  <div class="player-list">
    <label>Player List</label>
    <ul>
      <li
        v-if="player.id"
        class="current-player"
      >
        {{ player.name }}
        <span
          :style="{
            color: player.isReady ? 'green' : 'red',
          }"
        >
          {{ `${player.isReady ? "READY" : " NOT READY"}` }}
        </span>
      </li>
      <li
        v-for="item in playerList"
        :key="item.id"
        :class="{ 'current-player': item.id === player.id }"
      >
        {{ item.name }}
        <span
          :style="{
            color: item.isReady ? 'green' : 'red',
          }"
        >
          {{ `${item.isReady ? "READY" : " NOT READY"}` }}
        </span>
      </li>
    </ul>
    <button
      @click="readyUp()"
    >
      {{ readyButtonText }}
    </button>
  </div>
</template>

<script lang="ts" setup>
import { useGameStore } from '@/stores/game';
import { useLogStore } from '@/stores/logStore';
import { useSocketStore } from '@/stores/socketStore';
import {
  MessageType,
  type PlayerJoinedMessage,
  type PlayerLeftMessage,
  type PlayerListMessage,
  type PlayerReadyMessage,
} from '@/types/messages';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const gameStore = useGameStore();
const { player, playerList } = storeToRefs(gameStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message) => {
      if (!message) return;
      switch (message.type) {
        case MessageType.PlayerJoinedMsg:
          handlePlayerJoined(message as PlayerJoinedMessage);
          break;
        case MessageType.PlayerListMsg:
          handlePlayerList(message as PlayerListMessage);
          break;
        case MessageType.PlayerLeftMsg:
          handlePlayerLeft(message as PlayerLeftMessage);
          break;
        case MessageType.PlayerReadyMsg:
          handlePlayerReady(message as PlayerReadyMessage);
          break;
      }
    });
  }
});

const readyButtonText = computed(() => {
  return player.value.isReady ? 'Unready' : 'Ready Up';
})

const readyUp = () => {
  clientSocket.sendMessage({
    type: MessageType.PlayerReadyMsg,
    playerId: player.value.id,
    isReady: !player.value.isReady,
  });
};

const handlePlayerList = (message: PlayerListMessage) => {
  gameStore.setPlayerList(message.players);
};

const handlePlayerJoined = (message: PlayerJoinedMessage) => {
  gameStore.addPlayer({ id: message.playerId, name: message.name, isReady: false });

  logStore.addLogRecord(
    i18n.t('messages.playerJoined', { name: message.name }),
  );
};

const handlePlayerLeft = (message: PlayerLeftMessage) => {
  gameStore.removePlayer(message.playerId);
  const playerName = gameStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t('messages.playerLeft', { name: playerName }),
  );
};

const handlePlayerReady = (message: PlayerReadyMessage) => {
  gameStore.setPlayerReady(message.playerId, message.isReady);
  const playerName = gameStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t(message.isReady ? 'messages.ready' : 'messages.notReady', { name: playerName }),
  );
};
</script>
