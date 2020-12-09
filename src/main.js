import { createApp } from 'vue'
import App from './App.vue'
import store from './store'
import router from './router'
import mixin from './mixin'
import { setupI18n } from './i18n'
import en from '../locales/en'

const i18n = setupI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en',
  fallbackLocale: 'en',
  messages: { en }
})

createApp(App).use(router).use(store).use(i18n).mixin(mixin).mount('#app')
