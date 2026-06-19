<script lang="ts">
  import { page } from "$app/state";
  import Avatar from "./Avatar.svelte";
  import Plus from "lucide-svelte/icons/plus";

  let { navItems = [], adminSubItems = [], settingsSubItems = [], userOrgs = [], showSidebarOrgs = true }: {
    navItems: any[];
    adminSubItems: any[];
    settingsSubItems: any[];
    userOrgs: any[];
    showSidebarOrgs: boolean;
  } = $props();

  const isAdminPath = $derived(page.url.pathname.startsWith("/admin"));
  const isSettingsPath = $derived(page.url.pathname.startsWith("/user/settings"));
</script>

<div class="nav-links">
  {#each navItems as item}
    {@const Icon = item.icon}
    <a href={item.path} class="nav-item"
      class:active={page.url.pathname === item.path || (item.isParent && ((isAdminPath && !item.isSettings) || (isSettingsPath && item.isSettings)))}
    >
      <Icon size={20} />
      <span class="item-label">{item.name}</span>
    </a>

    {#if item.isParent && ((isAdminPath && !item.isSettings) || (isSettingsPath && item.isSettings))}
      <div class="sub-nav-links">
        {#each (item.isSettings ? settingsSubItems : adminSubItems) as sub, idx}
          {@const SubIcon = sub.icon}
          <a href={sub.path} class="sub-nav-item"
            class:active={page.url.pathname + page.url.search === sub.path || (page.url.pathname === sub.path.split('?')[0] && page.url.search === sub.path.substring(sub.path.indexOf('?'))) || (idx === 0 && page.url.pathname === sub.path.split('?')[0] && page.url.search === '')}
          >
            <SubIcon size={16} />
            <span class="item-label">{sub.name}</span>
          </a>
        {/each}
      </div>
    {/if}
  {/each}

  {#if showSidebarOrgs}
    {#if userOrgs.length > 0}
    <div class="nav-section-title">
      <span>Organizaciones</span>
      <a href="/org/create" class="mini-add" title="Nueva Organización" aria-label="Nueva Organización">
        <Plus size={18} />
      </a>
    </div>
    <div class="nav-org-list">
      {#each userOrgs as org}
        <a href="/org/{org.username}" class="nav-item org-item" class:active={page.url.pathname === `/org/${org.username}`}>
          <Avatar src={org.avatar_url} username={org.username} size={24} />
          <span class="item-label">{org.full_name || org.username}</span>
        </a>
      {/each}
    </div>
    {:else}
      <p class="nav-empty-hint">No eres parte de ninguna organización.</p>
    {/if}
  {/if}
</div>

<style>
  .nav-links {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex: 1;
  }
  .nav-item {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 0.85rem 1.25rem;
    color: var(--text-muted);
    text-decoration: none;
    border-radius: 0.5rem;
    transition: all 0.2s;
    font-weight: 600;
    font-size: 0.95rem;
    min-width: 0;
  }
  .nav-item:hover { color: var(--text-color); background: rgba(255,255,255,0.03); }
  .nav-item.active { color: var(--primary-color); background: rgba(var(--primary-rgb),0.08); }
  .item-label { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; flex: 1; }
  .sub-nav-links {
    display: flex; flex-direction: column; gap: 0.25rem;
    padding-left: 0.75rem; margin: 0.25rem 0;
    border-left: 1px solid var(--md-sys-color-outline-variant);
    margin-left: 1.25rem;
  }
  .sub-nav-item {
    display: flex; align-items: center; gap: 1rem;
    padding: 0.65rem 1rem; color: var(--text-muted);
    text-decoration: none; border-radius: 0.4rem;
    transition: all 0.2s; font-size: 0.85rem; font-weight: 600;
  }
  .sub-nav-item:hover { color: var(--text-color); background: rgba(255,255,255,0.03); }
  .sub-nav-item.active {
    color: var(--md-sys-color-on-secondary-container);
    background: var(--md-sys-color-secondary-container);
    box-shadow: var(--shadow-sm);
  }
  .nav-section-title {
    display: flex; align-items: center; justify-content: space-between;
    padding: 1.5rem 1.25rem 0.5rem;
    font-size: 0.75rem; font-weight: 800; text-transform: uppercase;
    color: var(--text-muted); letter-spacing: 0.05em;
  }
  .nav-org-list { display: flex; flex-direction: column; gap: 0.25rem; }
  .org-item { padding: 0.65rem 1.25rem; font-size: 0.9rem; }
  .mini-add { color: var(--text-muted); text-decoration: none; transition: color 0.2s; display: flex; }
  .mini-add:hover { color: var(--primary-color); }
  .nav-empty-hint {
    font-size: 0.75rem; color: var(--text-muted);
    padding: 0.5rem 1.25rem; opacity: 0.6; font-style: italic;
  }
  @media (max-width: 1024px) {
    .nav-links { flex-direction: row; width: auto; flex: 1; justify-content: space-around; gap: 0; }
    .nav-item { padding: 0.5rem; gap: 0; justify-content: center; flex: 1; }
    .item-label, .nav-section-title, .sub-nav-links, .nav-empty-hint { display: none; }
  }
</style>
