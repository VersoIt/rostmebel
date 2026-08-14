import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import { initYandexMetrica } from '@/utils/yandexMetrica';
import './style.css';

const app = createApp(App);
const pinia = createPinia();

initYandexMetrica();

app.use(pinia);
app.use(router);

app.mount('#app');
