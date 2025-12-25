<template>
  <div class="default-border">
    <form
      v-if="player.id !== null && player.sessionToken !== null"
      @submit.prevent='reconnect()'
    >
      <p>{{ $t('components.connect.reconnectMessage') }}</p>
      <button type='submit'>
        {{ $t('components.connect.actions.reconnect')}}
      </button>
    </form>
    <form
      v-else
      @submit.prevent='connect()'
    >
      <label for="playerName">{{ $t('components.connect.name') }}</label>
      <input
        id="playerName"
        ref="nameInput"
        v-model="name"
        type="text"
        required
      /><br><br>
      <button type='submit'>
        {{ $t('components.connect.actions.connect')}}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/gameStore';
import { useLogStore } from '@/stores/logStore';
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { useWordStore } from '@/stores/wordStore';
import { ErrCodes } from '@/types/errors';
import {
  type ConnectAckMessage,
  type ConnectMessage,
  type ErrorResponseMessage,
  type MessageBase,
  MessageType,
  type ReconnectAckMessage,
  type ReconnectMessage,
} from '@/types/messages';
import { storeToRefs } from 'pinia';
import { ref, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { toast } from 'vue3-toastify';

const emit = defineEmits<{
  updateDuration: [];
}>();
const i18n = useI18n();
const gameStore = useGameStore();
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
const wordStore = useWordStore();
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
        case MessageType.ReconnectAckMsg:
          handleReconnectAck(message as ReconnectAckMessage);
          break;
        case MessageType.ErrorResponseMsg:
          if ((message as ErrorResponseMessage).failedType === MessageType.ConnectMsg) {
            handleConnectError(message as ErrorResponseMessage)
          }
      }
    });
  }
});

const connect = () => {
  clientSocket.sendMessage<ConnectMessage>({
    type: MessageType.ConnectMsg,
    name: name.value
  });
};

const reconnect = () => {
  clientSocket.sendMessage<ReconnectMessage>({
    type: MessageType.ReconnectMsg,
    playerId: player.value.id!,
    sessionToken: player.value.sessionToken!,
  });
}

const handleConnectAck = (message: ConnectAckMessage) => {
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

const handleReconnectAck = (message: ReconnectAckMessage) => {
  // set player team and name
  playerStore.setPlayerTeam(message.playerId, message.team);
  playerStore.setPlayerName(message.name);
  // set current team and player roles
  gameStore.setCurrentTeam(message.currentTeam);
  gameStore.setHintGiverId(message.hintGiverId);
  gameStore.setGuesserId(message.guesserId);
  // set scores
  gameStore.setRedScore(message.redScore);
  gameStore.setBlueScore(message.blueScore);
  // set remaining duration and words to guess
  gameStore.setDuration(message.remainingDuration);
  emit('updateDuration');
  wordStore.addWords(message.words);
  // finally, set game state and player connected state to update UI
  gameStore.setGameState(message.state);
  playerStore.setConnected(true);
  logStore.addLogRecord(
    i18n.t('messages.connections.reconnected', { name: message.name }),
  );
}
</script>
