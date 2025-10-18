<template>
  <div class="menu-container">
    <menu class="menu-content">
      <form>
        <label for="gameIdInput">{{ $t('components.join.gameId') }}</label>
        <input
          id="gameIdInput"
          v-model="gameId"
          required
        >
        <label for="playerNameInput">{{ $t('components.join.name') }}</label>
        <input
          id="playerNameInput"
          v-model="playerName"
          required
        />
      </form>
      <button>{{ $t('components.menu.actions.join') }}</button>
      <router-link to="/">
        <button>{{ $t('components.menu.actions.back') }}</button>
      </router-link>
    </menu>
  </div>
</template>

<script setup lang="ts">
import { ref, type Ref } from 'vue';

const gameId: Ref<string> = ref('');
const playerName: Ref<string> = ref('');
</script>
