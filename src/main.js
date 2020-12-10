import { createApp } from 'vue'
import App from './App.vue'
import store from './store'
import router from './router'
import i18n from './i18n'
import mixin from './mixin'

createApp(App).use(router).use(store).use(i18n).mixin(mixin).mount('#app')
