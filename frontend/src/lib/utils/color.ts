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
  
  // Adjust for Dark Mode if needed
  let primaryHsl = { ...hsl };
  if (isDark) {
    primaryHsl.s = Math.min(primaryHsl.s, 50);
    primaryHsl.l = Math.max(primaryHsl.l, 75);
  } else {
    primaryHsl.l = Math.min(primaryHsl.l, 50); // Ensure primary is not too light for light mode
  }

  // SECONDARY: Same hue, lower saturation (Neutral-variant)
  const secondaryHsl = { 
    h: primaryHsl.h, 
    s: Math.min(primaryHsl.s * 0.4, 20), 
    l: isDark ? 80 : 40 
  };

  // TERTIARY: Hue shift +60, medium saturation
  const tertiaryHsl = { 
    h: (primaryHsl.h + 60) % 360, 
    s: Math.min(primaryHsl.s * 0.7, 50), 
    l: isDark ? 80 : 40 
  };

  const colors = {
    primary: hslToHex(primaryHsl.h, primaryHsl.s, primaryHsl.l),
    onPrimary: primaryHsl.l > 60 ? '#000000' : '#ffffff',
    primaryContainer: hslToHex(primaryHsl.h, primaryHsl.s * 0.3, isDark ? 25 : 90),
    onPrimaryContainer: hslToHex(primaryHsl.h, primaryHsl.s, isDark ? 90 : 15),

    secondary: hslToHex(secondaryHsl.h, secondaryHsl.s, secondaryHsl.l),
    onSecondary: secondaryHsl.l > 60 ? '#000000' : '#ffffff',
    secondaryContainer: hslToHex(secondaryHsl.h, secondaryHsl.s * 0.5, isDark ? 30 : 92),
    onSecondaryContainer: hslToHex(secondaryHsl.h, secondaryHsl.s, isDark ? 92 : 15),

    tertiary: hslToHex(tertiaryHsl.h, tertiaryHsl.s, tertiaryHsl.l),
    onTertiary: tertiaryHsl.l > 60 ? '#000000' : '#ffffff',
    tertiaryContainer: hslToHex(tertiaryHsl.h, tertiaryHsl.s * 0.5, isDark ? 30 : 92),
    onTertiaryContainer: hslToHex(tertiaryHsl.h, tertiaryHsl.s, isDark ? 92 : 15),
    
    surfaceVariant: hslToHex(primaryHsl.h, Math.min(primaryHsl.s * 0.2, 10), isDark ? 30 : 95),
    onSurfaceVariant: hslToHex(primaryHsl.h, 5, isDark ? 80 : 30),
    outline: hslToHex(primaryHsl.h, 5, isDark ? 60 : 50),
    outlineVariant: hslToHex(primaryHsl.h, 5, isDark ? 40 : 85),
  };

  return colors;
}
