<template>
  {{ teamTitle }}
  <ul class="player-list">
    <li
      v-for="item in players"
      :key="item.id"
      :class="{ 'current-player': item.id === player.id }"
    >
      {{ item.name }}
      <span
        v-if="gameState === GameState.InLobby"
        :style="{
          color: item.isReady ? 'green' : 'red',
        }"
      >
        {{ $t(item.isReady ? 'components.playerList.states.ready' : 'components.playerList.states.notReady') }}
      </span>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/gameStore';
import { GameState } from '@/types/messages';
import { Team, type OtherPlayer } from '@/types/player';
import { storeToRefs } from 'pinia';
import { computed, type PropType } from 'vue';
import { useI18n } from 'vue-i18n';

const componentProps = defineProps({
  players: {
    type: Array as PropType<OtherPlayer[]>,
    required: true,
  },
  team: {
    type: Number as PropType<Team>,
    required: true,
  },
});
const i18n = useI18n();
const gameStore = useGameStore();
const { player, gameState } = storeToRefs(gameStore);

const teamTitle = computed(() => {
  switch (componentProps.team) {
    case Team.Red:
      return i18n.t('components.playerList.teams.red');
    case Team.Blue:
      return i18n.t('components.playerList.teams.blue');
    default:
      return i18n.t('components.playerList.teams.unassigned');
  }
})

</script>
