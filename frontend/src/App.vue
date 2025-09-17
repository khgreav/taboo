<template>
  <h1>Taboo MP!</h1>
  <p>
    Player ID: {{ playerId }}<br>
    Player name: {{ name }}
  </p>
  <button>send</button>
</template>

<script setup lang='ts'>
import { usePlayerStore } from '@/stores/player';
import { useSocketStore } from '@/stores/socketStore';

const { playerId, name } = usePlayerStore();
const clientSocket = useSocketStore();
clientSocket.init('ws://localhost:8081');

clientSocket.$onAction(
  ({ name, after }) => {
    if (name === 'onMessage') {
      after((message: unknown) => {
        console.warn(message);
      })
    }
  }
);
</script>
