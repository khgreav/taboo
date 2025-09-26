import { createApp } from 'vue';
import { createPinia } from 'pinia';

import './styles/index.scss';

import App from './App.vue';
import router from './router';
import { useSocketStore } from './stores/socketStore';

const app = createApp(App);

app.use(createPinia());
app.use(router);

const clientSocket = useSocketStore();
clientSocket.init('ws://localhost:8081/ws');

app.mount('#app');
