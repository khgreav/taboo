<template>
  <div>
    <span v-if="player.id === hintGiverId">
      <button
        @click="startRound()"
      >
        {{ $t('actions.startRound') }}
      </button>
    </span>
    <span v-else>
      {{ $t('messages.waitingToStart', { name: playerStore.getPlayerName(hintGiverId!) }) }}
    </span>
  </div>
</template>

<script lang="ts" setup>
import { useGameStore } from '@/stores/gameStore';
import { useLogStore } from '@/stores/logStore';
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { GameState, MessageType, type MessageBase, type RoundStartedMessage } from '@/types/messages';
import { storeToRefs } from 'pinia';
import { useI18n } from 'vue-i18n';

const i18n = useI18n();
const gameStore = useGameStore();
const { hintGiverId } = storeToRefs(gameStore);
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message: MessageBase | null) => {
      if (!message) return;
      switch (message.type) {
        case MessageType.RoundStartedMsg:
          handleRoundStarted(message as RoundStartedMessage);
          break;
      }
    });
  }
});

const startRound = () => {
  clientSocket.sendMessage({
    type: MessageType.StartRoundMsg,
    playerId: player.value.id,
  });
}

const handleRoundStarted = (message: RoundStartedMessage) => {
  gameStore.setGameState(GameState.InRound);
  logStore.addLogRecord(
    i18n.t('messages.roundStarted', { name: playerStore.getPlayerName(message.playerId) }),
  );
}

</script>
