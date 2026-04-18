import { apiFetch } from "./api";

type Translations = Record<string, Record<string, string>>;

class I18nManager {
  locale = $state("en-US");
  translations = $state<Translations>({});
  loading = $state(false);

  constructor() {
    if (typeof window !== "undefined") {
      const saved = localStorage.getItem("ggit_locale");
      if (saved) {
        this.locale = saved;
      } else {
        // Detect from browser
        const browserLang = navigator.language;
        if (browserLang) {
            // Basic mapping if needed, e.g. "en" -> "en-US"
            this.locale = browserLang;
        }
      }
    }
  }

  async loadLocale(lang: string) {
    this.loading = true;
    try {
      const res = await apiFetch(`/api/v2/i18n/${lang}`);
      if (res.ok) {
        this.translations = await res.json();
        this.locale = lang;
        if (typeof window !== "undefined") {
          localStorage.setItem("ggit_locale", lang);
          document.documentElement.lang = lang.split('-')[0];
        }
      } else {
          // If 404 and it has a hyphen, try the base language
          if (lang.includes('-')) {
              await this.loadLocale(lang.split('-')[0]);
          }
      }
    } catch (e) {
      console.error("Failed to load locale", lang, e);
    } finally {
      this.loading = false;
    }
  }

  /**
   * t translates a key from the current locale.
   * Keys can be "common.home" or just "home" (defaults to common section).
   */
  t(key: string, ...args: any[]): string {
    let section = "common";
    let k = key;

    if (key.includes('.')) {
      const parts = key.split('.');
      section = parts[0];
      k = parts[1];
    }

    let text = this.translations[section]?.[k] || key;

    // Support for Go-style format verbs %s, %d, %v, %q
    let i = 0;
    return text.replace(/%[sdvq]/g, () => {
      const val = args[i++];
      return val !== undefined ? val : "";
    });
  }

  /**
   * setLocale changes the current language and reloads translations.
   */
  async setLocale(lang: string) {
      await this.loadLocale(lang);
  }
}

export const i18n = new I18nManager();
