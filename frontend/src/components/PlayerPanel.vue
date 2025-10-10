<template>
  <label>{{ $t('components.playerList.title') }}</label>
  <div>
    <PlayerList
      :players="redPlayers"
      :team="Team.Red"
    />
    <PlayerList
      :players="bluePlayers"
      :team="Team.Blue"
    />
    <PlayerList
      v-if="unassignedPlayers.length > 0"
      :players="unassignedPlayers"
      :team="Team.Unassigned"
    />
    <button
      v-if="gameState === GameState.InLobby && player.team !== Team.Unassigned"
      @click="changeReadyState()"
    >
      {{ readyButtonText }}
    </button>
    <button
      v-if="gameState === GameState.InLobby && player.team !== Team.Red && redPlayers.length < 2"
      @click="changeTeam(Team.Red)"
    >
      {{ $t('components.playerList.actions.red') }}
    </button>
    <button
      v-if="gameState === GameState.InLobby && player.team !== Team.Blue && bluePlayers.length < 2"
      @click="changeTeam(Team.Blue)"
    >
      {{ $t('components.playerList.actions.blue') }}
    </button>
    <button
      v-if="gameState === GameState.InLobby && player.team !== Team.Unassigned"
      @click="changeTeam(Team.Unassigned)"
    >
      {{ $t('components.playerList.actions.unassigned') }}
    </button>
  </div>
</template>

<script lang="ts" setup>
import { useGameStore } from '@/stores/gameStore';
import { useLogStore } from '@/stores/logStore';
import { useSocketStore } from '@/stores/socketStore';
import {
  GameState,
  MessageType,
  type PlayerJoinedMessage,
  type PlayerLeftMessage,
  type PlayerListMessage,
  type PlayerReadyMessage,
  type TeamChangedMessage,
} from '@/types/messages';
import { Team, type OtherPlayer } from '@/types/player';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import PlayerList from './PlayerList.vue';
import { usePlayerStore } from '@/stores/playerStore';
import { teamToString } from '@/utils/team';

const i18n = useI18n();
const playerStore = usePlayerStore();
const { player, playerList } = storeToRefs(playerStore);
const gameStore = useGameStore();
const { gameState } = storeToRefs(gameStore);
const logStore = useLogStore();
const clientSocket = useSocketStore();

clientSocket.$onAction(({ name, after }) => {
  if (name === 'onMessage') {
    after((message) => {
      if (!message) return;
      switch (message.type) {
        case MessageType.PlayerJoinedMsg:
          handlePlayerJoined(message as PlayerJoinedMessage);
          break;
        case MessageType.PlayerListMsg:
          handlePlayerList(message as PlayerListMessage);
          break;
        case MessageType.PlayerLeftMsg:
          handlePlayerLeft(message as PlayerLeftMessage);
          break;
        case MessageType.TeamChangedMsg:
          handleTeamChanged(message as TeamChangedMessage);
          break;
        case MessageType.PlayerReadyMsg:
          handlePlayerReady(message as PlayerReadyMessage);
          break;
      }
    });
  }
});

const unassignedPlayers = computed(() => {
  const players: OtherPlayer[] = [];
  if (player.value.team === Team.Unassigned) {
    players.push(player.value as OtherPlayer);
  }
  for (const player of playerList.value) {
    if (player.team === Team.Unassigned) {
      players.push(player);
    }
  }
  return players;
});

const redPlayers = computed(() => {
  const players: OtherPlayer[] = [];
  if (player.value.team === Team.Red) {
    players.push(player.value as OtherPlayer);
  }
  for (const player of playerList.value) {
    if (player.team === Team.Red) {
      players.push(player);
    }
  }
  return players;
});

const bluePlayers = computed(() => {
  const players: OtherPlayer[] = [];
  if (player.value.team === Team.Blue) {
    players.push(player.value as OtherPlayer);
  }
  for (const player of playerList.value) {
    if (player.team === Team.Blue) {
      players.push(player);
    }
  }
  return players;
});

const readyButtonText = computed(() => {
  if (player.value.isReady) {
    return i18n.t('components.playerList.actions.unready');
  }
  return i18n.t('components.playerList.actions.ready');
});

const changeTeam = (team: Team) => {
  clientSocket.sendMessage({
    type: MessageType.ChangeTeamMsg,
    playerId: player.value.id,
    team: team,
  });
};

const handleTeamChanged = (message: TeamChangedMessage) => {
  playerStore.setPlayerTeam(message.playerId, message.team)
  const playerName = playerStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t('messages.teamChanged', { name: playerName, team: teamToString(message.team, i18n) }),
  );
};

const changeReadyState = () => {
  clientSocket.sendMessage({
    type: MessageType.PlayerReadyMsg,
    playerId: player.value.id,
    isReady: !player.value.isReady,
  });
};

const handlePlayerList = (message: PlayerListMessage) => {
  playerStore.setPlayerList(message.players);
};

const handlePlayerJoined = (message: PlayerJoinedMessage) => {
  playerStore.addPlayer({
    id: message.playerId,
    name: message.name,
    team: Team.Unassigned,
    isReady: false,
  });

  logStore.addLogRecord(
    i18n.t('messages.playerJoined', { name: message.name }),
  );
};

const handlePlayerLeft = (message: PlayerLeftMessage) => {
  playerStore.removePlayer(message.playerId);
  const playerName = playerStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t('messages.playerLeft', { name: playerName }),
  );
};

const handlePlayerReady = (message: PlayerReadyMessage) => {
  playerStore.setPlayerReady(message.playerId, message.isReady);
  const playerName = playerStore.getPlayerName(message.playerId);
  if (!playerName) {
    return;
  }
  logStore.addLogRecord(
    i18n.t(message.isReady ? 'messages.ready' : 'messages.notReady', { name: playerName }),
  );
};
</script>
