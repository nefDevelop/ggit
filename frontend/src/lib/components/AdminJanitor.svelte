<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { tool } from "$lib/tool";
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";

  let loading = $state(true);
  let stats = $state<any[]>([]);

  async function fetchStats() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/admin/janitor");
      if (res.ok) {
        stats = res.data || [];
        stats.sort((a, b) => b.waste - a.waste);
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  onMount(fetchStats);
</script>

<div class="admin-janitor">
  <div class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">mop</span>
      <div>
        <h2>Ggit Janitor: Global Health</h2>
        <p class="subtitle">Monitorea el desperdicio de espacio y la salud de todos los repositorios.</p>
      </div>
    </div>
  </div>

  {#if loading}
    <div class="skeleton-grid">
      {#each Array(3) as _}
        <div class="skeleton-card"></div>
      {/each}
    </div>
  {:else if stats.length === 0}
    <Card class="empty-card">
      <span class="material-symbols-rounded">check_circle</span>
      <p>Todos los repositorios están limpios o no se han escaneado todavía.</p>
    </Card>
  {:else}
    <div class="stats-grid">
      {#each stats as repo}
        <Card class="repo-stat-card" tonal={repo.waste > 0}>
          <div class="repo-info">
            <div class="repo-main">
              <span class="material-symbols-rounded repo-icon">book</span>
              <div class="repo-names">
                <span class="owner">{repo.owner} /</span>
                <span class="name">{repo.name}</span>
              </div>
            </div>
            <div class="waste-badge" class:warning={repo.waste > 0}>
              <span class="material-symbols-rounded">delete</span>
              {tool.fileSize(repo.waste)}
            </div>
          </div>

          <div class="repo-details">
            <div class="detail-item">
              <span class="label">Problemas</span>
              <span class="value">{repo.issues_count}</span>
            </div>
            <div class="detail-item">
              <span class="label">Último Escaneo</span>
              <span class="value">{new Date(repo.last_check * 1000).toLocaleDateString()}</span>
            </div>
          </div>

          <div class="card-actions">
            <Button href="/{repo.owner}/{repo.name}/settings" variant="outline" size="sm" icon="settings">
              Gestionar
            </Button>
          </div>
        </Card>
      {/each}
    </div>
  {/if}
</div>

<style>
  .admin-janitor { display: flex; flex-direction: column; gap: 2.5rem; }
  
  .section-header { margin-bottom: 1rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .stats-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.5rem; }
  
  .repo-stat-card { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .repo-info { display: flex; justify-content: space-between; align-items: flex-start; }
  .repo-main { display: flex; gap: 1rem; align-items: center; }
  .repo-icon { color: var(--primary-color); opacity: 0.7; }
  .repo-names { display: flex; flex-direction: column; }
  .owner { font-size: 0.8rem; color: var(--text-muted); font-weight: 700; }
  .name { font-size: 1.2rem; font-weight: 900; }

  .waste-badge { display: flex; align-items: center; gap: 0.5rem; padding: 0.5rem 1rem; border-radius: 1rem; background: var(--surface-variant); font-weight: 800; font-size: 0.9rem; }
  .waste-badge.warning { background: rgba(var(--error-rgb), 0.1); color: var(--error); }
  .waste-badge span { font-size: 1.1rem; }

  .repo-details { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; padding: 1rem; background: rgba(0,0,0,0.1); border-radius: 1rem; }
  .detail-item { display: flex; flex-direction: column; gap: 0.2rem; }
  .label { font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); font-weight: 800; }
  .value { font-size: 0.95rem; font-weight: 700; }

  .card-actions { display: flex; justify-content: flex-end; }

  .empty-card { padding: 4rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-card span { font-size: 4rem; opacity: 0.2; }

  .skeleton-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.5rem; }
  .skeleton-card { height: 200px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }

  @keyframes pulse { 0% { opacity: 1; } 50% { opacity: 0.5; } 100% { opacity: 1; } }
</style>
