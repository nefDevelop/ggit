<script lang="ts">
  /**
   * Identicon.svelte
   * Genera un avatar abstracto determinista basado en un string (seed).
   * Estilo moderno con formas geométricas y colores armoniosos.
   */
  let { seed = '', size = 100 } = $props();

  // Función de hash simple para obtener un número basado en el string
  function getHash(str: string) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    return Math.abs(hash);
  }

  // Generar datos deterministas a partir del seed
  const hash = $derived(getHash(seed));
  
  // Colores vibrantes y modernos
  const palettes = [
    ['#6366f1', '#a855f7', '#ec4899'], // Indigo + Purple + Pink
    ['#3b82f6', '#06b6d4', '#10b981'], // Blue + Cyan + Emerald
    ['#f59e0b', '#ef4444', '#8b5cf6'], // Amber + Red + Violet
    ['#14b8a6', '#0ea5e9', '#6366f1'], // Teal + Sky + Indigo
    ['#d946ef', '#f43f5e', '#fb923c']  // Fuchsia + Rose + Orange
  ];

  const paletteIndex = $derived(hash % palettes.length);
  const palette = $derived(palettes[paletteIndex]);

  // Generamos 3 capas de formas
  const shapes = $derived.by(() => {
    const list = [];
    const count = 4;
    for (let i = 0; i < count; i++) {
        const shapeHash = getHash(seed + i);
        const type = shapeHash % 3; // 0: circle, 1: rect, 2: polygon
        const colorIndex = (shapeHash >> 2) % palette.length;
        const opacity = 0.4 + (shapeHash % 4) * 0.1;
        const x = (shapeHash % 60) + 20;
        const y = ((shapeHash >> 4) % 60) + 20;
        const r = ((shapeHash >> 6) % 30) + 15;
        const rot = (shapeHash >> 8) % 360;
        
        list.push({ type, color: palette[colorIndex], x, y, r, opacity, rot });
    }
    return list;
  });

  const bgColor = $derived(`rgba(${parseInt(palette[0].slice(1,3), 16)}, ${parseInt(palette[0].slice(3,5), 16)}, ${parseInt(palette[0].slice(5,7), 16)}, 0.15)`);
</script>

<svg
  width={size}
  height={size}
  viewBox="0 0 100 100"
  xmlns="http://www.w3.org/2000/svg"
  class="identicon"
  style="background: {bgColor};"
>
  <defs>
    <clipPath id="circle-clip">
      <circle cx="50" cy="50" r="50" />
    </clipPath>
  </defs>

  <g clip-path="url(#circle-clip)">
    {#each shapes as s}
      {#if s.type === 0}
        <circle cx={s.x} cy={s.y} r={s.r} fill={s.color} fill-opacity={s.opacity} />
      {:else if s.type === 1}
        <rect 
            x={s.x - s.r} y={s.y - s.r} width={s.r * 2} height={s.r * 2} 
            fill={s.color} fill-opacity={s.opacity} 
            transform="rotate({s.rot} {s.x} {s.y})"
            rx="8"
        />
      {:else}
        <circle cx={s.x} cy={s.y} r={s.r} fill="none" stroke={s.color} stroke-width="8" stroke-opacity={s.opacity} />
      {/if}
    {/each}
  </g>
</svg>

<style>
  .identicon {
    display: block;
    border-radius: 50%;
    /* Permite que el borde heredecurrentColor si es necesario */
    border: 1px solid currentColor;
    overflow: hidden;
    user-select: none;
  }
</style>
