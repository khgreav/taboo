import { type App } from 'vue';
import i18n from '@/plugins/i18n';
import registerSocket from '@/plugins/websocket';

export function registerPlugins(app: App) {
  app.use(i18n);
  registerSocket();
}
