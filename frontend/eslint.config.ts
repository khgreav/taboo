import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import pluginVue from 'eslint-plugin-vue'
import vueI18n from '@intlify/eslint-plugin-vue-i18n'
import skipFormatting from '@vue/eslint-config-prettier/skip-formatting'
import { globalIgnores } from 'eslint/config'

export default defineConfigWithVueTs(
  // Global settings for vue-i18n
  {
    settings: {
      'vue-i18n': {
        localeDir: './src/locales/**/*.{json,yaml,yml}',
        messageSyntaxVersion: '^11.0.0',
      },
    },
  },
  // Files to lint
  {
    files: ['**/*.{ts,tsx,vue}'],
    rules: {
      'vue/html-quotes': ['error', 'single'],
      '@intlify/vue-i18n/no-dynamic-keys': 'error',
      '@intlify/vue-i18n/no-unused-keys': [
        'error',
        {
          extensions: ['.js', '.ts', '.vue'],
        },
      ],
    },
  },
  // Ignore build artifacts
  globalIgnores(['**/dist/**', '**/coverage/**']),
  // Vue recommended rules
  pluginVue.configs['flat/recommended'],
  // TypeScript recommended rules
  vueTsConfigs.recommended,
  // Prettier compatibility
  skipFormatting,
  // vue-i18n recommended rules
  ...vueI18n.configs.recommended,
)
