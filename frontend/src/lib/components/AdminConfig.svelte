<script lang="ts">
  import { onMount } from 'svelte';
  
  let config = $state<any>(null);
  let loading = $state(true);

  async function fetchConfig() {
    loading = true;
    try {
      const res = await fetch('/api/v2/admin/config');
      if (res.ok) config = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  const sections = [
    { key: 'app', label: 'Application' },
    { key: 'server', label: 'Server' },
    { key: 'database', label: 'Database' },
    { key: 'security', label: 'Security' },
    { key: 'repository', label: 'Repository' },
    { key: 'git', label: 'Git' }
  ];

  onMount(fetchConfig);
</script>

<div class="admin-config">
  <header class="section-header">
    <h2>Server Configuration</h2>
    <p class="subtitle">Read-only view of the <code>app.ini</code> configuration file.</p>
  </header>

  {#if loading}
    <div class="skeleton-config"></div>
  {:else if config}
    <div class="config-sections">
      {#each sections as section}
        <div class="config-card island">
          <h3 class="card-title">{section.label}</h3>
          <div class="config-grid">
            {#each Object.entries(config[section.key] || {}) as [label, value]}
              <div class="config-item">
                <span class="key">{label}</span>
                <span class="val">{typeof value === 'object' ? JSON.stringify(value) : value}</span>
              </div>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .admin-config { display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header h2 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .subtitle { color: var(--text-muted); margin: 0.25rem 0 0 0; }

  .config-sections { display: flex; flex-direction: column; gap: 2rem; }
  .config-card { padding: 0; overflow: hidden; }
  .card-title { padding: 1rem 2rem; margin: 0; font-size: 1rem; background: rgba(255,255,255,0.02); border-bottom: 1px solid var(--border-color); }

  .config-grid { padding: 1.5rem 2rem; display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  .config-item { display: flex; flex-direction: column; gap: 0.2rem; }
  .key { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.6; }
  .val { font-size: 0.95rem; font-family: monospace; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

  .skeleton-config { height: 600px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
