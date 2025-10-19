<template>
  <div class="centered-div">
    <form @submit.prevent="sendConnect()">
      <label for="nameInput">{{ $t('components.connect.name') }}</label>
      <input
        v-model="name"
        ref="nameInput"
        type="text"
        required
      />
      <button type="submit">
        {{ $t('components.connect.action')}}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useLogStore } from '@/stores/logStore';
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { ErrCodes } from '@/types/errors';
import {
  type ConnectAckMessage,
  type ConnectMessage,
  type ErrorResponseMessage,
  type MessageBase,
  MessageType,
} from '@/types/messages';
import { storeToRefs } from 'pinia';
import { ref, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { toast } from 'vue3-toastify';

const i18n = useI18n();
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
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
        case MessageType.ErrorResponseMsg:
          if ((message as ErrorResponseMessage).failedType === MessageType.ConnectMsg) {
            handleConnectError(message as ErrorResponseMessage)
          }
      }
    });
  }
});

const sendConnect = () => {
  clientSocket.sendMessage<ConnectMessage>({
    type: MessageType.ConnectMsg,
    playerId: player.value.id,
    sessionToken: player.value.sessionToken,
    name: name.value
  });
};

const handleConnectAck = (message: ConnectAckMessage) => {
  if (player.value.id === message.playerId) {
    console.info('Connect ID matches, acknowledged.');
    return;
  }
  playerStore.setPlayerId(message.playerId);
  playerStore.setPlayerSessionToken(message.sessionToken);
  playerStore.setPlayerName(message.name);
  playerStore.setConnected(true);
  logStore.addLogRecord(
    i18n.t('messages.connections.connected', { name: message.name }),
  );
};

const handleConnectError = (message: ErrorResponseMessage) => {
  switch (message.errorCode) {
    case ErrCodes.PlayerNotFound:
      playerStore.clearSessionData();
      toast.error(
        i18n.t('messages.errors.playerNotFound')
      );
      break;
    default:
      toast.error(
        i18n.t('messages.errors.general')
      );
  }
}
</script>
