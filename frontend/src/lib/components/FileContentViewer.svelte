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

  let { repoData = null, username = "", reponame = "", realPath = "", showBlame = false, blameLoading = false, sidebarCollapsed = false, codeRows = [], formatSize, timeAgo, onToggleBlame, onToggleSidebar }: {
    repoData: any;
    username: string;
    reponame: string;
    realPath: string;
    showBlame: boolean;
    blameLoading: boolean;
    sidebarCollapsed: boolean;
    codeRows: any[];
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
          {#if blameLoading}<div class="mini-spinner"></div>{:else}<History size={16} /><span class="btn-text">Blame</span>{/if}
        </button>
        <a href={`/${username}/${reponame}/_edit/${repoData?.branch}/${realPath}`} class="tool-btn" title="Edit">
          <Edit size={16} /><span class="btn-text">Edit</span>
        </a>
      {/if}
      <a href={`/api/v2/${username}/${reponame}/raw/${repoData?.branch}/${realPath}`} target="_blank" class="tool-btn" title="Raw">
        <FileCode size={16} /><span class="btn-text">Raw</span>
      </a>
      <a href={`/api/v2/${username}/${reponame}/media/${repoData?.branch}/${realPath}`} class="tool-btn" download title="Download">
        <Download size={16} /><span class="btn-text">Download</span>
      </a>
      <button class="tool-btn btn-sidebar-toggle-inline" onclick={onToggleSidebar} title={sidebarCollapsed ? "Show Sidebar" : "Hide Sidebar"}>
        {#if sidebarCollapsed}<PanelLeftOpen size={18} />{:else}<PanelRightClose size={18} />{/if}
      </button>
    </div>
  </div>
  <div class="code-box">
    {#if repoData?.is_text_file}
      <div class="code-container" class:with-blame={showBlame}>
        {#each (codeRows) as row}
          <div class="code-line-row">
            {#if showBlame && row.blame}
              <div class="blame-cell" title={`${row.blame.author} - ${new Date(row.blame.time * 1000).toLocaleString()}`}>
                <span class="blame-author">{row.blame.author.split(' ')[0]}</span>
                <span class="blame-time">{timeAgo(row.blame.time * 1000)}</span>
              </div>
            {/if}
            <div class="ln-cell">{row.num}</div>
            <div class="code-cell"><pre class="hljs"><code>{@html row.content}</code></pre></div>
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
      <header class="readme-header small"><BookText size={16} /> README Preview</header>
      <div class="readme-content"><MarkdownRenderer content={repoData.readme_content} repoContext={{owner: username, name: reponame}} /></div>
    </div>
  {/if}
</div>

<style>
  .file-header { background: var(--surface-variant); padding: 0.85rem 1.25rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .file-meta { display: flex; align-items: center; gap: 0.75rem; font-weight: 800; font-size: 0.95rem; color: var(--text-color); }
  .filesize { font-weight: 400; opacity: 0.5; font-size: 0.8rem; }
  .file-tools.side-by-side { display: flex; align-items: center; gap: 0; }
  .tool-btn { background: var(--bg-color); border: 1px solid var(--md-sys-color-outline-variant); color: var(--text-muted); padding: 0 0.85rem; border-radius: 0.6rem; cursor: pointer; font-size: 0.8rem; font-weight: 800; transition: all 0.2s; text-decoration: none; display: flex; align-items: center; justify-content: center; gap: 0.5rem; height: 34px; box-sizing: border-box; }
  .tool-btn:hover { border-color: var(--primary-color); color: var(--primary-color); background: rgba(var(--primary-rgb),0.05); }
  .file-tools.side-by-side .tool-btn { border-radius: 0; border-left-width: 0; }
  .file-tools.side-by-side .tool-btn:first-child { border-top-left-radius: 0.6rem; border-bottom-left-radius: 0.6rem; border-left-width: 1px; }
  .file-tools.side-by-side .tool-btn:last-child { border-top-right-radius: 0.6rem; border-bottom-right-radius: 0.6rem; }
  .tool-btn.active { background: var(--primary-color) !important; color: white !important; }
  .mini-spinner { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 0.8s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  .btn-text { font-size: 0.8rem; }
  .code-box { background: var(--md-sys-color-surface-container-lowest); overflow-x: auto; border: 1px solid var(--md-sys-color-outline-variant); margin: 0; border-radius: 0.75rem; }
  .code-container { display: flex; flex-direction: column; }
  .code-line-row { display: flex; align-items: stretch; min-height: 1.6rem; }
  .code-line-row:hover { background: rgba(var(--primary-rgb),0.05); }
  .ln-cell { width: 60px; min-width: 60px; text-align: right; padding-right: 1.25rem; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; line-height: 1.6; color: var(--text-muted); opacity: 0.4; user-select: none; border-right: 1px solid var(--md-sys-color-outline-variant); background: var(--md-sys-color-surface-container-low); display: flex; align-items: flex-start; justify-content: flex-end; }
  .code-cell { flex: 1; padding: 0 1rem; display: flex; align-items: flex-start; overflow-x: auto; }
  .code-cell pre { margin: 0; padding: 0; font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; line-height: 1.6; background: transparent; }
  .code-cell pre code { background: transparent; padding: 0; }
  .blame-cell { width: 140px; min-width: 140px; border-right: 1px solid var(--md-sys-color-outline-variant); background: var(--md-sys-color-surface-container-low); padding: 0 0.75rem; font-size: 0.75rem; display: flex; align-items: center; gap: 0.5rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; opacity: 0.6; }
  .blame-author { font-weight: 800; color: var(--primary-color); }
  .blame-time { color: var(--text-muted); font-size: 0.7rem; }
  .binary-placeholder { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 4rem; gap: 1rem; color: var(--text-muted); }
  .file-readme-integrated { border-top: 1px solid var(--md-sys-color-outline-variant); padding: 1rem; }
  .readme-header.small { display: flex; align-items: center; gap: 0.5rem; font-weight: 800; font-size: 0.85rem; text-transform: uppercase; color: var(--text-muted); padding: 0.5rem 0; margin-bottom: 1rem; }
  .readme-content { padding: 0 1rem; }
  @media (max-width: 600px) { .code-box { margin: 0.5rem; } .file-header { padding: 0.6rem 0.75rem; } .tool-btn .btn-text { display: none; } .tool-btn { padding: 0 0.5rem; } }
  @media (max-width: 480px) { .code-box { margin: 0; border-radius: 0; border-left: none; border-right: none; } .code-cell { padding: 0 0.5rem; } .ln-cell { width: 45px; min-width: 45px; padding-right: 0.5rem; } }
</style>
