<template>
  <div class="player-list">
    <h2>Player List</h2>
    <ul>
      <li
        v-for="player in playerList"
        :key="player.id"
        :class="{ 'current-player': player.id === playerId }"
      >
        {{ player.name }}
      </li>
    </ul>
  </div>
</template>

<script lang="ts" setup>
import { useGameStore } from '@/stores/game';
import { useSocketStore } from '@/stores/socketStore';
import { MessageType, type PlayerJoinedMessage, type PlayerLeftMessage, type PlayerListMessage } from '@/types/messages';
import { storeToRefs } from 'pinia';

const gameStore = useGameStore();
const { playerId, playerList } = storeToRefs(gameStore);
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
      }
    });
  }
});

const handlePlayerList = (message: PlayerListMessage) => {
  gameStore.setPlayerList(message.players);
};

const handlePlayerJoined = (message: PlayerJoinedMessage) => {
  gameStore.addPlayer({ id: message.playerId, name: message.name });
};

const handlePlayerLeft = (message: PlayerLeftMessage) => {
  gameStore.removePlayer(message.playerId);
};
</script>
