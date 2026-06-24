<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { goto } from '$app/navigation';
  import { onMount, getContext } from 'svelte';
  import type { RepoContext } from "$lib/types";
  import { getLanguageColor } from '$lib/utils/languages';
  import RepoSearchResults from './RepoSearchResults.svelte';
  import MarkdownRenderer from './ui/MarkdownRenderer.svelte';
  import CommitStatusIcon from './CommitStatusIcon.svelte';
  import Avatar from './Avatar.svelte';
  import hljs from "$lib/hljs";
  import AlertCircle from "lucide-svelte/icons/alert-circle";
  import Terminal from "lucide-svelte/icons/terminal";
  import History from "lucide-svelte/icons/history";
  import PanelLeftOpen from "lucide-svelte/icons/panel-left-open";
  import PanelRightClose from "lucide-svelte/icons/panel-right-close";
  import Folder from "lucide-svelte/icons/folder";
  import FileText from "lucide-svelte/icons/file-text";
  import File from "lucide-svelte/icons/file";
  import Edit from "lucide-svelte/icons/edit";
  import Download from "lucide-svelte/icons/download";
  import Binary from "lucide-svelte/icons/binary";
  import BookText from "lucide-svelte/icons/book-text";
  import Rocket from "lucide-svelte/icons/rocket";
  import ChevronRight from "lucide-svelte/icons/chevron-right";
  import Copy from "lucide-svelte/icons/copy";
  import Check from "lucide-svelte/icons/check";
  import FileCode from "lucide-svelte/icons/file-code";
  import Scale from "lucide-svelte/icons/scale";
  import Database from "lucide-svelte/icons/database";
  import FileDirectoryListing from './FileDirectoryListing.svelte';
  import FileContentViewer from './FileContentViewer.svelte';
  import RepoFileSidebar from './RepoFileSidebar.svelte';

  let { path = '' } = $props();
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);
  const repoContext = getContext<RepoContext>('repo');

  let repoData = $state<any>(null);
  let loading = $state(true);
  let error = $state('');
  let searchQuery = $state('');
  let lastFetched = $state({ user: '', repo: '', path: '' });
  let copiedId = $state<string | null>(null);
  let abortController = new AbortController();

  function copyToClipboard(text: string, id: string) {
    navigator.clipboard.writeText(text);
    copiedId = id;
    setTimeout(() => { if (copiedId === id) copiedId = null; }, 2000);
  }

  let cloneUrl = $derived(repoData?.clone_url || repoData?.repository?.clone_url || '');

  $effect(() => {
    if (!path && repoContext.full) {
      repoData = repoContext.full;
      loading = false;
      lastFetched = { user: page.params.username || '', repo: page.params.reponame || '', path: '' };
    }
  });

  let codeRows = $derived.by(() => {
    if (!repoData?.content || !repoData?.is_text_file) return [];
    let lines: string[];
    if (repoData.highlighted_content) { lines = repoData.highlighted_content.split('\n'); }
    else {
      const language = repoData.name.split('.').pop() || '';
      let html: string;
      if (repoData.content.length > 100000) { html = repoData.content.replace(/[&<>"']/g, (m: string) => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":"&#39;"}[m] || m)); }
      else {
        try { html = hljs.highlight(repoData.content, { language: hljs.getLanguage(language) ? language : 'plaintext' }).value; }
        catch (e) { try { html = hljs.highlightAuto(repoData.content).value; } catch (e2) { html = repoData.content.replace(/[&<>"']/g, (m: string) => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":"&#39;"}[m] || m)); } }
      }
      lines = html.split('\n');
    }
    return lines.map((content, i) => ({ num: i + 1, content: content || '&nbsp;', blame: showBlame ? blameData[i] : null }));
  });

  function handleSearch(e: any) { searchQuery = e.detail; }

  onMount(() => {
    window.addEventListener('ggit:search', handleSearch);
    return () => { abortController.abort(); window.removeEventListener('ggit:search', handleSearch); };
  });

  let realPath = $derived(() => {
    let p = path;
    if (p.startsWith('src/')) { const parts = p.split('/'); return parts.slice(2).join('/'); }
    return p;
  });

  async function fetchRepo() {
    const currentPath = realPath();
    if (lastFetched.user === username && lastFetched.repo === reponame && lastFetched.path === currentPath) return;
    lastFetched = { user: username || '', repo: reponame || '', path: currentPath || '' };
    abortController.abort(); abortController = new AbortController();
    loading = true; error = '';
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/view/${currentPath}`, { signal: abortController.signal });
      if (res.ok) repoData = await res.json(); else error = 'Repository content not found.';
    } catch (e) { error = 'Failed to load repository.'; }
    finally { loading = false; }
  }

  $effect(() => { fetchRepo(); });

  function formatSize(bytes: number) {
    if (bytes === 0) return '0 B'; const k = 1024; const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  const langEntries = $derived(() => {
    const langs = repoData?.languages || repoContext.languages;
    if (!langs) return [];
    const total = Object.values(langs as Record<string, number>).reduce((a, b) => a + b, 0);
    if (total === 0) return [];
    return Object.entries(langs as Record<string, number>)
      .map(([name, bytes]) => ({ name, percent: ((bytes / total) * 100).toFixed(1), color: getLanguageColor(name) }))
      .sort((a, b) => parseFloat(b.percent) - parseFloat(a.percent));
  });

  let showAllLanguages = $state(false);
  let visibleLanguages = $derived(showAllLanguages ? langEntries() : langEntries().slice(0, 10));
  let sidebarCollapsed = $state(false);
  function toggleSidebar() { sidebarCollapsed = !sidebarCollapsed; localStorage.setItem('repo-sidebar-collapsed', sidebarCollapsed.toString()); }

  function timeAgo(dateInput: any) {
    if (!dateInput) return '';
    const date = new Date(dateInput);
    const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
    if (seconds < 60) return 'just now'; const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`; const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`; const days = Math.floor(hours / 24);
    if (days < 30) return `${days}d ago`;
    return date.toLocaleDateString();
  }

  onMount(() => { sidebarCollapsed = localStorage.getItem('repo-sidebar-collapsed') === 'true'; });

  let blameData = $state<any[]>([]); let showBlame = $state(false); let blameLoading = $state(false);
  async function toggleBlame() {
    if (showBlame) { showBlame = false; return; }
    if (blameData.length > 0) { showBlame = true; return; }
    blameLoading = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/blame/${realPath()}?ref=${repoData.branch}`, { signal: abortController.signal });
      if (res.ok) { blameData = await res.json(); showBlame = true; }
    } catch (e) { console.error('Failed to fetch blame', e); }
    finally { blameLoading = false; }
  }
