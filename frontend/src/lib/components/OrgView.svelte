<script lang="ts">
  import { onMount } from 'svelte';
  import OrgTeams from './OrgTeams.svelte';
  import OrgMembers from './OrgMembers.svelte';
  import { apiFetch } from '$lib/api';
  import Avatar from '$lib/components/Avatar.svelte';

  let { orgName } = $props();
  
  let org = $state<any>(null);
  let loading = $state(true);
  let activeTab = $state('repos');
  let repos = $state<any[]>([]);

  async function fetchOrg() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}`);
      if (res.ok) {
        org = await res.json();
        fetchRepos();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function fetchRepos() {
    try {
      const res = await apiFetch(`/api/v2/user/settings/repos?uid=${org.id}`); // Reusing user repos logic
      if (res.ok) repos = await res.json();
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchOrg);
</script>

{#if loading}
  <div class="skeleton-org"></div>
{:else if org}
  <div class="org-layout">
    <header class="org-header island island-interactive tonal">
      <Avatar src={org.avatar_url} username={org.username} size={100} class="org-avatar" />
      <div class="org-info">
        <h1>{org.full_name || org.username}</h1>
        <p class="org-desc">{org.description || 'No description provided.'}</p>
        <div class="org-meta">
          {#if org.location}
            <span class="meta-item">
              <span class="material-symbols-rounded">location_on</span>
              {org.location}
            </span>
          {/if}
          {#if org.website}
            <a href={org.website} class="meta-item link" target="_blank">
              <span class="material-symbols-rounded">language</span>
              {org.website}
            </a>
          {/if}
        </div>
      </div>
      <div class="org-actions">
        <a href="/org/{org.username}/settings" class="btn btn-outline">
          <span class="material-symbols-rounded">settings</span>
          Settings
        </a>
      </div>
    </header>

    <nav class="org-nav">
      <button class:active={activeTab === 'repos'} onclick={() => activeTab = 'repos'}>
        <span class="material-symbols-rounded">book</span>
        Repositories
        <span class="count">{repos.length}</span>
      </button>
      <button class:active={activeTab === 'teams'} onclick={() => activeTab = 'teams'}>
        <span class="material-symbols-rounded">groups</span>
        Teams
        <span class="count">{org.num_teams}</span>
      </button>
      <button class:active={activeTab === 'members'} onclick={() => activeTab = 'members'}>
        <span class="material-symbols-rounded">person</span>
        People
        <span class="count">{org.num_members}</span>
      </button>
    </nav>

    <main class="org-content">
      {#if activeTab === 'repos'}
        <div class="repos-grid">
          {#each repos as repo}
            <!-- Assuming a generic repo card style -->
            <div class="repo-card island island-interactive tonal hover-lift">
               <strong>{repo.name}</strong>
               <p>{repo.description || ''}</p>
            </div>
          {/each}
        </div>
      {:else if activeTab === 'teams'}
        <OrgTeams {orgName} />
      {:else if activeTab === 'members'}
        <OrgMembers {orgName} />
      {/if}
    </main>
  </div>
{/if}

<style>
  .org-layout { display: flex; flex-direction: column; gap: 2rem; }
  
  .org-header { display: flex; align-items: center; gap: 2rem; padding: 2.5rem; }
  :global(.org-avatar) { border-radius: 1.5rem !important; background: #222; }
  
  .org-info { flex: 1; display: flex; flex-direction: column; gap: 0.5rem; }
  .org-info h1 { margin: 0; font-size: 2rem; font-weight: 800; }
  .org-desc { font-size: 1.1rem; color: var(--text-muted); margin: 0; }
  
  .org-meta { display: flex; gap: 1.5rem; margin-top: 0.5rem; }
  .meta-item { display: flex; align-items: center; gap: 0.5rem; font-size: 0.9rem; color: var(--text-muted); }
  .meta-item.link { color: var(--primary-color); text-decoration: none; }
  .meta-item span { font-size: 1.2rem; }

  .org-nav { display: flex; gap: 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); padding: 0 1rem; }
  .org-nav button { background: transparent; border: none; padding: 1rem 0.5rem; color: var(--text-muted); font-weight: 800; display: flex; align-items: center; gap: 0.75rem; cursor: pointer; position: relative; transition: all 0.2s; }
  .org-nav button:hover { color: var(--text-color); }
  .org-nav button.active { color: var(--primary-color); }
  .org-nav button.active::after { content: ''; position: absolute; bottom: 0; left: 0; right: 0; height: 3px; background: var(--primary-color); border-radius: 3px 3px 0 0; }
  
  .count { background: rgba(255,255,255,0.05); padding: 0.1rem 0.5rem; border-radius: 1rem; font-size: 0.75rem; }

  .repos-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1.5rem; }
  .repo-card { padding: 1.5rem; display: flex; flex-direction: column; gap: 0.5rem; border: 1px solid var(--md-sys-color-outline-variant); }

  .skeleton-org { height: 600px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
