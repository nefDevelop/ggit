<script lang="ts">
  /**
   * Identicon.svelte
   * Genera un avatar abstracto determinista basado en un string (seed).
   * Estilo moderno con formas geométricas y colores armoniosos.
   */
  let { 
    seed = '', 
    size = 100,
    type = 'user' // 'user' or 'repo'
  } = $props();

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

  // Generamos formas según el tipo
  const shapes = $derived.by(() => {
    const list = [];
    if (type === 'repo') {
      // Estilo de cuadrícula simétrica para repositorios (estilo Pixel Art)
      // Generamos una matriz de 5x5, pero solo calculamos 3 columnas y espejamos
      const gridColor = palette[hash % palette.length];
      const gridHash = hash.toString(2).padStart(15, '0'); // Necesitamos 15 bits para 3x5
      
      for (let i = 0; i < 5; i++) { // Filas
        for (let j = 0; j < 3; j++) { // Columnas (centro y una mitad)
          const bit = gridHash[i * 3 + j] === '1';
          if (bit) {
            const size = 20;
            // Bloque original
            list.push({ x: j * size, y: i * size, size, color: gridColor });
            // Bloque espejado (si no es el centro)
            if (j < 2) {
              list.push({ x: (4 - j) * size, y: i * size, size, color: gridColor });
            }
          }
        }
      }
    } else {
      // Estilo orgánico para usuarios
      const count = 4;
      for (let i = 0; i < count; i++) {
          const shapeHash = getHash(seed + i);
          const stype = shapeHash % 3; // 0: circle, 1: rect, 2: polygon
          const colorIndex = (shapeHash >> 2) % palette.length;
          const opacity = 0.4 + (shapeHash % 4) * 0.1;
          const x = (shapeHash % 60) + 20;
          const y = ((shapeHash >> 4) % 60) + 20;
          const r = ((shapeHash >> 6) % 30) + 15;
          const rot = (shapeHash >> 8) % 360;
          
          list.push({ type: stype, color: palette[colorIndex], x, y, r, opacity, rot });
      }
    }
    return list;
  });

  const clipId = $derived(`clip-${hash}-${type}`);
  const bgColor = $derived(`rgba(${parseInt(palette[0].slice(1,3), 16)}, ${parseInt(palette[0].slice(3,5), 16)}, ${parseInt(palette[0].slice(5,7), 16)}, 0.15)`);
</script>

<svg
  width={size}
  height={size}
  viewBox="0 0 100 100"
  xmlns="http://www.w3.org/2000/svg"
  class="identicon"
  class:repo-type={type === 'repo'}
  style="background: {bgColor}; max-width: 100%; height: auto;"
>
  <defs>
    <clipPath id={clipId}>
      {#if type === 'repo'}
        <!-- Squircle / Rounded Square para repositorios -->
        <rect x="0" y="0" width="100" height="100" rx="24" />
      {:else}
        <!-- Círculo perfecto para usuarios -->
        <circle cx="50" cy="50" r="50" />
      {/if}
    </clipPath>
  </defs>

  <g clip-path="url(#{clipId})">
    {#if type === 'repo'}
      {#each shapes as s}
        <rect x={s.x} y={s.y} width={s.size} height={s.size} fill={s.color} fill-opacity="0.8" />
      {/each}
    {:else}
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
    {/if}
  </g>
</svg>

<style>
  .identicon {
    display: block;
    border-radius: 50%;
    border: 1px solid rgba(var(--primary-rgb), 0.1);
    overflow: hidden;
    user-select: none;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .identicon.repo-type {
    border-radius: 24%; /* Forma de 'Squircle' para repositorios */
  }
</style>
