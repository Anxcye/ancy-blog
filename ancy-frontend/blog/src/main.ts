import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style/index.scss'

import App from './App.vue'
import router from './router'
import { createHead } from '@vueuse/head'

const head = createHead()
const app = createApp(App)
app.use(head)
app.use(createPinia())

app.use(router)

app.mount('#app')
