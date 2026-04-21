<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let orgs = $state<any[]>([]);
  let loading = $state(true);

  async function fetchOrgs() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/orgs');
      if (res.ok) {
        orgs = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function leaveOrg(id: number, name: string) {
    if (!confirm(`Are you sure you want to leave the organization ${name}?`)) return;
    try {
      const res = await apiFetch(`/api/v2/user/settings/orgs/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        await fetchOrgs();
      } else {
        const data = await res.json();
        alert(data.message || 'Failed to leave organization');
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchOrgs);
</script>

<div class="settings-card island island-interactive">
  <header class="card-header">
    <h3 class="card-title">Organizations</h3>
    <p class="card-desc">Organizations you are a member of.</p>
  </header>

  <div class="card-body">
    {#if loading}
      <div class="skeletons">
        {#each Array(3) as _}<div class="skeleton-org"></div>{/each}
      </div>
    {:else if orgs.length === 0}
      <div class="empty-state">
        <span class="material-symbols-rounded">corporate_fare</span>
        <p>You are not a member of any organization yet.</p>
      </div>
    {:else}
      <div class="org-list">
        {#each orgs as org}
          <div class="org-item island island-interactive tonal">
             <div class="org-main">
                <img src={org.avatar_url} alt="" class="org-avatar" />
                <div class="org-info">
                   <a href="/org/{org.username}" class="org-name">{org.full_name || org.username}</a>
                   <span class="org-user">@{org.username}</span>
                </div>
             </div>
             <button class="btn btn-ghost error btn-sm" onclick={() => leaveOrg(org.id, org.username)}>
                <span class="material-symbols-rounded">logout</span> Leave
             </button>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  .card-desc { margin: 0.5rem 0 0; font-size: 0.9rem; color: var(--text-muted); }

  .card-body { padding: 1.5rem; }

  .org-list { display: flex; flex-direction: column; gap: 1rem; }
  .org-item { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.5rem; }
  
  .org-main { display: flex; align-items: center; gap: 1rem; }
  .org-avatar { width: 40px; height: 40px; border-radius: 8px; background: #222; }
  
  .org-info { display: flex; flex-direction: column; }
  .org-name { font-weight: 700; color: #fff; text-decoration: none; }
  .org-name:hover { color: var(--primary-color); }
  .org-user { font-size: 0.8rem; color: var(--text-muted); }

  .skeleton-org { height: 72px; background: var(--md3-surface-3); border-radius: 1rem; margin-bottom: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); }
  .empty-state span { font-size: 3rem; opacity: 0.2; }

  .btn-ghost.error { color: #ef4444; }
  .btn-ghost.error:hover { background: rgba(239, 68, 68, 0.1); }
</style>
