<script lang="ts">
  

import LayoutDashboard from "lucide-svelte/icons/layout-dashboard";
import Plus from "lucide-svelte/icons/plus";
import History from "lucide-svelte/icons/history";
import List from "lucide-svelte/icons/list";
import FolderOpen from "lucide-svelte/icons/folder-open";
import Inbox from "lucide-svelte/icons/inbox";
import Rocket from "lucide-svelte/icons/rocket";
import Folder from "lucide-svelte/icons/folder";
import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
import GitPullRequest from "lucide-svelte/icons/git-pull-request";
import AlertCircle from "lucide-svelte/icons/alert-circle";
import Building2 from "lucide-svelte/icons/building-2";
import ScrollText from "lucide-svelte/icons/scroll-text";


  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import { userPrefs } from '$lib/settings.svelte';
  import { onMount } from 'svelte';
  import RepoRow from './RepoRow.svelte';
  import GlobalSearch from './GlobalSearch.svelte';
  import ViewHeader from './ui/ViewHeader.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  import Button from './ui/Button.svelte';
  import DashboardInbox from './DashboardInbox.svelte';
  import DashboardSidebar from './DashboardSidebar.svelte';

  let data = $state<{
    repositories: any[],
    collaborative: any[],
    organizations: any[],
    gists: any[]
  }>({
    repositories: [],
    collaborative: [],
    organizations: [],
    gists: []
  });
  let loading = $state(true);
  let viewMode = $state<'list' | 'temporal' | 'folders' | 'inbox'>('temporal');
  let searchQuery = $state('');

  const oneWeekAgo = $derived(Date.now() - (7 * 24 * 60 * 60 * 1000));
  const allRepos = $derived.by(() => {
    const raw = [...(data?.repositories || []), ...(data?.collaborative || [])];
    const unique = new Map();
    raw.forEach(r => {
      if (r && !unique.has(r.id)) unique.set(r.id, r);
    });
    const result = Array.from(unique.values());
    console.log("DASHBOARD DEBUG - Unique Repos Count:", result.length, "from raw:", raw.length);
    return result;
  });
  
  const filteredRepos = $derived.by(() => {
    const repos = allRepos || [];
    if (!searchQuery) return repos;
    const q = searchQuery.toLowerCase();
    const result = repos.filter(r => 
      r?.name?.toLowerCase()?.includes(q) || 
      (r?.description && r.description.toLowerCase().includes(q))
    );
    console.log("DASHBOARD DEBUG - Filtered Repos Count:", result.length, "Query:", searchQuery);
    return result;
  });

  interface GroupedRepos {
    all?: any[];
    recent?: any[];
    older?: any[];
    folders?: Record<string, any[]>;
    ungrouped?: any[];
  }

  const groupedRepos = $derived.by((): GroupedRepos => {
    const repos = filteredRepos;
    // Robust sorting with null guards
    const all = [...repos].filter(r => r && r.name).sort((a, b) => a.name.localeCompare(b.name));
    
    if (viewMode === 'folders') {
      function folderDisplayName(repo: any, folderName: string): string | undefined {
        const name = repo.name;
        if (name && folderName && name.startsWith(folderName + '-')) {
          return name.substring(folderName.length + 1);
        }
        return undefined;
      }

      const folders: Record<string, any[]> = {};
      const ungrouped: any[] = [];
      repos.forEach(r => {
        if (r && r.project_folder) {
          if (!folders[r.project_folder]) folders[r.project_folder] = [];
          const displayName = folderDisplayName(r, r.project_folder);
          folders[r.project_folder].push(displayName ? { ...r, _displayName: displayName } : r);
        } else if (r) {
          if (userPrefs.current.autoFolder) {
            const dashIndex = r.name?.indexOf('-');
            if (dashIndex && dashIndex > 0) {
              const prefix = r.name.substring(0, dashIndex).trim();
              const shortName = r.name.substring(dashIndex + 1);
              if (!folders[prefix]) folders[prefix] = [];
              folders[prefix].push({ ...r, _displayName: shortName });
            } else {
              ungrouped.push(r);
            }
          } else {
            ungrouped.push(r);
          }
        }
      });
      return { folders, ungrouped: ungrouped.sort((a, b) => (a.name || '').localeCompare(b.name || '')) };
    }

    if (viewMode !== 'temporal') {
      return { all };
    }

    return {
      all,
      recent: repos.filter(r => r && r.updated_at && new Date(r.updated_at).getTime() > oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime()),
      older: repos.filter(r => r && r.updated_at && new Date(r.updated_at).getTime() <= oneWeekAgo).sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    };
  });

  let inboxItems = $state<any[]>([]);
  let inboxLoading = $state(false);

  async function loadInbox() {
    inboxLoading = true;
    try {
      const res = await apiFetch('/api/v2/user/dashboard/inbox');
      if (res.ok) {
        const json = await res.json();
        inboxItems = Array.isArray(json) ? json : [];
      }
    } catch (e) {
      console.error("Failed to fetch inbox", e);
      inboxItems = [];
    } finally {
      inboxLoading = false;
    }
  }

  async function loadDashboard() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/user/dashboard/data');
      if (res.ok) {
        const json = await res.json();
        console.log("DASHBOARD DEBUG - Raw JSON:", json);
        
        if (json) {
          // Robust mapping of fields from API to local state
          data.repositories = json.repositories || json.repos || json.data || (Array.isArray(json) ? json : []);
          data.collaborative = json.collaborative || [];
          data.organizations = json.organizations || json.orgs || [];
          
          console.log("DASHBOARD DEBUG - Mapped Repos:", data.repositories);
        }
      }
      // Silently prefetch inbox to show badge
      const inboxRes = await apiFetch('/api/v2/user/dashboard/inbox');
      if (inboxRes.ok) {
        const json = await inboxRes.json();
        inboxItems = Array.isArray(json) ? json : [];
      }

      // Fetch gists for sidebar
      const gistsRes = await apiFetch(`/api/v2/users/${authState.user?.username}/gists`);
      if (gistsRes.ok) {
        const json = await gistsRes.json();
        data.gists = Array.isArray(json) ? json : [];
      }
    } catch (e) {
      console.error("Failed to fetch dashboard data", e);
    } finally {
      loading = false;
    }
  }

  onMount(() => {
    const handleGlobalSearch = (e: any) => { searchQuery = e.detail; };
    window.addEventListener('ggit:search', handleGlobalSearch);
    loadDashboard();

    return () => window.removeEventListener('ggit:search', handleGlobalSearch);
  });

  $effect(() => {
    if (viewMode === 'inbox') loadInbox();
  });

  const hasSidebar = $derived(
    (data?.organizations?.length ?? 0) > 0 && userPrefs.current.showSidebarOrgs
  );

