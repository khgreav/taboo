<template>
  <div class="centered-div">
    <label for="nameInput">Select name</label>
    <input
      v-model="name"
      ref="nameInput"
      type="text"
      required
    />
      <button @click="sendConnect()">
    Connect
  </button>
  </div>

</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/game';
import { useSocketStore } from '@/stores/socketStore';
import {
  type ConnectAckMessage,
  type ConnectMessage,
  type MessageBase,
  MessageType,
} from '@/types/messages';
import { storeToRefs } from 'pinia';
import { ref, type Ref } from 'vue';

const gameStore = useGameStore();
const { playerId } = storeToRefs(gameStore);
// socket store
const clientSocket = useSocketStore();

const name: Ref<string> = ref('');

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message: MessageBase | null) => {
      if (!message) return;
      switch (message.type) {
        case MessageType.ConnectAckMsg:
          handleConnectAck(message as ConnectAckMessage);
          break;
      }
    });
  }
});

const sendConnect = () => {
  clientSocket.sendMessage<ConnectMessage>({
    type: MessageType.ConnectMsg,
    playerId: playerId.value,
    name: name.value
  });
};

const handleConnectAck = (message: ConnectAckMessage) => {
  if (playerId.value === message.playerId) {
    console.info('Connect ID matches, acknowledged.');
    return;
  }
  console.info('Player ID mismatch, assigning new:', playerId.value, message.playerId);
  gameStore.setPlayerId(message.playerId);
  gameStore.setConnected(true);
};
</script>
