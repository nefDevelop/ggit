<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import Folder from "lucide-svelte/icons/folder";
  import FileText from "lucide-svelte/icons/file-text";

  let { entries = [], username = "", reponame = "", branch = "", realPath = "", timeAgo }: {
    entries: any[];
    username: string;
    reponame: string;
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
  .file-row { cursor: pointer; transition: background 0.15s; }
  .file-row:hover { background: var(--md-sys-color-surface-container-high); }
  .name-cell { padding: 0.6rem 1rem; display: flex; align-items: center; gap: 0.75rem; }
  .name-cell a { color: var(--text-color); text-decoration: none; font-weight: 600; }
  .name-cell a:hover { color: var(--primary-color); }
  .icon-tree { color: #42a5f5; }
  .icon-blob { color: var(--text-muted); }
  .msg-cell { padding: 0.6rem 0.5rem; max-width: 400px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .msg-text { color: var(--text-muted); font-size: 0.85rem; }
  .date-cell { padding: 0.6rem 1rem; text-align: right; white-space: nowrap; color: var(--text-muted); font-size: 0.85rem; }
  @media (max-width: 700px) { .msg-cell, .date-cell { display: none; } }
</style>