</script>

<div class="workspace-area">
  {#if loading}
    <div class="technical-grid skeleton-mode">
      <div class="main-column">
        <div class="skeleton-pulse last-commit-skeleton"></div>
        <div class="content-view skeleton-table-view">
          {#each Array(8) as _}
            <div class="skeleton-row-item">
              <div class="skeleton-pulse icon-placeholder"></div>
              <div class="skeleton-pulse text-placeholder"></div>
              <div class="skeleton-pulse msg-placeholder"></div>
            </div>
          {/each}
        </div>
      </div>
      <aside class="side-column">
        <div class="side-section"><div class="skeleton-pulse title-placeholder"></div><div class="skeleton-pulse text-placeholder" style="height: 4rem;"></div></div>
        <div class="side-section"><div class="skeleton-pulse title-placeholder"></div><div class="skeleton-pulse lang-bar-placeholder"></div></div>
      </aside>
    </div>
  {:else if error}
    <div class="empty-state">
      <div class="empty-icon-ring"><AlertCircle size={40} /></div>
      <h2>No content found</h2>
      <p>{error || 'This repository is empty or the path does not exist.'}</p>
    </div>
  {:else if repoData}
    {#if searchQuery}
      <RepoSearchResults query={searchQuery} branch={repoData.branch || 'HEAD'} />
    {:else if repoData.is_bare}
      <div class="empty-state-push">
        <div class="empty-icon-ring"><Terminal size={40} /></div>
        <h2>Push your code</h2>
        <div class="setup-commands">
          <div class="command-block">
            <div class="command-header">
              <span class="cmd-label">Quick Setup</span>
              <button class="copy-btn" onclick={() => copyToClipboard(`git remote add origin ${cloneUrl}\ngit branch -M main\ngit push -u origin main`, 'quick')}>
                {#if copiedId === 'quick'}<Check size={16} class="success-icon" />{:else}<Copy size={16} />{/if}
              </button>
            </div>
            <pre><code>git remote add origin {cloneUrl}
git branch -M main
git push -u origin main</code></pre>
          </div>
        </div>
      </div>
    {:else}
      <div class="technical-grid" class:sidebar-collapsed={sidebarCollapsed}>
        <div class="main-column">
          {#if repoData.last_commit && repoData.type !== 'file'}
            <div class="last-commit-bar">
              <Avatar src={repoData.last_commit.author?.avatar_url || ''} username={repoData.last_commit.author?.name || 'U'}
                seed={repoData.last_commit.author?.email || repoData.last_commit.author?.name || 'unknown'} size={32} type="user" />
              <div class="commit-info">
                <div class="commit-summary-row">
                  <a href={`/${username}/${reponame}/commit/${repoData.last_commit.id}`} class="commit-msg-link">
                    <span class="commit-msg">{repoData.last_commit.summary}</span>
                  </a>
                </div>
                <div class="commit-meta-row">
                  <CommitStatusIcon status={repoData.last_commit.status} size={16} class="status-icon-inline" />
                  <span class="author">{repoData.last_commit.author?.name}</span>
                  <span class="separator">committed</span>
                  <span class="time">{timeAgo(repoData.last_commit.author?.when)}</span>
                </div>
              </div>
              <div class="commit-actions-right">
                <a href={`/${username}/${reponame}/commit/${repoData.last_commit.id}`} class="sha-link">
                  <code class="sha">{repoData.last_commit.id?.slice(0,7)}</code>
                </a>
                <button class="btn-sidebar-toggle" onclick={toggleSidebar}>
                  {#if sidebarCollapsed}<PanelLeftOpen size={20} />{:else}<PanelRightClose size={20} />{/if}
                </button>
              </div>
            </div>
          {/if}
          <div class="content-view">
            {#if repoData.type === 'directory'}
              <FileDirectoryListing entries={repoData.entries} {username} {reponame} branch={repoData.branch} realPath={realPath()} {timeAgo} />
            {:else if repoData.type === 'file'}
              <FileContentViewer {repoData} {username} {reponame} realPath={realPath()} {showBlame} {blameLoading} {sidebarCollapsed} {formatSize} {timeAgo} {codeRows} onToggleBlame={toggleBlame} onToggleSidebar={toggleSidebar} />
            {/if}
          </div>
          {#if repoData.readme}
            <div class="readme-section">
              <header class="readme-header"><BookText size={18} />{repoData.readme.name}</header>
              <div class="readme-content"><MarkdownRenderer content={repoData.readme.content} repoContext={{owner: username, name: reponame}} /></div>
            </div>
          {/if}
        </div>
        <RepoFileSidebar {repoContext} {repoData} {username} {reponame} langEntries={langEntries()} {visibleLanguages} {showAllLanguages} {formatSize} onToggleLanguages={() => showAllLanguages = !showAllLanguages} />
      </div>
    {/if}
  {/if}
</div>

<style>
  .workspace-area { height: 100%; }
  .technical-grid { display: grid; grid-template-columns: 1fr 280px; gap: 2.5rem; transition: all 0.2s; padding: 0; max-width: 100%; width: 100%; box-sizing: border-box; }
  .technical-grid.sidebar-collapsed { grid-template-columns: 100%; gap: 0; }
  .main-column { display: flex; flex-direction: column; min-width: 0; }
  .last-commit-bar { display: flex; align-items: center; gap: 1.25rem; padding: 1rem 1.25rem; background: var(--md-sys-color-secondary-container); border-radius: 1rem; margin-bottom: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); min-width: 0; }
  .commit-info { display: flex; flex-direction: column; gap: 0.2rem; flex: 1; min-width: 0; }
  .commit-summary-row { display: flex; align-items: center; min-width: 0; }
  .commit-msg { font-weight: 700; font-size: 0.95rem; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; display: block; }
  .commit-meta-row { display: flex; align-items: center; gap: 0.4rem; font-size: 0.8rem; color: var(--text-muted); opacity: 0.8; }
  .author { font-weight: 800; color: var(--text-color); }
  .separator { opacity: 0.5; }
  .commit-msg-link, .sha-link { text-decoration: none; color: inherit; min-width: 0; }
  .commit-msg-link:hover .commit-msg { color: var(--primary-color); text-decoration: underline; }
  :global(.status-icon-inline) { margin-right: 0.25rem; }
  .sha { font-family: monospace; background: var(--md-sys-color-tertiary-container); padding: 0.2rem 0.5rem; border-radius: 4px; color: var(--md-sys-color-on-tertiary-container); font-size: 0.85rem; }
  .sha-link:hover .sha { background: var(--tertiary); color: var(--on-tertiary); }
  .commit-actions-right { display: flex; align-items: center; gap: 1rem; flex-shrink: 0; }
  .btn-sidebar-toggle { background: transparent; border: none; color: var(--text-muted); cursor: pointer; display: flex; padding: 6px; border-radius: 8px; transition: all 0.2s; }
  .btn-sidebar-toggle:hover { background: rgba(var(--primary-rgb),0.1); color: var(--primary-color); }
  .content-view { background: var(--md-sys-color-secondary-container); border: 3px solid var(--md-sys-color-outline-variant); border-radius: 1rem; overflow: hidden; box-shadow: var(--shadow-lg); }
  .readme-section { margin-top: 2rem; background: var(--bg-color); border: 1px solid var(--md-sys-color-outline-variant); border-radius: 1.25rem; overflow: hidden; box-shadow: var(--shadow-sm); }
  .readme-header { padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); display: flex; align-items: center; gap: 0.75rem; font-weight: 800; font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); background: var(--surface-variant); }
  .readme-content { padding: 2rem; line-height: 1.6; }
  .skeleton-mode { pointer-events: none; }
  .skeleton-pulse { background: rgba(var(--primary-rgb),0.05); animation: pulse 1.5s infinite ease-in-out; border-radius: 0.5rem; }
  .last-commit-skeleton { height: 60px; width: 100%; margin-bottom: 1.5rem; border-radius: 1rem; }
  .skeleton-table-view { background: transparent; border-color: var(--md-sys-color-outline-variant); }
  .skeleton-row-item { display: flex; align-items: center; gap: 1rem; padding: 1rem 1.25rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .icon-placeholder { width: 20px; height: 20px; flex-shrink: 0; }
  .text-placeholder { height: 1rem; width: 150px; flex-shrink: 0; }
  .msg-placeholder { height: 0.8rem; width: 100%; opacity: 0.5; }
  .title-placeholder { height: 1rem; width: 100px; margin-bottom: 1rem; }
  .lang-bar-placeholder { height: 12px; width: 100%; border-radius: 6px; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.7; } 100% { opacity: 0.3; } }
  .empty-state { padding: 6rem 2rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1.25rem; color: var(--text-muted); }
  .empty-state h2 { margin: 0; font-size: 1.75rem; font-weight: 950; color: var(--text-color); }
  .empty-state p { max-width: 400px; margin: 0; font-size: 1.1rem; line-height: 1.5; }
  .empty-icon-ring { width: 80px; height: 80px; border: 2px dashed var(--md-sys-color-outline-variant); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem; color: var(--primary-color); background: rgba(var(--primary-rgb),0.05); }
  .empty-state-push .empty-icon-ring { border-color: var(--tertiary); color: var(--tertiary); margin-bottom: 2rem; }
  .empty-state-push { padding: 5rem 2rem; text-align: center; }
  .command-block { background: #000; padding: 1.5rem; border-radius: 1rem; text-align: left; margin-top: 2.5rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .cmd-label { font-size: 0.75rem; font-weight: 800; color: var(--primary-color); display: block; text-transform: uppercase; }
  .command-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem; }
  .copy-btn { background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); color: #fff; width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; cursor: pointer; transition: all 0.2s; }
  .copy-btn:hover { background: rgba(255,255,255,0.1); border-color: var(--primary-color); color: var(--primary-color); }
  .copy-btn .success-icon { color: #4ade80; }
  @media (max-width: 900px) { .technical-grid { grid-template-columns: 1fr; gap: 1.5rem; padding: 0; } .last-commit-bar { display: none; } .mobile-only { display: flex; } :global(.desktop-only) { display: none; } :global(.side-column) { padding-top: 1.5rem; } }
</style>
