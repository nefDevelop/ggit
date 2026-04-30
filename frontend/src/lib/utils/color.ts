export interface HSL {
  h: number;
  s: number;
  l: number;
}

export function hexToHsl(hex: string): HSL {
  let r = parseInt(hex.slice(1, 3), 16) / 255;
  let g = parseInt(hex.slice(3, 5), 16) / 255;
  let b = parseInt(hex.slice(5, 7), 16) / 255;

  let max = Math.max(r, g, b), min = Math.min(r, g, b);
  let h = 0, s, l = (max + min) / 2;

  if (max === min) {
    h = s = 0; // achromatic
  } else {
    let d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
    h /= 6;
  }

  return { h: h * 360, s: s * 100, l: l * 100 };
}

export function hslToHex(h: number, s: number, l: number): string {
  l /= 100;
  const a = s * Math.min(l, 1 - l) / 100;
  const f = (n: number) => {
    const k = (n + h / 30) % 12;
    const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
    return Math.round(255 * color).toString(16).padStart(2, '0');
  };
  return `#${f(0)}${f(8)}${f(4)}`;
}

/**
 * Generates an MD3-compliant palette from a seed color.
 * Follows MD3 Tonal Palette rules:
 * - Primary: Chroma 48
 * - Secondary: Chroma 16
 * - Tertiary: Hue + 60, Chroma 24
 * - Neutral: Chroma 4
 * - Neutral Variant: Chroma 8
 */
export function generateMd3Palette(seedHex: string, isDark: boolean = false) {
  const seedHsl = hexToHsl(seedHex);
  const h = seedHsl.h;

  // Helper to map MD3 Tones (0-100) to HSL Lightness
  const getTone = (tone: number) => tone; 

  // PALETTE DEFINITIONS (Saturation approximated for MD3 Chroma)
  const palettes = {
    primary: { h, s: Math.max(seedHsl.s, 48) },
    secondary: { h, s: 16 }, // Official MD3 Tonal Spot Chroma
    tertiary: { h: (h + 60) % 360, s: 24 }, // Official MD3 Tonal Spot Hue +60
    neutral: { h, s: 6 }, 
    neutralVariant: { h, s: 8 },
    error: { h: 0, s: 40 }
  };

  const colors: Record<string, string> = {
    // PRIMARY
    primary: hslToHex(palettes.primary.h, palettes.primary.s, isDark ? 80 : 40),
    onPrimary: hslToHex(palettes.primary.h, palettes.primary.s, isDark ? 20 : 100),
    primaryContainer: hslToHex(palettes.primary.h, palettes.primary.s, isDark ? 30 : 90),
    onPrimaryContainer: hslToHex(palettes.primary.h, palettes.primary.s, isDark ? 90 : 10),

    // SECONDARY
    secondary: hslToHex(palettes.secondary.h, palettes.secondary.s, isDark ? 80 : 40),
    onSecondary: hslToHex(palettes.secondary.h, palettes.secondary.s, isDark ? 20 : 100),
    secondaryContainer: hslToHex(palettes.secondary.h, palettes.secondary.s, isDark ? 30 : 90),
    onSecondaryContainer: hslToHex(palettes.secondary.h, palettes.secondary.s, isDark ? 90 : 10),

    // TERTIARY
    tertiary: hslToHex(palettes.tertiary.h, palettes.tertiary.s, isDark ? 80 : 40),
    onTertiary: hslToHex(palettes.tertiary.h, palettes.tertiary.s, isDark ? 20 : 100),
    tertiaryContainer: hslToHex(palettes.tertiary.h, palettes.tertiary.s, isDark ? 30 : 90),
    onTertiaryContainer: hslToHex(palettes.tertiary.h, palettes.tertiary.s, isDark ? 90 : 10),

    // ERROR
    error: hslToHex(palettes.error.h, palettes.error.s, isDark ? 80 : 40),
    onError: hslToHex(palettes.error.h, palettes.error.s, isDark ? 20 : 100),
    errorContainer: hslToHex(palettes.error.h, palettes.error.s, isDark ? 30 : 90),
    onErrorContainer: hslToHex(palettes.error.h, palettes.error.s, isDark ? 90 : 10),

    // NEUTRAL (Surface/Background) - Bumped tones to avoid "completely black"
    background: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 12 : 99),
    onBackground: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 90 : 12),
    surface: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 14 : 99),
    onSurface: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 90 : 12),
    
    // NEUTRAL VARIANT
    surfaceVariant: hslToHex(palettes.neutralVariant.h, palettes.neutralVariant.s, isDark ? 30 : 90),
    onSurfaceVariant: hslToHex(palettes.neutralVariant.h, palettes.neutralVariant.s, isDark ? 80 : 30),
    outline: hslToHex(palettes.neutralVariant.h, palettes.neutralVariant.s, isDark ? 60 : 50),
    outlineVariant: hslToHex(palettes.neutralVariant.h, palettes.neutralVariant.s, isDark ? 30 : 80),

    // SURFACE CONTAINERS
    surfaceContainerLowest: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 10 : 100),
    surfaceContainerLow: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 15 : 96),
    surfaceContainer: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 18 : 94),
    surfaceContainerHigh: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 22 : 92),
    surfaceContainerHighest: hslToHex(palettes.neutral.h, palettes.neutral.s, isDark ? 28 : 90),
  };

  return colors;
}
