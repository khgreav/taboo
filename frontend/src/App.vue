<template>
  <h1>Taboo MP!</h1>
  <ConnectName />
  <ClientInfo />
  <NameChange />
  <PlayerList />
</template>

<script setup lang="ts">
import PlayerList from '@/components/PlayerList.vue';
import NameChange from './components/NameChange.vue';
import ClientInfo from './components/ClientInfo.vue';
import ConnectName from './components/ConnectName.vue';

</script>
