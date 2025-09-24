<template>
  <h1>Taboo MP!</h1>
  <button
    class="bg-sky-700 hover:bg-sky-500 text-white px-4 py-2 border-2 border-black rounded"
    @click="sendConnect()"
  >
    Connect
  </button>
  <div>
    <p>
      Player ID: {{ playerId }}<br />
      Player name: {{ playerName }}
    </p>
    <input
      ref="nameChangeInput"
      type="text"
      v-model="nameChange"
      required
    />
    <button @click='setName(nameChange)'>Set Name</button>
  </div>
  <PlayerList />
</template>

<script setup lang="ts">
import PlayerList from '@/components/PlayerList.vue';
import { usePlayerStore } from '@/stores/player';
import { useSocketStore } from '@/stores/socketStore';
import {
  MessageType,
  type ConnectAckMessage,
  type MessageBase,
  type NameChangedMessage,
} from './types/messages';
import { storeToRefs } from 'pinia';
import { ref, type Ref } from 'vue';

// player store
const playerStore = usePlayerStore();
const { playerId, playerName } = storeToRefs(playerStore);
// socket store
const clientSocket = useSocketStore();
clientSocket.init('ws://localhost:8081/ws');

const nameChange: Ref<string> = ref('');

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message: MessageBase | null) => {
      if (!message) return;
      console.warn('Received message:', message);
      switch (message.type) {
        case MessageType.ConnectAckMsg:
          handleConnectAck(message as ConnectAckMessage);
          break;
        case MessageType.NameChangedMsg:
          handleNameChanged(message as NameChangedMessage);
          break;
        default:
          console.warn('Unhandled message type:', message.type);
      }
    });
  }
});

const sendConnect = () => {
  clientSocket.sendMessage({
    type: MessageType.ConnectMsg,
    playerId: playerId.value,
  });
};

const setName = (name: string) => {
  clientSocket.sendMessage({
    type: MessageType.ChangeNameMsg,
    playerId: playerId.value,
    name: name,
  });
};

const handleConnectAck = (message: ConnectAckMessage) => {
  if (playerId.value === message.playerId) {
    console.info('Connect ID matches, acknowledged.');
    return;
  }
  console.info('Player ID mismatch, assigning new:', playerId.value, message.playerId);
  playerStore.setPlayerId(message.playerId);
};

const handleNameChanged = (message: NameChangedMessage) => {
  console.info('Name changed to:', message.name);
  playerStore.setPlayerName(message.name);
};
</script>
