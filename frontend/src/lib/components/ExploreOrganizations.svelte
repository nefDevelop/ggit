<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Avatar from '$lib/components/Avatar.svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  
  

import MapPin from "lucide-svelte/icons/map-pin";
import Link from "lucide-svelte/icons/link";
import Building2 from "lucide-svelte/icons/building-2";


  
  let { q = '' } = $props();
  
  let orgs = $state<any[]>([]);
  let loading = $state(true);

  async function fetchOrgs() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/explore/organizations?q=${encodeURIComponent(q)}`);
      if (res.ok) {
        const data = await res.json();
        // Validamos que tengan username antes de renderizar para evitar /org/undefined
        orgs = (data.organizations || []).filter((o: any) => o.username);
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
  <div class="container">
    {#if loading && orgs.length === 0}
      <div class="skeletons">
        {#each Array(6) as _}
          <Card variant="low" class="skeleton-org" />
        {/each}
      </div>
    {:else if orgs.length > 0}
      <div class="orgs-grid">
        {#each orgs as org}
          <a href="/org/{org.username}" class="card-link">
            <Card hover variant="low" padding="md" class="org-card">
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
                     <Badge variant="ghost" icon={MapPin} size="sm">
                       {org.location}
                     </Badge>
                   {/if}
                   {#if org.website}
                      <Badge variant="ghost" icon={Link} size="sm">
                        {org.website.replace(/^https?:\/\//, '')}
                      </Badge>
                   {/if}
                </div>
              </div>
            </Card>
          </a>
        {/each}
      </div>
    {:else}
      <div class="empty-state tonal">
        <Building2 size={64} />
        <h3>No organizations found</h3>
        <p>Try searching for a different organization name.</p>
      </div>
    {/if}
  </div>
</div>

<style>
  .explore-orgs { padding: 2rem 0; }
  @media (max-width: 600px) {
    .explore-orgs { padding: 1rem; }
  }
  .container { max-width: 1200px; margin: 0 auto; width: 100%; }
  .orgs-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem; }
  .card-link { text-decoration: none; color: inherit; display: block; }
  
  :global(.org-card) { display: flex !important; gap: 1.25rem !important; align-items: flex-start !important; }

  :global(.org-avatar) { border-radius: 1rem !important; border: 1px solid var(--md-sys-color-outline-variant); }
  
  .org-info { display: flex; flex-direction: column; gap: 0.25rem; flex: 1; min-width: 0; }
  .org-name { font-size: 1.1rem; font-weight: 800; color: var(--md-sys-color-on-surface); }
  .org-username { font-size: 0.8rem; color: var(--md-sys-color-on-surface-variant); opacity: 0.7; }
  
  .org-desc { 
    font-size: 0.85rem; 
    color: var(--md-sys-color-on-surface-variant); 
    margin: 0.5rem 0; 
    line-height: 1.5; 
    overflow: hidden;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .org-meta { display: flex; flex-wrap: wrap; gap: 0.5rem; margin-top: 0.5rem; }

  .empty-state { padding: 5rem; text-align: center; background: var(--md-sys-color-surface-container-low); border-radius: 1.5rem; }
  .empty-state :global(svg) { font-size: 4rem; color: var(--md-sys-color-primary); opacity: 0.3; }

  .skeletons { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1rem; }
  :global(.skeleton-org) { height: 140px; animation: pulse 1.5s infinite; background: var(--md-sys-color-surface-container-low); }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
</style>
