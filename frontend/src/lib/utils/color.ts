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

export function generateMd3Palette(seedHex: string, isDark: boolean = false) {
  const hsl = hexToHsl(seedHex);
  
  // PRIMARY
  let p_h = hsl.h;
  let p_s = isDark ? Math.min(hsl.s, 40) : Math.min(hsl.s, 100);
  let p_l = isDark ? 80 : 40;

  // SECONDARY (Neutral-variant)
  let s_h = hsl.h;
  let s_s = Math.min(hsl.s * 0.25, 16);
  let s_l = isDark ? 80 : 45;

  // TERTIARY (Accent)
  let t_h = (hsl.h + 60) % 360;
  let t_s = Math.min(hsl.s * 0.5, 45);
  let t_l = isDark ? 80 : 45;

  // NEUTRAL
  let n_h = hsl.h;
  let n_s = Math.min(hsl.s * 0.08, 6);

  const colors: Record<string, string> = {
    // Core MD3
    primary: hslToHex(p_h, p_s, p_l),
    onPrimary: p_l > 60 ? '#001d31' : '#ffffff',
    primaryContainer: hslToHex(p_h, p_s * 0.4, isDark ? 30 : 90),
    onPrimaryContainer: hslToHex(p_h, p_s, isDark ? 90 : 15),

    secondary: hslToHex(s_h, s_s, s_l),
    onSecondary: s_l > 60 ? '#1b1b1f' : '#ffffff',
    secondaryContainer: hslToHex(s_h, s_s, isDark ? 30 : 90),
    onSecondaryContainer: hslToHex(s_h, s_s, isDark ? 90 : 15),

    tertiary: hslToHex(t_h, t_s, t_l),
    onTertiary: t_l > 60 ? '#1b1b1f' : '#ffffff',
    tertiaryContainer: hslToHex(t_h, t_s, isDark ? 30 : 90),
    onTertiaryContainer: hslToHex(t_h, t_s, isDark ? 90 : 15),

    error: isDark ? '#ffb4ab' : '#ba1a1a',
    onError: isDark ? '#690005' : '#ffffff',
    errorContainer: isDark ? '#93000a' : '#ffdad6',
    onErrorContainer: isDark ? '#ffdad6' : '#410002',

    background: isDark ? '#0d0e11' : '#f8f9ff',
    onBackground: isDark ? '#e2e2e6' : '#1b1b1f',
    surface: isDark ? '#0d0e11' : '#f8f9ff',
    onSurface: isDark ? '#e2e2e6' : '#1b1b1f',
    
    surfaceVariant: hslToHex(n_h, n_s * 2, isDark ? 30 : 90),
    onSurfaceVariant: isDark ? '#c4c6d0' : '#44474e',
    
    outline: isDark ? '#8e9099' : '#757780',
    outlineVariant: isDark ? '#44474e' : '#c5c6d0',

    surfaceContainerLowest: isDark ? '#0d0e11' : '#ffffff',
    surfaceContainerLow: isDark ? '#1a1c1e' : '#f7f3f8',
    surfaceContainer: isDark ? '#1b1d21' : '#ffffff',
    surfaceContainerHigh: isDark ? '#252a2f' : '#f2f2f7',
    surfaceContainerHighest: isDark ? '#30353b' : '#ebebed',
  };

  return colors;
}
