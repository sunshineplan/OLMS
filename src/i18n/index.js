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
  if (!SUPPORT_LOCALES.includes(locale))
    locale = i18n.global.fallbackLocale.value
  if (!i18n.global.availableLocales.includes(locale))
    i18n.global.setLocaleMessage(locale, await import(/* webpackChunkName: 'locale-[request]' */ `../../locales/${locale}.json`))
  i18n.global.locale.value = locale
}

export default i18n
