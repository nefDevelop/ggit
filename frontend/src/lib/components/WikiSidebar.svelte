<script lang="ts">
  import type { Snippet } from 'svelte';
  import { slide } from 'svelte/transition';
  import Button from "./ui/Button.svelte";
  import Plus from "lucide-svelte/icons/plus";
  import ChevronRight from "lucide-svelte/icons/chevron-right";
  import Folder from "lucide-svelte/icons/folder";
  import FileText from "lucide-svelte/icons/file-text";

  let { wikiPages = [], currentPath = "", expandedFolders = {}, repoOwner = "", repoName = "", onNewPage, onToggleFolder }: {
    wikiPages: any[];
    currentPath: string;
    expandedFolders: Record<string, boolean>;
    repoOwner: string;
    repoName: string;
    onNewPage?: () => void;
    onToggleFolder?: (path: string) => void;
  } = $props();
</script>

<aside class="wiki-sidebar">
  <div class="sidebar-header-minimal">
    <Button class="w-full" onclick={onNewPage} icon={Plus}>
      New Page
    </Button>
  </div>

  <nav class="wiki-nav">
    {#snippet renderTree(items: any[], depth = 0, parentPath = '')}
      {#each items as pageItem}
        {@const fullPath = parentPath ? `${parentPath}/${pageItem.name}` : pageItem.name}
        {#if pageItem.is_folder}
          <div class="folder-section">
            <button class="nav-item folder-label" style="padding-left: {depth * 0.75 + 0.8}rem"
              onclick={() => onToggleFolder?.(fullPath)}>
              <ChevronRight size={14} class="chevron {expandedFolders[fullPath] ? 'open' : ''}" />
              <Folder size={14} class="opacity-60" />
              <span>{pageItem.title}</span>
            </button>
            {#if expandedFolders[fullPath]}
              <div class="folder-children" transition:slide>
                {@render renderTree(pageItem.children || [], depth + 1, fullPath)}
              </div>
            {/if}
          </div>
        {:else}
          <a href="/{repoOwner}/{repoName}/wiki/{pageItem.name}" class="nav-item" class:active={currentPath === pageItem.name}
            style="padding-left: {depth * 0.75 + 0.8}rem">
            <FileText size={14} class="opacity-60" />
            <span>{pageItem.title || pageItem.name}</span>
          </a>
        {/if}
      {/each}
    {/snippet}
  </nav>
</aside>

<style>
  .wiki-sidebar { width: 280px; flex-shrink: 0; display: flex; flex-direction: column; gap: 1rem; }
  .sidebar-header-minimal { padding: 0 0.5rem; }
  .wiki-nav { display: flex; flex-direction: column; gap: 0.25rem; overflow-y: auto; max-height: calc(100vh - 200px); }
  .nav-item {
    display: flex; align-items: center; gap: 0.5rem;
    padding: 0.5rem 0.75rem; border-radius: 0.5rem;
    color: var(--text-muted); text-decoration: none;
    font-size: 0.9rem; transition: all 0.2s; cursor: pointer;
    background: transparent; border: none; width: 100%;
    text-align: left; font-family: inherit;
  }
  .nav-item:hover { background: rgba(255,255,255,0.03); color: var(--text-color); }
  .nav-item.active { background: var(--md-sys-color-primary-container); color: var(--md-sys-color-on-primary-container); }
  .folder-section { display: flex; flex-direction: column; }
  .folder-label { font-weight: 700; color: var(--text-color); }
  .chevron { transition: transform 0.2s; }
  .chevron.open { transform: rotate(90deg); }
  .folder-children { display: flex; flex-direction: column; }
  .opacity-60 { opacity: 0.6; }
  .w-full { width: 100%; justify-content: center; }
</style>
