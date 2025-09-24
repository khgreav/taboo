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
import { usePlayerStore } from '@/stores/player';
import { useSocketStore } from '@/stores/socketStore';
import { MessageType, type PlayerListMessage } from '@/types/messages';
import type { Player } from '@/types/player';
import { storeToRefs } from 'pinia';

const playerStore = usePlayerStore();
const { playerId } = storeToRefs(playerStore);
const gameStore = useGameStore();
const { playerList } = storeToRefs(gameStore);
const clientSocket = useSocketStore();

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message) => {
      if (!message) return;
      // Handle player list updates
      if (message.type === MessageType.PlayerListMsg) {
        handlePlayerList(message as PlayerListMessage);
      }
    });
  }
});

const handlePlayerList = (message: PlayerListMessage) => {
  message.data.sort((a: Player, b: Player) => {
    if (a.id === playerId.value) {
      return -1;
    }
    if (b.id === playerId.value) {
      return 1;
    }
    return a.name.localeCompare(b.name);
  });
  gameStore.setPlayerList(message.data);
};
</script>