</script>

<div class="dashboard-layout" class:full-width={!hasSidebar}>
    <div class="dashboard-main-content">
        <PageHeader 
            title="Dashboard" 
            subtitle="Gestiona tus proyectos, repositorios y actividad reciente."
            icon={LayoutDashboard}
        >
            {#snippet actions()}
                <Button href="/repo/create" icon={Plus}>
                    <span class="hide-mobile">Nuevo Repositorio</span>
                </Button>
            {/snippet}
        </PageHeader>

        <ViewHeader 
            bind:activeTab={viewMode}
            tabs={[
                { 
                    id: 'temporal', 
                    label: 'Activity', 
                    icon: History,
                    hasBadge: (groupedRepos.recent?.length ?? 0) > 0 && (groupedRepos.recent?.some(r => (Date.now() - new Date(r.updated_at).getTime()) < (24 * 60 * 60 * 1000)))
                },
                { id: 'list', label: 'Alphabetical', icon: List },
                { id: 'folders', label: 'Folders', icon: FolderOpen },
                { id: 'inbox', label: 'Inbox', icon: Inbox, hasBadge: (inboxItems?.length ?? 0) > 0 }
            ]}
        />

        <div class="repository-stream">
            {#if loading}
                <div class="skeleton-list">
                    {#each Array(5) as _}<div class="skeleton-row"></div>{/each}
                </div>
            {:else if allRepos.length === 0}
                <div class="empty-dashboard">
                    <Rocket size={48} />
                    <h3>Start your first project</h3>
                    <p>No repositories found. Create one to begin your journey.</p>
                </div>
            {:else}
                <div class="flat-repo-list">
                {#if viewMode === 'temporal'}
                    {#each groupedRepos.recent || [] as repo}
                        <RepoRow {repo} />
                    {/each}
                    {#if (groupedRepos.older?.length ?? 0) > 0}
                        <div class="divider-text">Older Projects</div>
                        {#each groupedRepos.older || [] as repo}
                            <RepoRow {repo} isOlder={true} />
                        {/each}
                    {/if}
                {:else if viewMode === 'list'}
                    {#each groupedRepos.all || [] as repo}
                        <RepoRow {repo} />
                    {/each}
                {:else if viewMode === 'folders'}
                    {#each Object.entries(groupedRepos.folders || {}) as [name, repos]}
                        <div class="folder-group">
                            <div class="folder-header">
                                <Folder size={20} />
                                <h3>{name}</h3>
                                <span class="badge">{repos.length}</span>
                            </div>
                            <div class="folder-repos">
                                {#each repos as repo}
                                    <RepoRow {repo} displayName={repo._displayName} />
                                {/each}
                            </div>
                        </div>
                    {/each}
                    {#if (groupedRepos.ungrouped?.length ?? 0) > 0}
                        <div class="divider-text">Ungrouped Projects</div>
                        {#each groupedRepos.ungrouped || [] as repo}
                            <RepoRow {repo} />
                        {/each}
                    {/if}
                {:else if viewMode === 'inbox'}
                    <DashboardInbox {inboxItems} {inboxLoading} />
                {/if}
                </div>
            {/if}
        </div>
    </div>

    {#if hasSidebar}
        <DashboardSidebar gists={data.gists || []} organizations={data.organizations || []} />
    {/if}
</div>

<style>
  .dashboard-layout {
    display: flex;
    gap: 2rem;
    height: 100%;
    width: 100%;
    box-sizing: border-box;
    flex: 1 1 100%;
  }

  .dashboard-layout.full-width {
    gap: 0;
  }

  .dashboard-main-content { 
    display: flex; 
    flex-direction: column; 
    flex: 1 1 100%;
    width: 100%;
    min-width: 0;
  }

  .repository-stream {
    padding: 2rem 1.5rem 4rem;
    width: 100%;
    animation: fadeIn 0.4s ease-out;
  }

  @media (max-width: 600px) {
    .repository-stream { padding: 1rem 0.75rem 5rem; }
    .folder-repos { padding-left: 0.75rem; margin-left: 0.5rem; }
  }

  @media (max-width: 480px) {
    .repository-stream { padding: 0.5rem 0.25rem 5rem; }
    .divider-text { padding: 1.5rem 0.5rem 0.25rem; }
    .folder-group { margin-bottom: 1rem; }
  }
  
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .dashboard-sidebar {
    width: 320px;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding-top: 1.5rem;
    padding-right: 1.5rem;
  }

  .sidebar-card {
    border-radius: 1.5rem;
    padding: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .sidebar-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 1.25rem;
    color: var(--text-muted);
  }

  .sidebar-header h4 {
    margin: 0;
    font-size: 0.9rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    flex: 1;
  }

  .mini-add-link { color: var(--text-muted); transition: color 0.2s; }
  .mini-add-link:hover { color: var(--primary-color); }

  .gist-mini-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .gist-mini-item { text-decoration: none; color: var(--text-color); font-size: 0.85rem; font-weight: 600; padding: 0.5rem 0.75rem; border-radius: 0.75rem; background: rgba(255,255,255,0.02); transition: all 0.2s; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .gist-mini-item:hover { background: var(--primary-color); color: white; transform: translateX(4px); }



  :global(.action-summary a) {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 700;
  }

  :global(.branch-pill) {
    font-family: monospace;
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    padding: 1px 4px;
    border-radius: 4px;
    font-size: 0.75rem;
  }

  :global(.branch-pill.deleted) {
    background: rgba(var(--error-rgb), 0.1);
    color: var(--error);
    text-decoration: line-through;
  }

  .org-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
    gap: 0.75rem;
  }

  .org-chip img {
    width: 40px;
    height: 40px;
    border-radius: 10px;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: transform 0.2s;
  }

  .org-chip:hover img {
    transform: scale(1.1) translateY(-2px);
    border-color: var(--primary-color);
  }

  .flat-repo-list { 
    display: flex; 
    flex-direction: column; 
    gap: 1.25rem;
    background: transparent; 
    border: none; 
    border-radius: 0; 
    overflow: visible; 
  }

  .inbox-stream { display: flex; flex-direction: column; gap: 1rem; width: 100%; }
  .inbox-card { display: flex; gap: 1.25rem; padding: 1.25rem; border-radius: 1.5rem; text-decoration: none; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1); }
  .inbox-card:hover { transform: translateY(-2px); border-color: var(--primary-color); background: var(--md-sys-color-surface-variant); }
  .inbox-icon { width: 44px; height: 44px; border-radius: 12px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
  .inbox-icon.issue { background: rgba(var(--error-rgb), 0.1); color: var(--error); }
  .inbox-icon.pull { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); }
  .inbox-content { display: flex; flex-direction: column; gap: 0.15rem; min-width: 0; }
  .inbox-repo-line { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.6; }
  .inbox-title-line { display: flex; gap: 0.5rem; align-items: baseline; }
  .inbox-title-line .index { font-weight: 400; opacity: 0.4; font-family: monospace; }
  .inbox-title-line .title { font-weight: 800; color: var(--text-color); font-size: 1.05rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .inbox-meta-line { display: flex; align-items: center; gap: 0.5rem; font-size: 0.8rem; color: var(--text-muted); margin-top: 0.25rem; }
  .mini-avatar { width: 16px; height: 16px; border-radius: 4px; }

  :global(.repo-row) {
    background: var(--md-sys-color-surface-container-low) !important;
    border: 1px solid var(--md-sys-color-outline-variant) !important;
    border-radius: 1.5rem !important;
    box-shadow: 0 2px 8px rgba(0,0,0,0.02);
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1) !important;
  }

  :global(.repo-row:hover) {
    background: var(--md-sys-color-surface-variant) !important;
    border-color: var(--primary-color) !important;
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0,0,0,0.06);
  }

  .divider-text { 
    padding: 2rem 0.5rem 0.5rem; 
    font-size: 0.75rem; 
    font-weight: 900; 
    text-transform: uppercase; 
    color: var(--text-muted); 
    letter-spacing: 0.15em; 
    background: transparent;
    border: none;
  }

  .folder-group { margin-bottom: 2rem; }
  .folder-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1rem; color: var(--text-muted); }
  .folder-header h3 { font-size: 1rem; font-weight: 800; margin: 0; color: var(--text-color); }
  .folder-header .badge { background: var(--md-sys-color-secondary-container); color: var(--md-sys-color-on-secondary-container); padding: 2px 8px; border-radius: 10px; font-size: 0.7rem; font-weight: 800; }
  .folder-repos { display: flex; flex-direction: column; gap: 1rem; border-left: 2px solid var(--md-sys-color-outline-variant); padding-left: 1.5rem; margin-left: 0.75rem; }

  .empty-dashboard { padding: 2rem; text-align: center; color: var(--text-muted); }

  .skeleton-list { display: flex; flex-direction: column; gap: 1px; background: var(--md-sys-color-outline-variant); border-radius: 1rem; overflow: hidden; }
  .skeleton-row { height: 80px; background: var(--bg-color); animation: pulse 1.5s infinite; }
  
  @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.7; } }

  @media (max-width: 1100px) {
    .dashboard-layout { padding: 0; }
    .dashboard-sidebar { display: none; }
  }

  @media (max-width: 600px) {
    .hide-mobile { display: none; }
  }
</style>
