<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let stats = $state<any>(null);
  let loading = $state(true);

  async function fetchStats() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/admin/stats');
      if (res.ok) {
        stats = await res.json();
      } else {
        stats = null;
      }
    } catch (e) {
      console.error(e);
      stats = null;
    } finally {
      loading = false;
    }
  }

  onMount(fetchStats);
</script>

<div class="admin-dashboard">
  <header class="section-header">
    <h2>System Statistics</h2>
    <p class="subtitle">Overview of current system usage and objects.</p>
  </header>

  {#if loading}
    <div class="stats-grid">
      {#each Array(8) as _}
        <div class="stat-card island skeleton"></div>
      {/each}
    </div>
  {:else if stats}
    <div class="stats-grid">
      <div class="stat-card island tonal">
        <span class="material-symbols-rounded icon">group</span>
        <div class="stat-content">
          <span class="label">Users</span>
          <span class="value">{stats.Counter.User}</span>
        </div>
      </div>
      <div class="stat-card island tonal">
        <span class="material-symbols-rounded icon">corporate_fare</span>
        <div class="stat-content">
          <span class="label">Organizations</span>
          <span class="value">{stats.Counter.Org}</span>
        </div>
      </div>
      <div class="stat-card island tonal">
        <span class="material-symbols-rounded icon">vpn_key</span>
        <div class="stat-content">
          <span class="label">Public Keys</span>
          <span class="value">{stats.Counter.PublicKey}</span>
        </div>
      </div>
      <div class="stat-card island tonal">
        <span class="material-symbols-rounded icon">book</span>
        <div class="stat-content">
          <span class="label">Repositories</span>
          <span class="value">{stats.Counter.Repo}</span>
        </div>
      </div>
      <div class="stat-card island tonal">
        <span class="material-symbols-rounded icon">star</span>
        <div class="stat-content">
          <span class="label">Stars</span>
          <span class="value">{stats.Counter.Star}</span>
        </div>
      </div>
      <div class="stat-card island tonal">
        <span class="material-symbols-rounded icon">visibility</span>
        <div class="stat-content">
          <span class="label">Watches</span>
          <span class="value">{stats.Counter.Watch}</span>
        </div>
      </div>
      <div class="stat-card island tonal">
        <span class="material-symbols-rounded icon">merge_type</span>
        <div class="stat-content">
          <span class="label">Attachments</span>
          <span class="value">{stats.Counter.Attachment}</span>
        </div>
      </div>
      <div class="stat-card island tonal">
        <span class="material-symbols-rounded icon">notifications</span>
        <div class="stat-content">
          <span class="label">Notices</span>
          <span class="value">{stats.Counter.Notice}</span>
        </div>
      </div>
    </div>
  {:else}
    <div class="error-zone island tonal">
       <span class="material-symbols-rounded">lock</span>
       <div class="error-content">
          <h3>Access Denied or API Error</h3>
          <p>We couldn't retrieve the system statistics. Please verify you're logged in as an administrator.</p>
          <div class="error-actions">
             <button class="btn btn-primary" onclick={fetchStats}>Retry Connection</button>
             <a href="/user/login" class="btn btn-outline">Go to Login</a>
          </div>
       </div>
    </div>
  {/if}
</div>

<style>
  .admin-dashboard { display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header h2 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .subtitle { color: var(--text-muted); margin: 0.25rem 0 0 0; }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 1.5rem;
  }

  .stat-card {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 1.5rem;
  }

  .icon {
    font-size: 2.5rem;
    color: var(--primary-color);
    opacity: 0.8;
  }

  .stat-content {
    display: flex;
    flex-direction: column;
  }

  .label {
    font-size: 0.75rem;
    font-weight: 800;
    text-transform: uppercase;
    color: var(--text-muted);
    letter-spacing: 0.05em;
  }

  .value {
    font-size: 1.75rem;
    font-weight: 800;
    color: var(--text-color);
  }

  .error-zone {
    display: flex;
    align-items: center;
    gap: 2rem;
    padding: 3rem;
    background: rgba(var(--primary-rgb), 0.05);
    border: 1px dashed rgba(var(--primary-rgb), 0.2);
  }

  .error-zone span { font-size: 4rem; color: var(--primary-color); opacity: 0.5; }
  .error-content h3 { margin: 0 0 0.5rem 0; font-size: 1.5rem; }
  .error-content p { margin: 0 0 1.5rem 0; color: var(--text-muted); max-width: 400px; }
  .error-actions { display: flex; gap: 1rem; }

  .skeleton {
    height: 100px;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
