import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import { setupPermissionDirective } from './directives/permission'
import i18n from './locales'

const app = createApp(App)
const pinia = createPinia()

app.use(naive)
app.use(pinia)
app.use(router)
app.use(i18n)

// 注册权限指令
setupPermissionDirective(app)

app.mount('#app')
