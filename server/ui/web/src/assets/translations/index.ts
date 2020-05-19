import VueI18n from "vue-i18n";
import { ENGLISH_TRANSLATION } from "./en";
import { FRENCH_TRANSLATION } from "./fr";

export const TRANSLATOR_OPTIONS: VueI18n.I18nOptions = {
  locale: "en",
  fallbackLocale: "en",
  messages: {
    en: ENGLISH_TRANSLATION,
    fr: FRENCH_TRANSLATION
  }
};
