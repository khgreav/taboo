import { type App } from 'vue';
import i18n from '@/plugins/i18n';
import pinia from '@/plugins/pinia';
import toastify, { ToastOptions } from '@/plugins/toast';
import router from '@/router';
import registerSocket from '@/plugins/websocket';

export function registerPlugins(app: App) {
  app
    .use(pinia)
    .use(router)
    .use(i18n)
    .use(toastify, ToastOptions)
  registerSocket();
}
