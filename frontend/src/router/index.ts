import DefaultLayout from '@/layouts/DefaultLayout.vue';
import CreateGame from '@/pages/CreateGame.vue';
import GameView from '@/pages/GameView.vue';
import JoinGame from '@/pages/JoinGame.vue';
import MainMenu from '@/pages/MainMenu.vue';
import NotFound from '@/pages/NotFound.vue';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: DefaultLayout,
      children: [
        {
          path: '',
          name: 'Menu',
          component: MainMenu,
        },
        {
          path: 'create',
          name: 'CreateGame',
          component: CreateGame,
        },
        {
          path: 'join',
          name: 'JoinGame',
          component: JoinGame,
        },
        {
          path: 'game/:id',
          name: 'Game',
          component: GameView,
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: NotFound,
    }
  ],
});

export default router;
