<template>
  <div class="menu-container">
    <menu class="menu-content">
      <form>
        <label for="playerNameInput">{{ $t('components.join.name') }}</label>
        <input
          id="playerNameInput"
          v-model="playerName"
          required
        />
      </form>
      <button @click="onSubmit()">{{ $t('components.menu.actions.create') }}</button>
      <router-link to="/">
        <button>{{ $t('components.menu.actions.back') }}</button>
      </router-link>
    </menu>
  </div>
</template>

<script setup lang="ts">
import { GameService } from '@/services/gameService';
import { ref, type Ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const service = new GameService();
const playerName: Ref<string> = ref('');

async function onSubmit() {
  try {
    const data = await service.createGame({ playerName: playerName.value});
    console.warn('Created game with ID:', data.gameId);
    await router.push({
      name: 'Game',
      params: { id: data.gameId },
    });
  } catch {
    // error handling
  }
}
</script>
