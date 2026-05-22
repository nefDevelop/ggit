<script lang="ts">
  import RepoFileBrowser from '$lib/components/RepoFileBrowser.svelte';
  import { page } from '$app/state';
  import { getContext } from 'svelte';
  import { Book, Scale, Database } from 'lucide-svelte';

  const repoContext: any = getContext('repo');
  let repoData = $derived(repoContext.data);
  let fullRepoRes = $derived(repoContext.full);
  
  // Format bytes to MiB
  function formatBytes(bytes: number) {
    if (bytes === 0 || bytes === undefined) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KiB', 'MiB', 'GiB', 'TiB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  // Calculate size in bytes
  // If backend returns size in KB (as Gitea/Gogs sometimes do), multiply by 1024.
  // The backend's countObject.Size + countObject.SizePack is usually in KiB.
  // Wait, let's just assume bytes if the label says 'Size' or if it's KiB, let's multiply by 1024 if we know it's KiB.
  // In Gogs/Gitea, Repository.Size is in KiB. Let's do `repoData?.size * 1024`.
</script>

<svelte:head>
  <title>{page.params.username}/{page.params.reponame} · Ggit</title>
</svelte:head>

<div class="repo-grid">
  <div class="main-content">
    <RepoFileBrowser path="" />
  </div>
  
  <div class="repo-sidebar">
    <div class="sidebar-section">
      <h3 class="sidebar-title">About</h3>
      {#if repoData?.description}
        <p class="repo-description">{repoData.description}</p>
      {/if}
      
      <div class="sidebar-items">
        <!-- If there is a readme in the root, it gets rendered in RepoFileBrowser usually, or we just link to it -->
        {#if fullRepoRes?.readme || true}
        <a href="#readme" class="sidebar-item">
          <Book size={16} />
          <span>Readme</span>
        </a>
        {/if}
        
        {#if repoData?.license}
          <div class="sidebar-item">
            <Scale size={16} />
            <span>{repoData.license}</span>
          </div>
        {/if}
        
        {#if repoData?.size !== undefined}
          <div class="sidebar-item">
            <Database size={16} />
            <span>{formatBytes(repoData.size * 1024)}</span>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .repo-grid {
    display: grid;
    grid-template-columns: minmax(0, 1fr) 300px;
    gap: 1.5rem;
    align-items: start;
    max-width: 1200px;
    margin: 0 auto;
    padding-top: 1rem;
  }
  
  @media (max-width: 1280px) {
    .repo-grid {
      padding-left: 1rem;
      padding-right: 1rem;
    }
  }

  @media (max-width: 800px) {
    .repo-grid {
      grid-template-columns: 1fr;
    }
    .repo-sidebar {
      margin-top: 1rem;
    }
  }

  .sidebar-section {
    background: var(--surface-variant, rgba(0,0,0,0.02));
    border: 1px solid var(--border-color, rgba(0,0,0,0.05));
    border-radius: 12px;
    padding: 1.25rem;
  }
  
  @media (prefers-color-scheme: dark) {
    .sidebar-section {
      background: var(--surface-variant, rgba(255,255,255,0.02));
      border: 1px solid var(--border-color, rgba(255,255,255,0.1));
    }
  }

  .sidebar-title {
    font-size: 1rem;
    font-weight: 700;
    margin: 0 0 1rem 0;
    color: var(--text-color, inherit);
  }

  .repo-description {
    font-size: 0.9rem;
    color: var(--text-muted, #666);
    margin-bottom: 1rem;
    line-height: 1.4;
  }
  
  @media (prefers-color-scheme: dark) {
    .repo-description {
      color: var(--text-muted, #aaa);
    }
  }

  .sidebar-items {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .sidebar-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    color: var(--text-muted, #666);
    font-size: 0.9rem;
    text-decoration: none;
    transition: color 0.2s;
  }
  
  @media (prefers-color-scheme: dark) {
    .sidebar-item {
      color: var(--text-muted, #aaa);
    }
  }

  a.sidebar-item:hover {
    color: var(--primary-color, #3b82f6);
  }
</style>
