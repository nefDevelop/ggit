<script lang="ts">
  import { onMount } from 'svelte';
  
  let sources = $state<any[]>([]);
  let loading = $state(true);

  async function fetchAuths() {
    loading = true;
    try {
      const res = await fetch('/api/v2/admin/auths');
      if (res.ok) sources = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteAuth(id: number) {
    if (!confirm('Are you sure you want to delete this authentication source?')) return;
    try {
      const res = await fetch(`/api/v2/admin/auths/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        sources = sources.filter(s => s.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchAuths);
</script>

<div class="admin-auths">
  <header class="section-header">
    <div class="header-text">
       <h2>Authentication Sources</h2>
       <p class="subtitle">External authentication methods like LDAP, SMTP, or OAuth2.</p>
    </div>
    <button class="btn btn-primary">
      <span class="material-symbols-rounded">add_moderator</span>
      Add Source
    </button>
  </header>

  <div class="auth-grid">
    {#each sources as source}
      <div class="auth-card island tonal">
        <div class="card-top">
          <div class="type-icon">
            <span class="material-symbols-rounded">
              {source.type === 1 || source.type === 2 ? 'hub' : source.type === 3 ? 'mail' : 'security'}
            </span>
          </div>
          <div class="auth-info">
             <strong>{source.name}</strong>
             <span class="type-name">Type: {source.type_name || 'Legacy'}</span>
          </div>
          <div class="status-badge" class:active={source.is_actived}>
            {source.is_actived ? 'Active' : 'Disabled'}
          </div>
        </div>
        
        <div class="card-actions">
           <button class="btn btn-outline small">View Config</button>
           <button class="btn btn-danger-text small" onclick={() => deleteAuth(source.id)}>Delete</button>
        </div>
      </div>
    {/each}
  </div>

  {#if !loading && sources.length === 0}
    <div class="empty-state island">
      <span class="material-symbols-rounded">no_encryption</span>
      <h3>No authentication sources</h3>
      <p>Configure LDAP or SMTP to allow external user authentication.</p>
    </div>
  {/if}
</div>

<style>
  .admin-auths { display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header { display: flex; justify-content: space-between; align-items: flex-end; }
  .section-header h2 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .subtitle { color: var(--text-muted); margin: 0.25rem 0 0 0; }

  .auth-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.5rem; }
  
  .auth-card { padding: 1.5rem; display: flex; flex-direction: column; gap: 1.5rem; }
  .card-top { display: flex; align-items: flex-start; gap: 1rem; }
  
  .type-icon { width: 44px; height: 44px; border-radius: 12px; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); display: flex; align-items: center; justify-content: center; }
  .type-icon span { font-size: 1.5rem; }

  .auth-info { flex: 1; display: flex; flex-direction: column; }
  .auth-info strong { font-size: 1.1rem; }
  .type-name { font-size: 0.75rem; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; font-weight: 800; }

  .status-badge { font-size: 0.7rem; font-weight: 800; color: var(--text-muted); padding: 0.2rem 0.5rem; background: rgba(255,255,255,0.05); border-radius: 4px; text-transform: uppercase; }
  .status-badge.active { background: rgba(16, 185, 129, 0.1); color: #10b981; }

  .card-actions { display: flex; justify-content: flex-end; gap: 0.5rem; align-items: center; }
  .small { padding: 0.4rem 0.8rem; font-size: 0.8rem; }
  .btn-danger-text { color: #ef4444; }

  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); }
  .empty-state span { font-size: 4rem; opacity: 0.2; }
</style>
