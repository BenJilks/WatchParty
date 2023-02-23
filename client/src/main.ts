import { createApp } from 'vue'
import App from './App.vue'

import './assets/main.css'
import { open_socket_client } from '@/socket_client'

const client_future = open_socket_client()
createApp(App)
    .provide('client_future', client_future)
    .mount('#app');

(window as any).godot_send = (type: string, message: object) =>
    client_future.then(client => client.send(type, message));

(window as any).godot_on = (type: string, callback: (message: object) => void) =>
    client_future.then(client => client.on(type, callback));
