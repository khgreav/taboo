<template>
  <p>
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
  </p>
</template>

<script lang="ts" setup>
import { useGameStore } from '@/stores/gameStore';
import { usePlayerStore } from '@/stores/playerStore';
import { useSocketStore } from '@/stores/socketStore';
import { MessageType} from '@/types/messages';
import { storeToRefs } from 'pinia';

const gameStore = useGameStore();
const { hintGiverId } = storeToRefs(gameStore);
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
const clientSocket = useSocketStore();

const startRound = () => {
  clientSocket.sendMessage({
    type: MessageType.StartRoundMsg,
    playerId: player.value.id,
  });
}

</script>
