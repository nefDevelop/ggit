<script lang="ts">
  import { userPrefs } from "$lib/settings.svelte";
  import ScrollText from "lucide-svelte/icons/scroll-text";
  import Building2 from "lucide-svelte/icons/building-2";
  import Plus from "lucide-svelte/icons/plus";

  let { gists = [], organizations = [] }: {
    gists: any[];
    organizations: any[];
  } = $props();
</script>

<aside class="dashboard-sidebar">
  {#if gists.length > 0}
    <div class="sidebar-card tonal">
      <header class="sidebar-header">
        <ScrollText size={18} />
        <h4>Recent Gists</h4>
        <a href="/gists/new" class="mini-add-link" title="New Gist" aria-label="New Gist"><Plus size={16} /></a>
      </header>
      <div class="gist-mini-list">
        {#each gists.slice(0, 5) as gist}
          <a href="/gists/{gist.uuid}" class="gist-mini-item">
            <span class="gist-mini-title">{gist.description || 'Untitled'}</span>
          </a>
        {/each}
      </div>
    </div>
  {/if}

  {#if organizations.length > 0 && userPrefs.current.showSidebarOrgs}
    <div class="sidebar-card tonal">
      <header class="sidebar-header">
        <Building2 size={18} />
        <h4>Organizations</h4>
        <a href="/org/create" class="mini-add-link" title="New Org" aria-label="New Organization"><Plus size={16} /></a>
      </header>
      <div class="org-grid">
        {#each organizations as org}
          <a href="/org/{org.username}" class="org-chip" title={org.full_name || org.username}>
            <img src={org.avatar_url} alt={org.username} loading="lazy" />
          </a>
        {/each}
      </div>
    </div>
  {/if}
</aside>

<style>
  .dashboard-sidebar {
    display: flex; flex-direction: column; gap: 1rem;
    width: 260px; flex-shrink: 0;
  }
  .sidebar-card {
    padding: 1.25rem; border-radius: 1.5rem;
    background: var(--md-sys-color-surface-container-low);
  }
  .sidebar-header {
    display: flex; align-items: center; gap: 0.5rem; margin-bottom: 1rem;
  }
  .sidebar-header h4 { margin: 0; font-size: 0.85rem; font-weight: 800; flex: 1; }
  .mini-add-link {
    color: var(--text-muted); text-decoration: none; transition: color 0.2s;
    display: flex; padding: 0.25rem;
  }
  .mini-add-link:hover { color: var(--primary-color); }
  .gist-mini-list { display: flex; flex-direction: column; gap: 0.25rem; }
  .gist-mini-item {
    padding: 0.5rem; border-radius: 0.5rem; color: var(--text-muted);
    text-decoration: none; transition: all 0.2s; font-size: 0.9rem;
  }
  .gist-mini-item:hover { background: rgba(255,255,255,0.03); color: var(--text-color); }
  .org-grid { display: flex; gap: 0.5rem; flex-wrap: wrap; }
  .org-chip img {
    width: 36px; height: 36px; border-radius: 50%;
    transition: transform 0.2s;
  }
  .org-chip:hover img { transform: scale(1.1); }
  @media (max-width: 1200px) {
    .dashboard-sidebar { width: 100%; flex-direction: row; flex-wrap: wrap; }
    .sidebar-card { flex: 1; min-width: 200px; }
  }
</style>
