<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { onMount, setContext } from 'svelte';
  import type { RepoContext } from '$lib/types';
  import GlobalSearch from '$lib/components/GlobalSearch.svelte';
  import { userPrefs } from '$lib/settings.svelte';
  

import LayoutGrid from "lucide-svelte/icons/layout-grid";
import Lock from "lucide-svelte/icons/lock";
import Scale from "lucide-svelte/icons/scale";
import Download from "lucide-svelte/icons/download";
import ChevronDown from "lucide-svelte/icons/chevron-down";
import Copy from "lucide-svelte/icons/copy";
import Database from "lucide-svelte/icons/database";
import GitBranch from "lucide-svelte/icons/git-branch";
import FileArchive from "lucide-svelte/icons/file-archive";
import Archive from "lucide-svelte/icons/archive";
import Code from "lucide-svelte/icons/code";
import Terminal from "lucide-svelte/icons/terminal";
import Settings from "lucide-svelte/icons/settings";
import History from "lucide-svelte/icons/history";
import Tag from "lucide-svelte/icons/tag";
import Rocket from "lucide-svelte/icons/rocket";
import AlertCircle from "lucide-svelte/icons/alert-circle";
import GitMerge from "lucide-svelte/icons/git-merge";
import Book from "lucide-svelte/icons/book";
import Gavel from "lucide-svelte/icons/gavel";
import Calendar from "lucide-svelte/icons/calendar";


  
  import { fade } from 'svelte/transition';
  import Avatar from '$lib/components/Avatar.svelte';
  
  let { children } = $props();
  
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);
  let currentPath = $derived(page.url.pathname);

  let repoData = $state<any>(null);
  let fullRepoRes = $state<any>(null);
  let rootEntries = $state<any[]>([]);
  let gitStats = $state<any>({ commits: 0, branches: 0, tags: 0, releases: 0 });
  let languages = $state<any>(null);
  let loading = $state(true);
  let error = $state('');

  let lastFetched = { user: '', repo: '' };

  async function fetchRepoInfo() {
    if (lastFetched.user === username && lastFetched.repo === reponame) return;
    lastFetched = { user: username || '', repo: reponame || '' };

    if (!repoData) loading = true;
    error = '';
    try {
      const repoRes = await apiFetch(`/api/v2/${username}/${reponame}/view/`);
      if (repoRes.ok) {
        const res = await repoRes.json();
        repoData = res.repository;
        fullRepoRes = res; // Store the full response
        rootEntries = res.entries || [];
        if (res.stats) {
            gitStats = res.stats;
        }
        if (res.languages) {
            languages = res.languages;
        }
      } else {
        error = 'Repository not found or access denied.';
      }
    } catch (e) {
      error = 'Failed to load repository information.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    fetchRepoInfo();
  });

  setContext<RepoContext>('repo', {
    get data() { return repoData; },
    get full() { return fullRepoRes; }, // New full response
    get entries() { return rootEntries; },
    get stats() { return gitStats; },
    get languages() { return languages; },
    get loading() { return loading; }
  });

  let activeTab = $derived(
    currentPath.endsWith('/commits') || currentPath.includes('/commit/') ? 'commits' :
    currentPath.endsWith('/issues') || currentPath.includes('/issues/') ? 'issues' :
    currentPath.endsWith('/pulls') || currentPath.includes('/pulls/') ? 'pulls' :
    currentPath.endsWith('/wiki') || currentPath.includes('/wiki/') ? 'wiki' :
    currentPath.endsWith('/adr') || currentPath.includes('/adr/') ? 'adr' :
    currentPath.endsWith('/branches') ? 'branches' :
    currentPath.endsWith('/releases') ? 'releases' :
    currentPath.endsWith('/tags') ? 'tags' :
    currentPath.endsWith('/labels') ? 'labels' :
    currentPath.endsWith('/milestones') ? 'milestones' : 
    currentPath.endsWith('/settings') ? 'settings' : 'code'
  );

  let showCloneDropdown = $state(false);
  let cloneProtocol = $state<'https' | 'ssh'>('https');

  onMount(() => {
    const saved = localStorage.getItem('clone-protocol');
    if (saved === 'ssh' || saved === 'https') cloneProtocol = saved;
  });

  $effect(() => {
    localStorage.setItem('clone-protocol', cloneProtocol);
  });

  function copyCloneUrl(url: string) {
    navigator.clipboard.writeText(url);
    showCloneDropdown = false;
  }
</script>

