import { createApp } from 'vue';
import App from './App.vue';
import './style.css';
import 'vfonts/Lato.css';
import naive from 'naive-ui';

const app = createApp(App);
app.use(naive);
app.mount('#app');
