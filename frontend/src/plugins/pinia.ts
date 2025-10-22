import { createPinia } from 'pinia';
import { createPersistedState } from 'pinia-plugin-persistedstate';

const pinia = createPinia();
pinia.use(createPersistedState({
  storage: import.meta.env.DEV ? sessionStorage : localStorage,
}));

export default pinia;
