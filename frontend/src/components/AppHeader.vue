<template>
  <div class="header header-border">
    <div class="header-center">
      <span class="header-text">{{ $t('title') }}</span>
    </div>
    <div class="header-right">
      <button
        class="icon-button"
        @click="toggleTheme()"
      >
        <svg v-if="theme === Theme.Light" viewBox="2 2 20 20" width="24" height="24" aria-hidden="true">
          <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
        </svg>
        <svg v-else viewBox="0 0 24 24" width="24" height="24" aria-hidden="true">
          <circle cx="12" cy="12" r="5" />
          <line x1="12" y1="1" x2="12" y2="3" />
          <line x1="12" y1="21" x2="12" y2="23" />
          <line x1="1" y1="12" x2="3" y2="12" />
          <line x1="21" y1="12" x2="23" y2="12" />
          <line x1="4.2" y1="4.2" x2="5.6" y2="5.6" />
          <line x1="18.4" y1="18.4" x2="19.8" y2="19.8" />
          <line x1="4.2" y1="19.8" x2="5.6" y2="18.4" />
          <line x1="18.4" y1="5.6" x2="19.8" y2="4.2" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useTheme } from '@/composables/useTheme';
import { Theme } from '@/types/theme';

const { theme, toggleTheme } = useTheme();
</script>
