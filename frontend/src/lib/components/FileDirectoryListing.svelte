<script lang="ts">
  import { goto } from '$app/navigation';
  import Folder from "lucide-svelte/icons/folder";
  import FileText from "lucide-svelte/icons/file-text";

  let { entries = [], username = "", reponame = "", branch = "", realPath = "", timeAgo }: {
    entries: any[];
    // Optional because it has a default. SvelteKit types route params as
    // string | undefined, so a required string prop rejected every caller.
    username?: string;
    reponame?: string;
    branch: string;
    realPath: string;
    timeAgo: (d: any) => string;
  } = $props();
</script>

<table class="file-table">
  <tbody>
    {#each entries as entry}
      <tr class="file-row clickable" onclick={() => goto(`/${username}/${reponame}/src/${branch}/${realPath ? realPath + '/' : ''}${entry.name}`)}>
        <td class="name-cell">
          {#if entry.type === 'tree'}
            <Folder size={18} class="icon-tree" />
          {:else}
            <FileText size={18} class="icon-blob" />
          {/if}
          <a href={`/${username}/${reponame}/src/${branch}/${realPath ? realPath + '/' : ''}${entry.name}`} onclick={(e) => e.preventDefault()}>
            {entry.name}
          </a>
        </td>
        <td class="msg-cell">
          <span class="msg-text">{entry.last_commit?.summary || ''}</span>
        </td>
        <td class="date-cell">{timeAgo(entry.last_commit?.author?.when) || ''}</td>
      </tr>
    {/each}
  </tbody>
</table>

<style>
  .file-table { width: 100%; border-collapse: collapse; }
  .file-row { border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .file-row:last-child { border-bottom: none; }
  .file-row td { padding: 0.85rem 1.25rem; vertical-align: middle; }
  .file-row.clickable { cursor: pointer; }
  .file-row:hover { background: rgba(var(--primary-rgb), 0.04); }
  .name-cell { display: flex; align-items: center; gap: 1rem; min-width: 200px; }
  .name-cell a { color: var(--text-color); text-decoration: none; font-weight: 700; font-size: 0.95rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .name-cell a:hover { color: var(--primary-color); }
  :global(.icon-tree) { color: var(--status-warning-soft); }
  :global(.icon-blob) { color: var(--text-muted); opacity: 0.7; }
  .msg-cell { font-size: 0.85rem; color: var(--text-muted); max-width: 0; width: 60%; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .msg-text { opacity: 0.7; }
  .date-cell { text-align: right; font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; white-space: nowrap; }
  @media (max-width: 600px) { .msg-cell, .date-cell { display: none; } .name-cell { width: 100%; min-width: 0; } .file-row td { padding: 0.85rem 1rem; } }
</style>
