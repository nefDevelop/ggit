<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  import BarChart3 from "lucide-svelte/icons/bar-chart-3";
  import Users from "lucide-svelte/icons/users";
  import Building2 from "lucide-svelte/icons/building-2";
  import Book from "lucide-svelte/icons/book";
  import Key from "lucide-svelte/icons/key";
  import Star from "lucide-svelte/icons/star";
  import Bell from "lucide-svelte/icons/bell";
  import GitMerge from "lucide-svelte/icons/git-merge";
  import Eye from "lucide-svelte/icons/eye";
  import CloudOff from "lucide-svelte/icons/cloud-off";
  import RefreshCw from "lucide-svelte/icons/refresh-cw";

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

<Card class="settings-premium-card" padding="none">
  <PageHeader 
    title="Panel de Estadísticas"
    subtitle="Resumen global de objetos y uso del sistema"
    icon={BarChart3}
  />

  {#if loading}
    <div class="stats-grid">
      {#each Array(8) as _}
        <div class="stat-premium-card skeleton-card">
          <div class="card-icon skeleton-pulse"></div>
          <div class="card-data">
            <div class="skeleton-pulse val-skeleton"></div>
            <div class="skeleton-pulse lab-skeleton"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else if stats}
    <div class="stats-grid">
      <div class="stat-premium-card island tonal-hover role-primary">
        <div class="card-icon">
          <Users size={32} />
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.User}</span>
          <span class="label">Usuarios</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover role-secondary">
        <div class="card-icon">
          <Building2 size={32} />
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Org}</span>
          <span class="label">Organizaciones</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover role-tertiary">
        <div class="card-icon">
          <Book size={32} />
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Repo}</span>
          <span class="label">Repositorios</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover role-primary">
        <div class="card-icon">
          <Key size={32} />
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.PublicKey}</span>
          <span class="label">Claves SSH</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover role-tertiary">
        <div class="card-icon">
          <Star size={32} />
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Star}</span>
          <span class="label">Favoritos</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover role-error">
        <div class="card-icon">
          <Bell size={32} />
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Notice}</span>
          <span class="label">Avisos</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover role-secondary">
        <div class="card-icon">
          <GitMerge size={32} />
        </div>
        <div class="card-data">
          <span class="value">{stats.Counter.Attachment}</span>
          <span class="label">Adjuntos</span>
        </div>
      </div>

      <div class="stat-premium-card island tonal-hover role-primary">
        <div class="card-icon">
          <Eye size={32} />
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
         <CloudOff size={48} />
      </div>
      <div class="error-text">
        <h3>Error de Conexión</h3>
        <p>No se han podido recuperar las estadísticas del sistema. Verifica que tienes privilegios de administrador y que el backend está respondiendo correctamente.</p>
        <Button variant="primary" icon={RefreshCw} onclick={fetchStats}>Reintentar Conexión</Button>
      </div>
    </div>
  {/if}
</Card>

<style>




  .stats-grid { 
    display: grid; 
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); 
    gap: 1.5rem; 
    animation: fadeIn 0.4s ease-out;
  }
  
  .stat-premium-card { display: flex; align-items: center; gap: 1.5rem; padding: 2rem; border-radius: 1.25rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); background: rgba(255,255,255,0.01); }
  .stat-premium-card:hover { transform: translateY(-8px) scale(1.02); background: rgba(255,255,255,0.03); box-shadow: var(--shadow-xl); }
  
  .stat-premium-card.role-primary:hover { border-color: var(--primary-color); }
  .stat-premium-card.role-secondary:hover { border-color: var(--secondary); }
  .stat-premium-card.role-tertiary:hover { border-color: var(--tertiary); }
  .stat-premium-card.role-error:hover { border-color: var(--error); }

  .stat-premium-card.role-primary .card-icon { color: var(--primary-color); background: var(--md-sys-color-primary-container); }
  .stat-premium-card.role-secondary .card-icon { color: var(--secondary); background: var(--md-sys-color-secondary-container); }
  .stat-premium-card.role-tertiary .card-icon { color: var(--tertiary); background: var(--md-sys-color-tertiary-container); }
  .stat-premium-card.role-error .card-icon { color: var(--error); background: var(--md-sys-color-error-container); }

  .card-icon { width: 64px; height: 64px; border-radius: 0.75rem; display: flex; align-items: center; justify-content: center; border: 1px solid var(--md-sys-color-outline-variant); }

  .card-data { display: flex; flex-direction: column; gap: 0.1rem; }
  .card-data .value { font-size: 2rem; font-weight: 950; color: var(--text-color); letter-spacing: -0.02em; }
  .card-data .label { font-size: 0.8rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.1em; }

  .premium-error-zone { padding: 4rem; display: flex; align-items: center; gap: 3rem; background: rgba(var(--error-rgb), 0.05); border: 2px dashed var(--md-sys-color-outline-variant); border-radius: 1.5rem; }
  .error-visual { width: 100px; height: 100px; border-radius: 0.75rem; background: rgba(var(--error-rgb), 0.1); color: var(--error); display: flex; align-items: center; justify-content: center; }
  .error-text h3 { margin: 0 0 0.75rem 0; font-size: 1.75rem; font-weight: 900; }
  .error-text p { margin: 0 0 2rem 0; color: var(--text-muted); font-size: 1.1rem; max-width: 500px; line-height: 1.5; }

  .skeleton-card { border-color: var(--md-sys-color-outline-variant); background: rgba(255,255,255,0.01); pointer-events: none; }
  .skeleton-pulse { background: var(--md-sys-color-surface-variant); animation: pulse 1.5s infinite ease-in-out; border-radius: 0.5rem; }
  .val-skeleton { height: 2rem; width: 4rem; margin-bottom: 0.5rem; }
  .lab-skeleton { height: 0.8rem; width: 6rem; }

  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.7; } 100% { opacity: 0.3; } }

  @media (max-width: 600px) {
    .premium-error-zone { flex-direction: column; text-align: center; padding: 2rem; }
    .error-text p { max-width: 100%; }
  }
</style>
