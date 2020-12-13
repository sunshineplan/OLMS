import { createI18n } from 'vue-i18n'

const SUPPORT_LOCALES = ['en', 'zh']

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {}
})

export async function loadLocaleMessages(locale) {
  if (!i18n.global.availableLocales.includes(locale)) {
    if (!SUPPORT_LOCALES.includes(locale)) {
      return false
    }
    i18n.global.setLocaleMessage(locale, await import(/* webpackChunkName: 'locale-[request]' */ `../../locales/${locale}.json`))
  }
}

export default i18n
