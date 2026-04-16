<script lang="ts">
  import { onMount } from 'svelte';
  
  let { q = '' } = $props();
  
  let orgs = $state<any[]>([]);
  let loading = $state(true);

  async function fetchOrgs() {
    loading = true;
    try {
      const res = await fetch(`/api/v2/explore/organizations?q=${encodeURIComponent(q)}`);
      if (res.ok) {
        const data = await res.json();
        orgs = data.organizations || [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    fetchOrgs();
  });

  onMount(fetchOrgs);
</script>

<div class="explore-orgs">
  {#if loading && orgs.length === 0}
    <div class="skeletons">
      {#each Array(6) as _}
        <div class="skeleton-org island"></div>
      {#/each}
    </div>
  {:else if orgs.length > 0}
    <div class="orgs-grid">
      {#each orgs as org}
        <a href="/org/{org.username}" class="org-card island tonal hover-lift">
          <img src={org.avatar_url} alt="" class="org-avatar" />
          <div class="org-info">
            <header>
               <strong class="org-name">{org.full_name || org.username}</strong>
               <span class="org-username">@{org.username}</span>
            </header>
            {#if org.description}
                <p class="org-desc">{org.description}</p>
            {/if}
            <div class="org-meta">
               {#if org.location}
                 <div class="meta-item">
                    <span class="material-symbols-rounded">location_on</span>
                    {org.location}
                 </div>
               {/if}
               {#if org.website}
                 <div class="meta-item">
                    <span class="material-symbols-rounded">link</span>
                    {org.website.replace('http://', '').replace('https://', '')}
                 </div>
               {/if}
            </div>
          </div>
        </a>
      {/each}
    </div>
  {:else}
    <div class="empty-state island tonal">
      <span class="material-symbols-rounded">corporate_fare</span>
      <h3>No organizations found</h3>
      <p>Try searching for a different organization name.</p>
    </div>
  {/if}
</div>

<style>
  .orgs-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem; }
  .org-card { display: flex; gap: 1.5rem; padding: 1.75rem; text-decoration: none; color: inherit; align-items: flex-start; border: 1px solid var(--border-color); height: 100%; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
  .org-card:hover { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .org-avatar { width: 72px; height: 72px; border-radius: 16px; box-shadow: var(--shadow-md); background: #111; }
  
  .org-info { display: flex; flex-direction: column; gap: 0.5rem; flex: 1; }
  .org-info header { display: flex; flex-direction: column; }
  .org-name { font-size: 1.2rem; font-weight: 800; color: #fff; }
  .org-username { font-size: 0.85rem; color: var(--text-muted); font-weight: 600; }
  
  .org-desc { font-size: 0.9rem; color: var(--text-muted); margin: 0.5rem 0; line-height: 1.5; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }

  .org-meta { display: flex; flex-wrap: wrap; gap: 1rem; margin-top: 0.5rem; }
  .meta-item { display: flex; align-items: center; gap: 0.4rem; font-size: 0.75rem; color: var(--text-muted); font-weight: 700; opacity: 0.8; }
  .meta-item span { font-size: 1rem; color: var(--primary-color); }

  .empty-state { padding: 5rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state span { font-size: 4rem; color: var(--primary-color); opacity: 0.3; }
  .empty-state h3 { margin: 0; }
  .empty-state p { color: var(--text-muted); }

  .skeletons { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem; }
  .skeleton-org { height: 140px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.5; } 100% { opacity: 0.3; } }
</style>
