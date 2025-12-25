import { Theme } from '@/types/theme';
import { ref, type Ref } from 'vue';

export function useTheme() {
  const theme: Ref<Theme> = ref<Theme>(Theme.Light);

  const toggleTheme = () => {
    const isDark = document.documentElement.classList.toggle(Theme.Dark);
    theme.value = isDark ? Theme.Dark : Theme.Light;
    localStorage.setItem('taboo-app-theme', theme.value);
  };

  const preferredTheme = localStorage.getItem('taboo-app-theme') as Theme | null;
  if (preferredTheme) {
    theme.value = preferredTheme;
    document.documentElement.classList.toggle(Theme.Dark, preferredTheme === Theme.Dark);
  } else {
    const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    document.documentElement.classList.toggle(Theme.Dark, prefersDark);
  }

  return {
    theme,
    toggleTheme,
  };
}
