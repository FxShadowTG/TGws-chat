import { createApp } from 'vue'
import App from './App.vue'
import router from '@/router'
import store from '@/store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElIcon from '@element-plus/icons-vue'
import VueCookies from 'vue-cookies'

const app = createApp(App)

Object.keys(ElIcon).forEach((key) => {
    app.component(key, ElIcon[key])
  })

app.use(router).use(store).use(ElementPlus).use(VueCookies).mount('#app')
