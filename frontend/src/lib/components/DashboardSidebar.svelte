<script lang="ts">
  import { userPrefs } from "$lib/settings.svelte";
  import Avatar from "./Avatar.svelte";
  import ScrollText from "lucide-svelte/icons/scroll-text";
  import Building2 from "lucide-svelte/icons/building-2";
  import FileText from "lucide-svelte/icons/file-text";
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
          <a href="/gists/{gist.uuid}" class="gist-mini-item" title={gist.description || 'Untitled'}>
            <FileText size={16} class="gist-mini-icon" />
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
      <div class="org-list">
        {#each organizations as org}
          <a href="/org/{org.username}" class="org-row" title={org.full_name || org.username}>
            <Avatar src={org.avatar_url} username={org.username} size={32} />
            <span class="org-text">
              <span class="org-name">{org.full_name || org.username}</span>
              {#if org.full_name && org.full_name !== org.username}
                <span class="org-username">@{org.username}</span>
              {/if}
            </span>
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
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  .sidebar-header {
    display: flex; align-items: center; gap: 0.5rem; margin-bottom: 1rem;
    color: var(--md-sys-color-on-surface-variant);
  }
  .sidebar-header h4 {
    margin: 0; font-size: 0.75rem; font-weight: 800; flex: 1;
    text-transform: uppercase; letter-spacing: 0.06em;
  }
  .mini-add-link {
    color: var(--text-muted); text-decoration: none; transition: color 0.2s;
    display: flex; padding: 0.25rem;
  }
  .mini-add-link:hover { color: var(--primary-color); }
  .gist-mini-list { display: flex; flex-direction: column; gap: 0.25rem; }
  .gist-mini-item {
    display: flex; align-items: center; gap: 0.6rem;
    padding: 0.6rem 0.75rem; border-radius: 0.9rem; color: var(--text-muted);
    text-decoration: none; transition: all 0.2s; font-size: 0.9rem;
  }
  .gist-mini-item:hover { background: rgba(var(--primary-rgb), 0.08); color: var(--text-color); }
  .gist-mini-icon { flex-shrink: 0; }
  .gist-mini-item:hover .gist-mini-icon { color: var(--primary-color); }
  .gist-mini-title {
    flex: 1; min-width: 0; overflow: hidden;
    text-overflow: ellipsis; white-space: nowrap;
  }
  .org-list { display: flex; flex-direction: column; gap: 0.25rem; }
  .org-row {
    display: flex; align-items: center; gap: 0.7rem;
    padding: 0.5rem 0.6rem; border-radius: 0.9rem;
    text-decoration: none; transition: background 0.2s;
  }
  .org-row:hover { background: rgba(var(--primary-rgb), 0.08); }
  .org-text { display: flex; flex-direction: column; min-width: 0; }
  .org-name {
    font-size: 0.9rem; font-weight: 700; color: var(--text-color);
    overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  }
  .org-username { font-size: 0.75rem; color: var(--text-muted); }
  @media (max-width: 1200px) {
    .dashboard-sidebar { width: 100%; flex-direction: row; flex-wrap: wrap; }
    .sidebar-card { flex: 1; min-width: 200px; }
  }
  /* The rule that hides the sidebar below 1100px used to live in
     DashboardView.svelte. Scoped styles put a per-component class on the
     elements, so a selector in the parent never matched an element rendered by
     the child and the media query was dead code: the sidebar was never hidden.
     It has to be declared here, where the element is. */
  @media (max-width: 1100px) {
    .dashboard-sidebar { display: none; }
  }
</style>
