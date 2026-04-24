<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let stats = $state<any>(null);
  let loading = $state(true);

  async function fetchStats() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/admin/stats");
      if (res.ok) stats = await res.json();
    } catch (e) { console.error(e); } finally { loading = false; }
  }

  onMount(fetchStats);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">analytics</span>
      <div>
        <h2>Panel de Estadísticas</h2>
        <p class="subtitle">Resumen global de objetos y uso del sistema</p>
      </div>
    </div>
  </header>

  {#if loading}
    <div class="skeleton-dashboard"></div>
  {:else if stats}
    <div class="stats-grid">
      <div class="stat-premium-card island tonal-hover">
        <div class="card-icon" style="--color: #3b82f6">
          <span class="material-symbols-rounded">group</span>
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.User}</span>
          <span class="label">Usuarios</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover">
        <div class="card-icon" style="--color: #10b981">
          <span class="material-symbols-rounded">corporate_fare</span>
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Org}</span>
          <span class="label">Organizaciones</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover">
        <div class="card-icon" style="--color: #f59e0b">
          <span class="material-symbols-rounded">book</span>
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Repo}</span>
          <span class="label">Repositorios</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover">
        <div class="card-icon" style="--color: #9c27b0">
          <span class="material-symbols-rounded">vpn_key</span>
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.PublicKey}</span>
          <span class="label">Claves SSH</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover">
        <div class="card-icon" style="--color: #06b6d4">
          <span class="material-symbols-rounded">star</span>
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Star}</span>
          <span class="label">Favoritos</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover">
        <div class="card-icon" style="--color: #f43f5e">
          <span class="material-symbols-rounded">notifications</span>
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Notice}</span>
          <span class="label">Avisos</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover">
        <div class="card-icon" style="--color: #8b5cf6">
          <span class="material-symbols-rounded">merge_type</span>
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Attachment}</span>
          <span class="label">Adjuntos</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover">
        <div class="card-icon" style="--color: #ec4899">
          <span class="material-symbols-rounded">visibility</span>
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Watch}</span>
          <span class="label">Seguimientos</span>
        </div>
      </div>
    </div>
  {:else}
    <div class="premium-error-zone island">
      <div class="error-visual">
         <span class="material-symbols-rounded">cloud_off</span>
      </div>
      <div class="error-text">
        <h3>Error de Conexión</h3>
        <p>No se han podido recuperar las estadísticas del sistema. Verifica que tienes privilegios de administrador y que el backend está respondiendo correctamente.</p>
        <Button variant="primary" icon="refresh" onclick={fetchStats}>Reintentar Conexión</Button>
      </div>
    </div>
  {/if}
</Card>

<style>
  :global(.settings-premium-card) { border-radius: 3rem; background: rgba(var(--surface-rgb), 0.5); backdrop-filter: blur(24px); border: 1px solid var(--md-sys-color-outline-variant); }
  .section-header { margin-bottom: 4rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  @media (max-width: 600px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { font-size: 0.9rem; }
  }

  .stats-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1.5rem; }
  
  .stat-premium-card { display: flex; align-items: center; gap: 1.5rem; padding: 2rem; border-radius: 2.5rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); background: rgba(255,255,255,0.01); }
  .stat-premium-card:hover { transform: translateY(-8px) scale(1.02); background: rgba(255,255,255,0.03); border-color: var(--primary-color); box-shadow: var(--shadow-xl); }

  .card-icon { width: 64px; height: 64px; border-radius: 1.5rem; background: rgba(255,255,255,0.03); color: var(--color); display: flex; align-items: center; justify-content: center; border: 1px solid rgba(255,255,255,0.05); }
  .card-icon span { font-size: 2.5rem; }

  .card-data { display: flex; flex-direction: column; gap: 0.1rem; }
  .card-data .value { font-size: 2rem; font-weight: 950; color: var(--text-color); letter-spacing: -0.02em; }
  .card-data .label { font-size: 0.8rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.1em; }

  .premium-error-zone { padding: 4rem; display: flex; align-items: center; gap: 3rem; background: rgba(var(--error-rgb), 0.05); border: 2px dashed var(--md-sys-color-outline-variant); border-radius: 3rem; }
  .error-visual { width: 100px; height: 100px; border-radius: 50%; background: rgba(var(--error-rgb), 0.1); color: var(--error); display: flex; align-items: center; justify-content: center; }
  .error-visual span { font-size: 4rem; }
  .error-text h3 { margin: 0 0 0.75rem 0; font-size: 1.75rem; font-weight: 900; }
  .error-text p { margin: 0 0 2rem 0; color: var(--text-muted); font-size: 1.1rem; max-width: 500px; line-height: 1.5; }

  .skeleton-dashboard { height: 500px; background: var(--surface-variant); border-radius: 3rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }

  @media (max-width: 600px) {
    .premium-error-zone { flex-direction: column; text-align: center; padding: 2rem; }
    .error-text p { max-width: 100%; }
  }
</style>
