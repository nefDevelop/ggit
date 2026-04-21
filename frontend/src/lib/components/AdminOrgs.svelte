<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let orgs = $state<any[]>([]);
  let loading = $state(true);
  let page = $state(1);

  async function fetchOrgs() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/orgs?page=${page}`);
      if (res.ok) {
        orgs = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  onMount(fetchOrgs);
</script>

<div class="admin-orgs island tonal">
  <header class="section-header">
    <div class="title">
      <span class="material-symbols-rounded">corporate_fare</span>
      <h2>Organization Oversight</h2>
    </div>
  </header>

  <div class="table-container">
    <table class="modern-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Organization</th>
          <th>Full Name</th>
          <th>Owner</th>
          <th>Teams</th>
          <th>Members</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          {#each Array(5) as _}
            <tr class="skeleton-row">
              <td colspan="7"><div class="skeleton-bar"></div></td>
            </tr>
          {/each}
        {:else}
          {#each orgs as org}
            <tr>
              <td>{org.id}</td>
              <td><b>{org.username}</b></td>
              <td>{org.full_name || ''}</td>
              <td>{org.owner?.username}</td>
              <td>{org.num_teams}</td>
              <td>{org.num_members}</td>
              <td class="actions">
                <a href="/org/{org.username}" class="icon-btn" title="View"><span class="material-symbols-rounded">visibility</span></a>
                <button class="icon-btn danger" title="Delete"><span class="material-symbols-rounded">delete</span></button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>

<style>
  .admin-orgs { padding: 2rem; display: flex; flex-direction: column; gap: 2rem; box-shadow: var(--shadow-xl); }
  
  .section-header { display: flex; justify-content: space-between; align-items: center; }
  .title { display: flex; align-items: center; gap: 0.75rem; }
  .title span { font-size: 1.75rem; color: var(--primary-color); }
  .title h2 { margin: 0; font-size: 1.5rem; font-weight: 800; }

  .table-container { overflow-x: auto; }
  .modern-table { width: 100%; border-collapse: collapse; text-align: left; }
  .modern-table th { padding: 1rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); color: var(--text-muted); font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.1em; font-weight: 800; }
  .modern-table td { padding: 1.25rem 1rem; border-bottom: 1px solid rgba(255,255,255,0.03); font-size: 0.95rem; }

  .actions { display: flex; gap: 0.5rem; }
  .icon-btn { background: transparent; border: none; color: var(--text-muted); cursor: pointer; padding: 0.4rem; border-radius: 0.5rem; transition: all 0.2s; text-decoration: none; display: flex; align-items: center; }
  .icon-btn:hover { background: rgba(255,255,255,0.05); color: #fff; }
  .icon-btn.danger:hover { background: rgba(239, 68, 68, 0.1); color: #ef4444; }

  .skeleton-bar { height: 1.5rem; background: rgba(255,255,255,0.03); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }
</style>
