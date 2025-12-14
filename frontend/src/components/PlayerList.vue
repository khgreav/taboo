<template>
  <div
    :class="{
      'team': true,
      'last-team': team === Team.Blue && gameState !== GameState.InLobby,
    }"
  >
    <div class="team-title">
      {{ teamTitle }} <span v-if="team !== Team.Unassigned">{{ teamScore }}</span>
    </div>
    <div class="team-players">
      <ul>
        <li
          v-for="item in players"
          :key="item.id"
          :class="{
            'current-player': item.id === player.id,
            'disconnected-player': !item.connected,
          }"
        >
          {{ item.name }}
          <span
            v-if="gameState === GameState.InLobby && item.team !== Team.Unassigned"
            :style="{
              color: item.isReady ? 'green' : 'red',
            }"
          >
            {{ item.isReady ? $t('components.playerList.states.ready') : $t('components.playerList.states.notReady') }}
          </span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useGameStore } from '@/stores/gameStore';
import { usePlayerStore } from '@/stores/playerStore';
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
const playerStore = usePlayerStore();
const { player } = storeToRefs(playerStore);
const gameStore = useGameStore();
const { gameState, redScore, blueScore } = storeToRefs(gameStore);

const teamScore = computed(() => {
  if (componentProps.team === Team.Red) {
    return redScore.value;
  } else {
    return blueScore.value;
  }
})

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
