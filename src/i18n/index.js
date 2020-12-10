import { createI18n } from 'vue-i18n'

const SUPPORT_LOCALES = ['en', 'zh']

export default createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en',
  fallbackLocale: 'en',
})

export async function loadLocaleMessages(i18n, locale) {
  if (!i18n.global.availableLocales.includes(locale)) {
    if (!SUPPORT_LOCALES.includes(locale)) {
      return false
    }
    const messages = await import( /* webpackChunkName: 'locale-[request]' */ `../../locales/${locale}.json`)
    i18n.global.setLocaleMessage(locale, messages.default)
  }
}