<div class="workspace-shell">
  <div class="master-container">
    {#if loading && !repoData}
      <div class="skeleton-master">
        <div class="repo-header-region">
          <div class="tonal-search-bar">
            <div class="skeleton-pulse search-skeleton"></div>
          </div>

          <div class="repo-identity-bar">
            <div class="skeleton-pulse identity-skeleton"></div>
            <div class="skeleton-pulse actions-skeleton"></div>
          </div>

          <div class="rich-tabs">
            <div class="skeleton-pulse tabs-skeleton"></div>
          </div>
        </div>

        <div class="master-body-skeleton">
          <div class="skeleton-pulse body-skeleton"></div>
        </div>
      </div>
    {:else if error}
      <div class="tonal error-tonal">
        <AlertCircle size={48} />
        <p>{error}</p>
        <a href="/dashboard" class="btn btn-tonal">Back to Dashboard</a>
      </div>
    {:else if repoData}
      <div class="master-tonal tonal">
      <div class="repo-header-region">
        <div class="tonal-search-bar">
          <GlobalSearch />
        </div>

        <header class="repo-identity-bar">
          <div class="repo-header-main">
            <Avatar 
              src={repoData.avatar_url} 
              username={repoData.name} 
              size={36} 
              type="repo" 
              class="header-repo-avatar"
            />
            <div class="breadcrumb-nav">
              <a href="/dashboard" class="bc-home"><LayoutGrid size={20} /></a>
              <span class="bc-sep">/</span>
              <a href={`/${username}`} class="bc-item">{username}</a>
              <span class="bc-sep">/</span>
              <h1 class="bc-item active">{reponame}</h1>
              {#if repoData.private}<Lock size={16} class="lock-icon" />{/if}
            </div>
          </div>

          <div class="repo-meta-actions">
            {#if repoData.license}
              <div class="license-tag">
                <Scale size={16} />
                {repoData.license}
              </div>
            {/if}
            
            <div class="clone-zone">
              <div class="split-btn-group">
                <button class="btn btn-sm btn-primary main-clone-btn" onclick={() => showCloneDropdown = !showCloneDropdown}>
                  <Download size={18} />
                  <span class="clone-text">Clone</span>
                </button>
                <button class="btn btn-sm btn-primary arrow-btn" onclick={() => showCloneDropdown = !showCloneDropdown}>
                  <ChevronDown size={20} />
                </button>
              </div>

              {#if showCloneDropdown}
                <div class="clone-dropdown tonal shadow-xl">
                  <div class="dropdown-section">
                    <div class="protocol-tabs">
                      <button class="proto-btn" class:active={cloneProtocol === 'https'} onclick={() => cloneProtocol = 'https'}>HTTPS</button>
                      <button class="proto-btn" class:active={cloneProtocol === 'ssh'} onclick={() => cloneProtocol = 'ssh'}>SSH</button>
                    </div>
                    <div class="clone-copy-row">
                      <code class="url-display">{cloneProtocol === 'https' ? repoData.clone_url : repoData.ssh_url}</code>
                      <button class="copy-action" onclick={() => copyCloneUrl(cloneProtocol === 'https' ? repoData.clone_url : repoData.ssh_url)}>
                        <Copy size={18} />
                      </button>
                    </div>
                    <p class="protocol-hint">
                      {cloneProtocol === 'https' ? 'Use Git or checkout with SVN using the web URL.' : 'Use a password-protected SSH key.'}
                    </p>
                  </div>

                  <div class="dropdown-divider"></div>

                  <div class="dropdown-section">
                    <span class="section-title">Features</span>
                    <div class="feature-list">
                      <div class="feature-item" title="Git Large File Storage">
                        <Database size={18} />
                        <span>Git LFS enabled</span>
                      </div>
                      {#if repoData.has_submodules}
                      <div class="feature-item">
                        <GitBranch size={18} />
                        <span>Submodules included</span>
                      </div>
                      {/if}
                    </div>
                  </div>

                  <div class="dropdown-divider"></div>

                  <div class="dropdown-section">
                    <span class="section-title">Download Source</span>
                    <div class="archive-grid-compact">
                      <a href={`/api/v2/${username}/${reponame}/archive/${repoData.default_branch}.zip`} class="archive-tile">
                        <FileArchive size={20} />
                        <span>ZIP</span>
                      </a>
                      <a href={`/api/v2/${username}/${reponame}/archive/${repoData.default_branch}.tar.gz`} class="archive-tile">
                        <Archive size={20} />
                        <span>TAR.GZ</span>
                      </a>
                    </div>
                  </div>

                  <div class="dropdown-divider"></div>

                  <div class="dropdown-section">
                    <span class="section-title">Open in IDE</span>
                    <div class="archive-grid-compact">
                      <a href={`vscode://vscode.git/clone?url=${encodeURIComponent(cloneProtocol === 'https' ? repoData.clone_url : repoData.ssh_url)}`} class="archive-tile">
                        <Code size={20} />
                        <span>VS Code</span>
                      </a>
                      <a href={`vscodium://vscode.git/clone?url=${encodeURIComponent(cloneProtocol === 'https' ? repoData.clone_url : repoData.ssh_url)}`} class="archive-tile">
                        <Terminal size={20} />
                        <span>VSCodium</span>
                      </a>
                    </div>
                  </div>
                </div>
              {/if}
            </div>
            <a href={`/${username}/${reponame}/settings`} class="settings-link" class:active={activeTab === 'settings'} title="Settings">
              <Settings size={20} />
            </a>
          </div>
        </header>

        <nav class="rich-tabs">
          <a href={`/${username}/${reponame}`} class="tab" class:active={activeTab === 'code'}>
            <Code size={18} />
            {#if activeTab === 'code'}<span class="tab-label">Code</span>{/if}
          </a>

          {#if userPrefs.current.repoTabs.commits}
            <a href={`/${username}/${reponame}/commits`} class="tab" class:active={activeTab === 'commits'}>
              <History size={18} />
              <span class="tab-count">{gitStats.commits}</span>
              {#if activeTab === 'commits'}<span class="tab-label">Commits</span>{/if}
            </a>
          {/if}

          {#if userPrefs.current.repoTabs.branch}
            <a href={`/${username}/${reponame}/branches`} class="tab" class:active={activeTab === 'branches'}>
              <GitBranch size={18} />
              <span class="tab-count">{gitStats.branches}</span>
              {#if activeTab === 'branches'}<span class="tab-label">Branches</span>{/if}
            </a>
          {/if}

          {#if userPrefs.current.repoTabs.tags}
            <a href={`/${username}/${reponame}/tags`} class="tab" class:active={activeTab === 'tags'}>
              <Tag size={18} />
              <span class="tab-count">{gitStats.tags}</span>
              {#if activeTab === 'tags'}<span class="tab-label">Tags</span>{/if}
            </a>
          {/if}

          {#if userPrefs.current.repoTabs.releases}
            <a href={`/${username}/${reponame}/releases`} class="tab" class:active={activeTab === 'releases'}>
              <Rocket size={18} />
              <span class="tab-count">{gitStats.releases || 0}</span>
              {#if activeTab === 'releases'}<span class="tab-label">Releases</span>{/if}
            </a>
          {/if}

          {#if repoData.enable_issues && userPrefs.current.repoTabs.issue}
            <a href={repoData.enable_external_tracker ? repoData.external_tracker_url : `/${username}/${reponame}/issues`} 
               target={repoData.enable_external_tracker ? "_blank" : ""}
               rel={repoData.enable_external_tracker ? "noopener noreferrer" : ""}
               class="tab" class:active={activeTab === 'issues'}>
              <AlertCircle size={18} />
              <span class="tab-count">{repoData.open_issues_count || 0}</span>
              {#if activeTab === 'issues'}<span class="tab-label">Issues</span>{/if}
            </a>
          {/if}

          {#if repoData.enable_pulls && userPrefs.current.repoTabs.pr}
            <a href={`/${username}/${reponame}/pulls`} class="tab" class:active={activeTab === 'pulls'}>
              <GitMerge size={18} />
              <span class="tab-count">{repoData.open_pulls_count || 0}</span>
              {#if activeTab === 'pulls'}<span class="tab-label">PRs</span>{/if}
            </a>
          {/if}

          {#if repoData.enable_wiki && userPrefs.current.repoTabs.wiki}
            <a href={repoData.enable_external_wiki ? repoData.external_wiki_url : `/${username}/${reponame}/wiki`} 
               target={repoData.enable_external_wiki ? "_blank" : ""}
               rel={repoData.enable_external_wiki ? "noopener noreferrer" : ""}
               class="tab" class:active={activeTab === 'wiki'}>
              <Book size={18} />
              {#if activeTab === 'wiki'}<span class="tab-label">Wiki</span>{/if}
            </a>
          {/if}

          {#if userPrefs.current.repoTabs.adr}
            <a href={`/${username}/${reponame}/adr`} class="tab" class:active={activeTab === 'adr'}>
              <Gavel size={18} />
              {#if activeTab === 'adr'}<span class="tab-label">ADR</span>{/if}
            </a>
          {/if}

          {#if userPrefs.current.repoTabs.labels}
            <a href={`/${username}/${reponame}/labels`} class="tab" class:active={activeTab === 'labels'}>
              <Tag size={18} />
              {#if activeTab === 'labels'}<span class="tab-label">Labels</span>{/if}
            </a>
          {/if}

          {#if userPrefs.current.repoTabs.milestones}
            <a href={`/${username}/${reponame}/milestones`} class="tab" class:active={activeTab === 'milestones'}>
              <Calendar size={18} />
              {#if activeTab === 'milestones'}<span class="tab-label">Milestones</span>{/if}
            </a>
          {/if}

          <a href={`/${username}/${reponame}/settings`} class="tab mobile-settings-tab" class:active={activeTab === 'settings'}>
            <Settings size={18} />
            {#if activeTab === 'settings'}<span class="tab-label">Settings</span>{/if}
          </a>
        </nav>
      </div>

      <div class="master-body">
        {@render children()}
      </div>
    </div>
  {/if}
</div>
</div>

<style>
  .workspace-shell { height: 100%; width: 100%; display: flex; flex-direction: column; overflow: hidden; }
  .master-container { flex: 1; display: grid; grid-template-columns: 1fr; grid-template-rows: 1fr; overflow: hidden; }
  .skeleton-master, .master-tonal, .error-tonal { grid-area: 1 / 1 / 2 / 2; display: flex; flex-direction: column; overflow: hidden; }
  .master-tonal { z-index: 2; }
  .skeleton-master { z-index: 1; }

  .tonal-search-bar { padding: 1.5rem 1.5rem 0.75rem; background: rgba(255,255,255,0.01); }

  .repo-identity-bar { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1.5rem; }
  .repo-header-main { display: flex; align-items: center; gap: 1rem; }
  :global(.header-repo-avatar) { flex-shrink: 0; }
  .breadcrumb-nav { display: flex; align-items: baseline; gap: 0.5rem; flex-wrap: wrap; }
  .bc-item { color: var(--text-muted); text-decoration: none; font-weight: 500; font-size: 1.1rem; }
  h1.bc-item.active { color: var(--text-color); font-weight: 800; font-size: 1.25rem; margin: 0; word-break: break-all; }
  .bc-home { color: var(--text-muted); display: flex; opacity: 0.6; align-self: center; }
  .bc-sep { opacity: 0.2; }

  .repo-meta-actions { display: flex; align-items: center; gap: 1.5rem; }
  .license-tag { display: flex; align-items: center; gap: 0.4rem; font-size: 0.8rem; font-weight: 700; color: var(--text-muted); }
  .settings-link { color: var(--text-muted); display: flex; padding: 4px; transition: color 0.2s; text-decoration: none; }
  .settings-link:hover { color: var(--primary-color); }

  .clone-zone { position: relative; }
  .split-btn-group { display: flex; align-items: stretch; border-radius: 12px; overflow: hidden; box-shadow: var(--shadow-sm); transition: transform 0.2s; }
  .split-btn-group:hover { transform: translateY(-1px); }
  .main-clone-btn { border-radius: 0; padding-right: 0.75rem; border-right: 1px solid rgba(255,255,255,0.1); font-size: 0.9rem; }
  .arrow-btn { border-radius: 0; padding-left: 0.4rem; padding-right: 0.4rem; width: auto; }

  .clone-dropdown { position: absolute; top: calc(100% + 0.5rem); right: 0; width: 340px; z-index: 1000; padding: 1.25rem; background: var(--md3-surface-2); border-radius: 1.25rem; border: 1px solid var(--border-color); animation: slideUp 0.2s ease-out; }
  @keyframes slideUp { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }

  .dropdown-section { display: flex; flex-direction: column; gap: 0.75rem; }
  .section-title { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.7; letter-spacing: 0.05em; }
  .dropdown-divider { height: 1px; background: var(--border-color); margin: 0.75rem -1.25rem; opacity: 0.5; }

  .protocol-tabs { display: flex; gap: 4px; background: rgba(0,0,0,0.1); padding: 3px; border-radius: 10px; }
  .proto-btn { flex: 1; border: none; background: transparent; color: var(--text-muted); padding: 0.4rem; border-radius: 8px; font-size: 0.8rem; font-weight: 700; cursor: pointer; transition: all 0.2s; }
  .proto-btn.active { background: var(--primary-color); color: white; }

  .clone-copy-row { display: flex; align-items: center; gap: 0.5rem; background: rgba(0,0,0,0.2); padding: 0.5rem 0.75rem; border-radius: 10px; border: 1px solid var(--border-color); }
  .url-display { flex: 1; font-size: 0.8rem; font-family: 'JetBrains Mono', monospace; overflow: hidden; text-overflow: ellipsis; opacity: 0.8; }
  .copy-action { background: transparent; border: none; color: var(--primary-color); cursor: pointer; display: flex; padding: 4px; border-radius: 6px; transition: background 0.2s; }
  .copy-action:hover { background: rgba(var(--primary-rgb), 0.1); }
  .protocol-hint { font-size: 0.75rem; color: var(--text-muted); margin: 0; }

  .feature-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .feature-item { display: flex; align-items: center; gap: 0.75rem; color: var(--text-color); font-size: 0.85rem; opacity: 0.9; }
  .feature-item span:first-child { color: var(--primary-color); font-size: 1.1rem; }

  .archive-grid-compact { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; }
  .archive-tile { display: flex; align-items: center; justify-content: center; gap: 0.5rem; padding: 0.75rem; background: rgba(var(--primary-rgb), 0.05); border: 1px solid var(--border-color); border-radius: 12px; text-decoration: none; color: var(--text-color); font-size: 0.8rem; font-weight: 700; transition: all 0.2s; }
  .archive-tile:hover { background: rgba(var(--primary-rgb), 0.1); border-color: var(--primary-color); }

  .repo-header-region { background-color: var(--surface-variant); border-bottom: none; }
  
  @media (max-width: 800px) {
    .repo-header-region {
      padding-bottom: 0.25rem;
    }
  }

  .rich-tabs { display: flex; gap: 0.25rem; padding: 0.5rem 1rem 0; overflow-x: auto; scrollbar-width: none; -ms-overflow-style: none; }
  .rich-tabs::-webkit-scrollbar { display: none; }
  .tab { display: flex; align-items: center; gap: 0.6rem; padding: 0.85rem 1.25rem; color: var(--text-muted); text-decoration: none; font-size: 0.85rem; font-weight: 700; transition: all 0.2s; border-radius: 0.75rem 0.75rem 0 0; }
  .tab:hover { color: var(--text-color); background: rgba(255,255,255,0.03); }
  .tab.active { color: var(--primary-color); background: var(--bg-color); border: none; }
  .tab-count { font-variant-numeric: tabular-nums; background: rgba(0,0,0,0.05); padding: 0.1rem 0.4rem; border-radius: 4px; font-size: 0.75rem; color: var(--text-muted); }
  .tab.active .tab-count { background: var(--primary-container); color: var(--on-primary-container); }
  .tab-label { text-transform: uppercase; letter-spacing: 0.05em; font-size: 0.75rem; }

  .mobile-settings-tab { display: none; }
  
  @media (max-width: 800px) {
    .tab:not(.active) .tab-label, .tab:not(.active) .tab-count { display: none; }
    .tab { padding: 0.85rem 1rem; flex-shrink: 0; }
    .repo-identity-bar { flex-direction: row; justify-content: space-between; align-items: center; gap: 0.5rem; }
    .repo-meta-actions { width: auto; gap: 0.5rem; align-items: center; }
    .repo-meta-actions .settings-link { display: none; }
    .mobile-settings-tab { display: flex; margin-left: auto; }
    .clone-text, .arrow-btn { display: none; }
    .clone-dropdown { width: calc(100vw - 3rem); right: -1.5rem; }
  }

  .master-body { flex: 1; overflow-y: auto; background: var(--bg-color); padding: 1rem; box-sizing: border-box; }
  .master-body-skeleton { flex: 1; padding: 2rem; display: flex; flex-direction: column; background: var(--bg-color); }
  @media (max-width: 600px) {
    .master-body { padding: 0.5rem; }
  }

  @media (max-width: 480px) {
    .master-body { padding: 0.5rem; }
  }
  .skeleton-master { background: var(--md-sys-color-surface-container-low); height: 100vh; display: flex; flex-direction: column; }
  .skeleton-pulse { background: var(--md-sys-color-surface-container-high); animation: skeleton-pulse 1.5s infinite ease-in-out; border-radius: 0.5rem; }
  
  .search-skeleton { height: 44px; width: 100%; max-width: 800px; border-radius: 14px; margin: 1rem; }
  .identity-skeleton { height: 32px; width: 320px; margin: 1rem; }
  .actions-skeleton { height: 36px; width: 180px; border-radius: 12px; margin: 1rem; }
  .tabs-skeleton { height: 48px; width: 100%; border-radius: 12px 12px 0 0; margin-top: auto; }
  .body-skeleton { height: 100%; width: 100%; border-radius: 1rem; background: var(--md-sys-color-surface-container-low); }

  @keyframes skeleton-pulse { 0% { opacity: 0.4; } 50% { opacity: 0.8; } 100% { opacity: 0.4; } }
  .error-tonal { height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; color: #f87171; }
</style>
