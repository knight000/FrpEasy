import { createApp } from 'vue'
import { createPinia } from 'pinia'
import vuetify from './plugins/vuetify'
import App from './App.vue'

console.log('Starting Vue app...')

const app = createApp(App)

app.use(createPinia())
app.use(vuetify)

app.mount('#app')

console.log('Vue app mounted successfully')
