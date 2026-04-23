let currentLocale = $state("en-US");
let dictionary = $state<Record<string, any>>({});
let loading = $state(true);

export const i18n = {
  get locale() {
    return currentLocale;
  },
  get loading() {
    return loading;
  },

  // Función para cambiar de idioma (Carga perezosa)
  async init(lang: string = "en-US") {
    loading = true;
    try {
      // Vite dividirá cada JSON en un archivo (chunk) separado
      const module = await import(`../locales/${lang}.json`);
      dictionary = module.default || module;
      currentLocale = lang;

      if (typeof window !== "undefined") {
        localStorage.setItem("ggit-locale", lang);
      }
    } catch (e) {
      console.error(`No se pudo cargar el idioma: ${lang}`, e);
    } finally {
      loading = false;
    }
  },

  // Función de traducción
  t(key: string, ...args: any[]): string {
    const parts = key.split('.');
    let value: any = dictionary;
    
    for (const part of parts) {
      if (value && typeof value === 'object') {
        value = value[part];
      } else {
        value = undefined;
        break;
      }
    }

    if (typeof value !== 'string') {
      return key;
    }

    // Reemplazar variables estilo Gogs (%s, %d, %v)
    if (args && args.length > 0) {
      let i = 0;
      return value.replace(/%[sdv]/g, (match) => {
        return i < args.length ? String(args[i++]) : match;
      });
    }

    return value;
  },
};
