<template>
  <div class="menu-container">
    <menu class="menu-content">
      <router-link to="/create">
        <button>{{ $t('components.menu.actions.create') }}</button>
      </router-link>
      <router-link to="/join">
        <button>{{ $t('components.menu.actions.join') }}</button>
      </router-link>
    </menu>
  </div>
</template>

<script setup lang="ts">
</script>
