<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Avatar from '$lib/components/Avatar.svelte';
  
  let { q = '' } = $props();
  
  let orgs = $state<any[]>([]);
  let loading = $state(true);

  async function fetchOrgs() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/explore/organizations?q=${encodeURIComponent(q)}`);
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
        <div class="skeleton-org"></div>
      {/each}
    </div>
  {:else if orgs.length > 0}
    <div class="orgs-grid">
      {#each orgs as org}
        <a href="/org/{org.username}" class="org-card tonal hover-lift">
          <Avatar src={org.avatar_url} username={org.username} size={64} class="org-avatar" />
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
    <div class="empty-state tonal">
      <span class="material-symbols-rounded">corporate_fare</span>
      <h3>No organizations found</h3>
      <p>Try searching for a different organization name.</p>
    </div>
  {/if}
</div>

<style>
  .orgs-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem; grid-auto-rows: min-content; }
  .org-card { 
    display: flex; 
    gap: 1.25rem; 
    padding: 1.5rem; 
    text-decoration: none; 
    color: inherit; 
    align-items: flex-start; 
    border: 1px solid var(--md-sys-color-outline-variant); 
    border-radius: 2rem;
    background: var(--md-sys-color-surface-container-low);
    transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1); 
    position: relative;
  }
  .org-card:hover { 
    border-color: var(--primary-color); 
    background: var(--md-sys-color-surface-container);
    transform: translateY(-4px);
    box-shadow: 0 12px 24px -10px rgba(0, 0, 0, 0.2);
  }

  :global(.org-avatar) { border-radius: 12px !important; box-shadow: var(--shadow-sm); transition: all 0.3s; border: 1px solid var(--md-sys-color-outline-variant); }
  .org-card:hover :global(.org-avatar) { border-color: var(--primary-color); transform: scale(1.05); }
  
  .org-info { display: flex; flex-direction: column; gap: 0.25rem; flex: 1; }
  .org-info header { display: flex; flex-direction: column; }
  .org-name { font-size: 1.05rem; font-weight: 800; color: var(--md-sys-color-on-surface); }
  .org-username { font-size: 0.8rem; color: var(--text-muted); font-weight: 600; opacity: 0.7; }
  
  .org-desc { font-size: 0.85rem; color: var(--text-muted); margin: 0.4rem 0; line-height: 1.5; display: -webkit-box; line-clamp: 2; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; opacity: 0.8; }

  .org-meta { display: flex; flex-wrap: wrap; gap: 0.75rem; margin-top: 0.4rem; }
  .meta-item { display: flex; align-items: center; gap: 0.3rem; font-size: 0.7rem; color: var(--text-muted); font-weight: 700; opacity: 0.6; }
  .meta-item span { font-size: 0.9rem; color: var(--primary-color); }

  .empty-state { padding: 5rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state span { font-size: 4rem; color: var(--primary-color); opacity: 0.3; }
  .empty-state h3 { margin: 0; }
  .empty-state p { color: var(--text-muted); }

  .skeletons { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1rem; }
  .skeleton-org { height: 130px; border-radius: 1rem; animation: pulse 1.5s infinite; background: var(--md-sys-color-surface-container); }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.5; } 100% { opacity: 0.3; } }
</style>
