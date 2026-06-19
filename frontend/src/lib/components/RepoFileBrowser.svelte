<script lang="ts">
  import AlertCircle from "lucide-svelte/icons/alert-circle";
  import Terminal from "lucide-svelte/icons/terminal";
  import PanelLeftOpen from "lucide-svelte/icons/panel-left-open";
  import PanelRightClose from "lucide-svelte/icons/panel-right-close";
  import BookText from "lucide-svelte/icons/book-text";
  import Copy from "lucide-svelte/icons/copy";
  import Check from "lucide-svelte/icons/check";

  import { apiFetch } from "$lib/api";
  import { page } from '$app/state';
  import { getContext, onMount } from 'svelte';
  import type { RepoContext } from "$lib/types";
  import RepoSearchResults from './RepoSearchResults.svelte';
  import CommitStatusIcon from './CommitStatusIcon.svelte';
  import Avatar from './Avatar.svelte';
  import MarkdownRenderer from './ui/MarkdownRenderer.svelte';
  import { getLanguageColor } from '$lib/utils/languages';
  import hljs from "$lib/hljs";
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

  function copyToClipboard(text: string, id: string) {
    navigator.clipboard.writeText(text);
    copiedId = id;
    setTimeout(() => { if (copiedId === id) copiedId = null; }, 2000);
  }

  let cloneUrl = $derived(repoData?.clone_url || repoData?.repository?.clone_url || '');

  let abortController = new AbortController();

  function formatSize(bytes: number) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  function timeAgo(dateInput: any) {
    if (!dateInput) return '';
    const date = new Date(dateInput);
    const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
    if (seconds < 60) return 'just now';
    const minutes = Math.floor(seconds / 60);
    if (minutes < 60) return `${minutes}m ago`;
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    if (days < 30) return `${days}d ago`;
    return date.toLocaleDateString();
  }

  let realPath = $derived(() => {
    let p = path;
    if (p.startsWith('src/')) { const parts = p.split('/'); return parts.slice(2).join('/'); }
    return p;
  });

  let codeRows = $derived(() => {
    if (!repoData?.content || !repoData?.is_text_file) return [];
    let lines: string[];
    if (repoData.highlighted_content) {
      lines = repoData.highlighted_content.split('\n');
    } else {
      const language = repoData.name.split('.').pop() || '';
      let html: string;
      if (repoData.content.length > 100000) {
        html = repoData.content.replace(/[&<>"']/g, (m: string) => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":"&#39;"}[m] || m));
      } else {
        try {
          html = hljs.highlight(repoData.content, { language: hljs.getLanguage(language) ? language : 'plaintext' }).value;
        } catch (e) {
          try { html = hljs.highlightAuto(repoData.content).value; } catch (e2) {
            html = repoData.content.replace(/[&<>"']/g, (m: string) => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":"&#39;"}[m] || m));
          }
        }
      }
      lines = html.split('\n');
    }
    return lines.map((content: string, i: number) => ({
      num: i + 1,
      content: content || '&nbsp;',
      blame: showBlame ? blameData[i] : null,
    }));
  });

  const langEntries = $derived(() => {
    const langs = repoData?.languages || repoContext?.languages;
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
  function toggleSidebar() { sidebarCollapsed = !sidebarCollapsed; }

  let blameData = $state<any[]>([]);
  let showBlame = $state(false);
  let blameLoading = $state(false);

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

  $effect(() => { fetchRepo(); });

  async function fetchRepo() {
    const currentPath = realPath();
    if (lastFetched.user === username && lastFetched.repo === reponame && lastFetched.path === currentPath) return;
    lastFetched = { user: username || '', repo: reponame || '', path: currentPath || '' };
    abortController.abort();
    abortController = new AbortController();
    loading = true; error = '';
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/view/${currentPath}`, { signal: abortController.signal });
      if (res.ok) repoData = await res.json();
      else error = 'Repository content not found.';
    } catch (e) { error = 'Failed to load repository.'; }
    finally { loading = false; }
  }

  onMount(() => {
    if (!path && repoContext?.full) {
      repoData = repoContext.full;
      loading = false;
      lastFetched = { user: username, repo: reponame, path: '' };
    }
    window.addEventListener('ggit:search', (e: any) => { searchQuery = e.detail; });
    return () => { abortController.abort(); };
  });

</script>

<div class="workspace-area">
  {#if loading}
    <div class="technical-grid skeleton-mode">
      <div class="main-column">
        <div class="skeleton-pulse last-commit-skeleton"></div>
        <div class="content-view skeleton-table-view">{#each Array(10) as _}<div class="skeleton-row"></div>{/each}</div>
      </div>
    </div>
  {:else if error}
    <div class="empty-state error-tonal"><AlertCircle size={48} /><p>{error}</p></div>
  {:else if repoData?.is_bare}
    <div class="empty-state"><Terminal size={48} /><h3>Empty Repository</h3>
      <p>This repository is empty. Push your code to get started.</p>
      <div class="push-guide tonal">
        <div class="guide-copy-row">
          <code class="copy-text">git init</code>
          <button class="copy-btn" onclick={() => copyToClipboard('git init', 'init')}>
            {#if copiedId === 'init'}<Check size={16} />{:else}<Copy size={16} />{/if}
          </button>
        </div>
        <pre><code>git remote add origin {cloneUrl}
git branch -M main
git push -u origin main</code></pre>
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
              <button class="btn-sidebar-toggle" onclick={toggleSidebar} title={sidebarCollapsed ? "Show Sidebar" : "Hide Sidebar"} aria-label={sidebarCollapsed ? "Show Sidebar" : "Hide Sidebar"}>
                {#if sidebarCollapsed}<PanelLeftOpen size={20} />{:else}<PanelRightClose size={20} />{/if}
              </button>
            </div>
          </div>
        {/if}

        <div class="content-view">
          {#if repoData.type === 'directory'}
            <FileDirectoryListing entries={repoData.entries} {username} {reponame} branch={repoData.branch} realPath={realPath()} {timeAgo} />
          {:else if repoData.type === 'file'}
            <FileContentViewer {repoData} {username} {reponame} realPath={realPath()} {showBlame} {blameLoading} {sidebarCollapsed} {formatSize} {timeAgo} onToggleBlame={toggleBlame} onToggleSidebar={toggleSidebar} />
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
</div>

<style>
  .workspace-area { height: 100%; display: flex; flex-direction: column; }
  .technical-grid {
    display: flex; gap: 2rem; align-items: flex-start; flex: 1;
    padding: 1rem 1.5rem; min-height: 0;
  }
  .technical-grid.sidebar-collapsed { gap: 0; }
  .main-column { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-mode { opacity: 0.5; pointer-events: none; }
  .skeleton-pulse { background: var(--md-sys-color-surface-container-high); border-radius: 0.75rem; animation: pulse 1.5s infinite; }
  .last-commit-skeleton { height: 60px; }
  .skeleton-table-view { display: flex; flex-direction: column; gap: 0.5rem; }
  .skeleton-row { height: 40px; background: var(--md-sys-color-surface-container-high); border-radius: 0.5rem; }
  @keyframes pulse { 50% { opacity: 0.5; } }
  .empty-state {
    display: flex; flex-direction: column; align-items: center;
    justify-content: center; padding: 4rem 2rem; gap: 1rem;
    opacity: 0.7; text-align: center; height: 100%;
  }
  .empty-state h3 { margin: 0; font-size: 1.5rem; }
  .empty-state p { margin: 0; color: var(--text-muted); max-width: 400px; }
  .error-tonal { background: rgba(var(--error-rgb),0.05); border-radius: 1rem; }
  .push-guide {
    background: var(--md-sys-color-surface-container-high);
    padding: 1.5rem; border-radius: 1rem; text-align: left;
    max-width: 500px; margin-top: 1rem;
  }
  .guide-copy-row { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.5rem; }
  .guide-copy-row code { flex: 1; padding: 0.5rem; background: rgba(0,0,0,0.2); border-radius: 0.5rem; font-family: "JetBrains Mono", monospace; }
  .copy-btn { background: transparent; border: none; color: var(--text-muted); cursor: pointer; padding: 0.5rem; border-radius: 0.5rem; }
  .copy-btn:hover { background: var(--md-sys-color-surface-container); color: var(--text-color); }
  .push-guide pre { margin: 0; }
  .push-guide pre code { font-family: "JetBrains Mono", monospace; font-size: 0.85rem; line-height: 1.5; }
  .last-commit-bar {
    display: flex; align-items: center; gap: 1rem;
    padding: 0.75rem 1rem;
    background: var(--md-sys-color-surface-container);
    border-radius: 0.75rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  .commit-info { flex: 1; min-width: 0; }
  .commit-summary-row { margin-bottom: 0.25rem; }
  .commit-msg-link { text-decoration: none; color: inherit; }
  .commit-msg-link:hover .commit-msg { color: var(--primary-color); }
  .commit-msg { font-weight: 700; font-size: 0.9rem; }
  .commit-meta-row { display: flex; align-items: center; gap: 0.5rem; font-size: 0.8rem; color: var(--text-muted); }
  .commit-meta-row .author { font-weight: 600; }
  .commit-meta-row .separator { opacity: 0.6; }
  .commit-actions-right { display: flex; align-items: center; gap: 0.5rem; }
  .sha-link { text-decoration: none; }
  .sha {
    font-family: "JetBrains Mono", monospace; font-size: 0.85rem;
    color: var(--primary-color); font-weight: 700;
  }
  .btn-sidebar-toggle {
    background: transparent; border: none; color: var(--text-muted);
    cursor: pointer; padding: 0.5rem; border-radius: 0.5rem; display: flex;
  }
  .btn-sidebar-toggle:hover { background: var(--md-sys-color-surface-container-high); color: var(--text-color); }
  .content-view { border: 1px solid var(--md-sys-color-outline-variant); border-radius: 0.75rem; overflow: hidden; }
  .readme-section { border: 1px solid var(--md-sys-color-outline-variant); border-radius: 0.75rem; overflow: hidden; }
  .readme-header {
    display: flex; align-items: center; gap: 0.5rem;
    padding: 0.75rem 1rem; font-weight: 800; font-size: 0.85rem;
    text-transform: uppercase; color: var(--text-muted);
    background: var(--md-sys-color-surface-container);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }
  .readme-content { padding: 1rem; }
  @media (max-width: 900px) {
    .technical-grid { flex-direction: column; }
    .last-commit-bar { flex-wrap: wrap; }
  }
</style>
