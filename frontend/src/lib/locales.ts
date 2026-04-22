import { i18n } from '$lib/i18n.svelte';

/**
 * Mapping of locale codes to their native names.
 * This list is used to generate the language selector.
 */
export const localeNames: Record<string, string> = {
    'bg-BG': 'Български',
    'cs-CZ': 'Čeština',
    'de-DE': 'Deutsch',
    'en-GB': 'English (UK)',
    'en-US': 'English (US)',
    'es-ES': 'Español',
    'fa-IR': 'فارسی',
    'fi-FI': 'Suomi',
    'fr-FR': 'Français',
    'gl-ES': 'Galego',
    'hu-HU': 'Magyar',
    'id-ID': 'Bahasa Indonesia',
    'it-IT': 'Italiano',
    'ja-JP': '日本語',
    'ko-KR': '한국어',
    'lv-LV': 'Latviešu',
    'mn-MN': 'Монгол',
    'nl-NL': 'Nederlands',
    'pl-PL': 'Polski',
    'pt-BR': 'Português (Brasil)',
    'pt-PT': 'Português (Portugal)',
    'ro-RO': 'Română',
    'ru-RU': 'Русский',
    'sk-SK': 'Slovenčina',
    'sr-SP': 'Српски',
    'sv-SE': 'Svenska',
    'tr-TR': 'Türkçe',
    'uk-UA': 'Українська',
    'vi-VN': 'Tiếng Việt',
    'zh-CN': '简体中文',
    'zh-HK': '繁體中文 (香港)',
    'zh-TW': '繁體中文 (台灣)'
};

/**
 * Helper function to get available languages based on JSON files in locales folder.
 */
export function getAvailableLanguages() {
    return Object.entries(localeNames).map(([code, name]) => ({
        code,
        name
    }));
}
