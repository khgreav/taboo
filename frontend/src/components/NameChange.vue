<template>
  <div>
    <input
      ref="nameChangeInput"
      type="text"
      v-model="nameChange"
      required
    />
    <button @click='setName(nameChange)'>Set Name</button>
  </div>
</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/game';
import { useSocketStore } from '@/stores/socketStore';
import { MessageType, type MessageBase, type NameChangedMessage } from '@/types/messages';
import { storeToRefs } from 'pinia';
import { ref, type Ref } from 'vue';

// player store
const gameStore = useGameStore();
const { playerId } = storeToRefs(gameStore);
// socket store
const clientSocket = useSocketStore();
const nameChange: Ref<string> = ref('');

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message: MessageBase | null) => {
      if (!message) return;
      switch (message.type) {
        case MessageType.NameChangedMsg:
          handleNameChanged(message as NameChangedMessage);
          break;
      }
    });
  }
});

const setName = (name: string) => {
  clientSocket.sendMessage({
    type: MessageType.ChangeNameMsg,
    playerId: playerId.value,
    name: name,
  });
};

const handleNameChanged = (message: NameChangedMessage) => {
  console.info('Name changed to:', message.name);
  gameStore.setPlayerName(message.name);
};
</script>
