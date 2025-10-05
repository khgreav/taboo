<template>
  <div>
    <label for="nameChangeInput">Change name</label>
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
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { MessageType, type MessageBase, type NameChangedMessage } from '@/types/messages';
import { storeToRefs } from 'pinia';
import { ref, type Ref } from 'vue';

// player store
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
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
    playerId: player.value.id,
    name: name,
  });
};

const handleNameChanged = (message: NameChangedMessage) => {
  playerStore.setPlayerName(message.name);
};
</script>
