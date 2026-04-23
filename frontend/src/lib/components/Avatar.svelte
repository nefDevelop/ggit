<script lang="ts">
  import Identicon from './Identicon.svelte';
  
  /**
   * Avatar.svelte
   * Componente universal para avatares en Ggit.
   * Si 'src' es nulo o vacío, muestra un Identicon determinista.
   */
  let { 
    src = '', 
    username = '?', 
    seed = '', 
    size = 40, 
    class: className = '' 
  } = $props();

  // El seed por defecto es el username si no se pasa nada más
  const finalSeed = $derived(seed || username);
  
  // Determinamos si es un avatar real o debemos generar uno
  const useIdenticon = $derived(!src || src === '');
</script>

<div class="avatar-container {className}" style="width: {size}px; height: {size}px;">
  {#if useIdenticon}
    <Identicon seed={finalSeed} size={size} />
  {:else}
    <img src={src} alt={username} class="avatar-img" />
  {/if}
</div>

<style>
  .avatar-container {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    overflow: hidden;
    flex-shrink: 0;
    aspect-ratio: 1/1;
  }

  .avatar-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }
</style>
