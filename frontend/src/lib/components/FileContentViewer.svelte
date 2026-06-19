<script lang="ts">
  import { goto } from '$app/navigation';
  import File from "lucide-svelte/icons/file";
  import History from "lucide-svelte/icons/history";
  import Edit from "lucide-svelte/icons/edit";
  import FileCode from "lucide-svelte/icons/file-code";
  import Download from "lucide-svelte/icons/download";
  import PanelLeftOpen from "lucide-svelte/icons/panel-left-open";
  import PanelRightClose from "lucide-svelte/icons/panel-right-close";
  import Binary from "lucide-svelte/icons/binary";
  import BookText from "lucide-svelte/icons/book-text";
  import MarkdownRenderer from "./ui/MarkdownRenderer.svelte";

  let { repoData = null, username = "", reponame = "", realPath = "", showBlame = false, blameLoading = false, sidebarCollapsed = false, formatSize, timeAgo, onToggleBlame, onToggleSidebar }: {
    repoData: any;
    username: string;
    reponame: string;
    realPath: string;
    showBlame: boolean;
    blameLoading: boolean;
    sidebarCollapsed: boolean;
    formatSize: (n: number) => string;
    timeAgo: (d: any) => string;
    onToggleBlame: () => void;
    onToggleSidebar: () => void;
  } = $props();
</script>

<div class="file-viewer">
  <div class="file-header">
    <div class="file-meta">
      <File size={18} />
      <span class="filename">{repoData?.name}</span>
      <span class="filesize">{formatSize(repoData?.size || 0)}</span>
    </div>
    <div class="file-tools side-by-side">
      {#if repoData?.is_text_file}
        <button class="tool-btn" class:active={showBlame} onclick={onToggleBlame} disabled={blameLoading}>
          {#if blameLoading}
            <div class="mini-spinner"></div>
          {:else}
            <History size={16} />
            <span class="btn-text">Blame</span>
          {/if}
        </button>
        <a href={`/${username}/${reponame}/_edit/${repoData?.branch}/${realPath}`} class="tool-btn" title="Edit" aria-label="Edit">
          <Edit size={16} /><span class="btn-text">Edit</span>
        </a>
      {/if}
      <a href={`/api/v2/${username}/${reponame}/raw/${repoData?.branch}/${realPath}`} target="_blank" class="tool-btn" title="Raw" aria-label="Raw">
        <FileCode size={16} /><span class="btn-text">Raw</span>
      </a>
      <a href={`/api/v2/${username}/${reponame}/media/${repoData?.branch}/${realPath}`} class="tool-btn" download title="Download" aria-label="Download">
        <Download size={16} /><span class="btn-text">Download</span>
      </a>
      <button class="tool-btn btn-sidebar-toggle-inline" onclick={onToggleSidebar} title={sidebarCollapsed ? "Show Sidebar" : "Hide Sidebar"} aria-label={sidebarCollapsed ? "Show Sidebar" : "Hide Sidebar"}>
        {#if sidebarCollapsed}
          <PanelLeftOpen size={18} />
        {:else}
          <PanelRightClose size={18} />
        {/if}
      </button>
    </div>
  </div>

  <div class="code-box">
    {#if repoData?.is_text_file}
      <div class="code-container" class:with-blame={showBlame}>
        {#each (repoData?.codeRows || []) as row}
          <div class="code-line-row">
            {#if showBlame && row.blame}
              <div class="blame-cell" title={`${row.blame.author} - ${new Date(row.blame.time * 1000).toLocaleString()}`}>
                <span class="blame-author">{row.blame.author.split(' ')[0]}</span>
                <span class="blame-time">{timeAgo(row.blame.time * 1000)}</span>
              </div>
            {/if}
            <div class="ln-cell">{row.num}</div>
            <div class="code-cell">
              <pre class="hljs"><code>{@html row.content}</code></pre>
            </div>
          </div>
        {/each}
      </div>
    {:else}
      <div class="binary-placeholder">
        <Binary size={48} />
        <p>Binary file — {formatSize(repoData?.size || 0)}</p>
      </div>
    {/if}
  </div>

  {#if repoData?.readme_content}
    <div class="file-readme-integrated">
      <header class="readme-header small">
        <BookText size={16} />
        README Preview
      </header>
      <div class="readme-content">
        <MarkdownRenderer content={repoData.readme_content} repoContext={{owner: username, name: reponame}} />
      </div>
    </div>
  {/if}
</div>

<style>
  .file-viewer { display: flex; flex-direction: column; }
  .file-header {
    display: flex; justify-content: space-between; align-items: center;
    padding: 0.75rem 1rem; background: var(--md-sys-color-surface-container);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    flex-wrap: wrap; gap: 0.5rem;
  }
  .file-meta { display: flex; align-items: center; gap: 0.75rem; }
  .filename { font-weight: 700; font-family: "JetBrains Mono", monospace; }
  .filesize { color: var(--text-muted); font-size: 0.85rem; }
  .file-tools { display: flex; gap: 4px; }
  .tool-btn {
    display: flex; align-items: center; gap: 0.35rem;
    padding: 0.4rem 0.7rem; border-radius: 0.4rem;
    border: none; background: transparent; color: var(--text-muted);
    cursor: pointer; font-size: 0.8rem; font-weight: 600;
    text-decoration: none; transition: all 0.2s;
  }
  .tool-btn:hover { background: var(--md-sys-color-surface-container-high); color: var(--text-color); }
  .tool-btn.active { background: var(--primary-container); color: var(--on-primary-container); }
  .code-box { overflow-x: auto; background: var(--md-sys-color-surface-container-lowest); }
  .code-container { min-width: fit-content; }
  .code-line-row { display: flex; font-family: "JetBrains Mono", monospace; font-size: 13px; line-height: 1.6; }
  .code-line-row:hover { background: var(--md-sys-color-surface-container); }
  .blame-cell {
    display: flex; gap: 0.5rem; padding: 0 0.5rem;
    width: 180px; min-width: 180px; cursor: default;
    color: var(--text-muted); font-size: 12px; user-select: none;
    border-right: 1px solid var(--md-sys-color-outline-variant);
    overflow: hidden; white-space: nowrap;
  }
  .blame-author { font-weight: 600; overflow: hidden; text-overflow: ellipsis; }
  .blame-time { font-weight: 400; opacity: 0.7; }
  .ln-cell {
    width: 40px; min-width: 40px; text-align: right; padding-right: 1rem;
    color: var(--text-muted); user-select: none; opacity: 0.5;
    border-right: 1px solid var(--md-sys-color-outline-variant);
  }
  .code-cell { flex: 1; padding: 0 0 0 1rem; }
  .code-cell pre { margin: 0; background: transparent; padding: 0; }
  .binary-placeholder {
    display: flex; flex-direction: column; align-items: center;
    justify-content: center; padding: 4rem; gap: 1rem; color: var(--text-muted);
  }
  .file-readme-integrated {
    border-top: 1px solid var(--md-sys-color-outline-variant); padding: 1rem;
  }
  .readme-header {
    display: flex; align-items: center; gap: 0.5rem;
    padding: 0.75rem 1rem; font-weight: 800; font-size: 0.85rem;
    text-transform: uppercase; color: var(--text-muted);
  }
  .readme-header.small { padding: 0.5rem 0; margin-bottom: 1rem; }
  .readme-content { padding: 0 1rem; }
</style>
