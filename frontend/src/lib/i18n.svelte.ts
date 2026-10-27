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
      // Intentar cargar el idioma solicitado
      // Si es un código corto (ej: 'es'), intentar buscar el más apropiado
      let targetLang = lang;
      
      // Mapeo básico de redundancia para códigos cortos comunes
      const fallbacks: Record<string, string> = {
        'es': 'es-ES',
        'en': 'en-US',
        'zh': 'zh-CN',
        'pt': 'pt-BR',
        'fr': 'fr-FR',
        'de': 'de-DE'
      };

      if (fallbacks[targetLang]) {
        targetLang = fallbacks[targetLang];
      }

      try {
        const module = await import(`../locales/${targetLang}.json`);
        dictionary = module.default || module;
        currentLocale = targetLang;
      } catch (e) {
        console.warn(`No se encontró el idioma específico: ${targetLang}, intentando fallback...`);
        // Si falló el específico, intentar con el genérico si era un código largo
        if (targetLang.includes('-')) {
          const generic = targetLang.split('-')[0];
          if (fallbacks[generic]) {
            const module = await import(`../locales/${fallbacks[generic]}.json`);
            dictionary = module.default || module;
            currentLocale = fallbacks[generic];
          } else {
            throw e; // Re-lanzar si no hay más opciones
          }
        } else {
          throw e;
        }
      }

      if (typeof window !== "undefined") {
        localStorage.setItem("ggit-locale", currentLocale);
      }
    } catch (e) {
      console.error(`Error crítico cargando idioma ${lang}:`, e);
      // Fallback definitivo a inglés si todo falla
      try {
        const module = await import(`../locales/en-US.json`);
        dictionary = module.default || module;
        currentLocale = "en-US";
      } catch (inner) {
        console.error("No se pudo cargar ni siquiera el idioma de fallback (en-US)");
      }
    } finally {
      loading = false;
      if (typeof document !== "undefined") {
        document.documentElement.lang = currentLocale;
      }
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
