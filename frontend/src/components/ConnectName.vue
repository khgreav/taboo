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
import { useGameStore } from '@/stores/gameStore';
import { useLogStore } from '@/stores/logStore';
import { useSocketStore } from '@/stores/socketStore';
import {
  type ConnectAckMessage,
  type ConnectMessage,
  type MessageBase,
  MessageType,
} from '@/types/messages';
import { storeToRefs } from 'pinia';
import { ref, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const gameStore = useGameStore();
const { player } = storeToRefs(gameStore);
const logStore = useLogStore();
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
    playerId: player.value.id,
    name: name.value
  });
};

const handleConnectAck = (message: ConnectAckMessage) => {
  if (player.value.id === message.playerId) {
    console.info('Connect ID matches, acknowledged.');
    return;
  }
  gameStore.setPlayerId(message.playerId);
  gameStore.setPlayerName(message.name);
  gameStore.setConnected(true);
  logStore.addLogRecord(
    i18n.t('messages.connected', { name: message.name }),
  );
};
</script>
