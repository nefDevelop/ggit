<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import MarkdownRenderer from "$lib/components/ui/MarkdownRenderer.svelte";

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let releases = $state<any[]>([]);
  let loading = $state(true);

  async function fetchReleases() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/releases`);
      if (res.ok) releases = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  function formatDate(d: string) {
    return new Date(d).toLocaleDateString(undefined, {
      day: "numeric",
      month: "short",
      year: "numeric",
    });
  }

  onMount(fetchReleases);
</script>

<div class="release-manager">
  <header class="view-header">
    <div class="title-group">
      <span class="material-symbols-rounded">sell</span>
      <h2>Releases</h2>
    </div>
    <a href="/{username}/{reponame}/releases/new" class="btn btn-primary">
      Draft a new release
    </a>
  </header>

  {#if loading}
    <div class="skeleton-list">
      {#each Array(3) as _}
        <div class="skeleton-item"></div>
      {/each}
    </div>
  {:else if releases.length === 0}
    <div class="empty-state tonal">
      <span class="material-symbols-rounded">new_releases</span>
      <h3>No releases yet</h3>
      <p>
        Releases provide a way to package software and provide binary assets to
        your users.
      </p>
    </div>
  {:else}
    <div class="release-list">
      {#each releases as rel}
        <article class="release-card hover-lift" class:draft={rel.IsDraft}>
          <div class="release-meta">
            <div class="tag-status">
              <span class="tag-badge">
                <span class="material-symbols-rounded">label</span>
                {rel.TagName}
              </span>
              {#if rel.IsDraft}
                <span class="status-badge draft">Draft</span>
              {:else if rel.IsPrerelease}
                <span class="status-badge pre">Pre-release</span>
              {:else}
                <span class="status-badge latest">Latest</span>
              {/if}
            </div>
            <time>{formatDate(rel.Created)}</time>
          </div>

          <div class="release-content">
            <header>
              <h3>{rel.Title || rel.TagName}</h3>
              <div class="commit-info">
                <span class="material-symbols-rounded">history</span>
                <code>{rel.Sha1.slice(0, 7)}</code>
              </div>
            </header>

            <div class="release-note-body">
              <MarkdownRenderer 
                 content={rel.Note} 
                 repoContext={{owner: username, name: reponame}}
              />
            </div>

            <footer class="release-assets">
              <h4>Assets</h4>
              <div class="asset-grid">
                <a
                  href="/{username}/{reponame}/archive/{rel.TagName}.zip"
                  class="asset-link"
                >
                  <span class="material-symbols-rounded">zip_box</span>
                  Source code (zip)
                </a>
                <a
                  href="/{username}/{reponame}/archive/{rel.TagName}.tar.gz"
                  class="asset-link"
                >
                  <span class="material-symbols-rounded">folder_zip</span>
                  Source code (tar.gz)
                </a>
                {#each rel.Attachments || [] as asset}
                  <a href={asset.DownloadURL} class="asset-link">
                    <span class="material-symbols-rounded">deployed_code</span>
                    {asset.Name}
                    <span class="size"
                      >({(asset.Size / 1024 / 1024).toFixed(2)} MB)</span
                    >
                  </a>
                {/each}
              </div>
            </footer>
          </div>
        </article>
      {/each}
    </div>
  {/if}
</div>

<style>
  .release-manager {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    animation: fadeIn 0.4s ease-out;
  }

  .view-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .title-group {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  .title-group span {
    font-size: 1.75rem;
    color: var(--primary-color);
  }
  .view-header h2 {
    margin: 0;
    font-size: 1.75rem;
    font-weight: 800;
  }

  .release-list {
    display: flex;
    flex-direction: column;
    gap: 2.5rem;
  }

  .release-card {
    display: grid;
    grid-template-columns: 200px 1fr;
    gap: 2rem;
    padding: 0;
    overflow: hidden;
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  .release-card.draft {
    border-style: dashed;
    opacity: 0.8;
  }

  .release-meta {
    padding: 2rem;
    background: rgba(0, 0, 0, 0.1);
    border-right: 1px solid var(--md-sys-color-outline-variant);
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .tag-status {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  .tag-badge {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-weight: 700;
    color: var(--primary-color);
    font-size: 0.95rem;
  }
  .tag-badge span {
    font-size: 1.1rem;
  }

  .status-badge {
    font-size: 0.7rem;
    font-weight: 800;
    text-transform: uppercase;
    padding: 0.2rem 0.6rem;
    border-radius: 4px;
    width: fit-content;
  }
  .status-badge.draft {
    background: rgba(255, 255, 255, 0.05);
    color: var(--text-muted);
  }
  .status-badge.pre {
    background: rgba(245, 158, 11, 0.1);
    color: #f59e0b;
  }
  .status-badge.latest {
    background: rgba(34, 197, 94, 0.1);
    color: #22c55e;
  }

  .release-meta time {
    font-size: 0.85rem;
    color: var(--text-muted);
    font-weight: 600;
  }

  .release-content {
    padding: 2.5rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
  .release-content header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .release-content h3 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 700;
  }

  .commit-info {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--text-muted);
    font-size: 0.85rem;
  }
  .commit-info code {
    background: rgba(255, 255, 255, 0.03);
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
  }

  .release-assets {
    margin-top: 1rem;
    border-top: 1px solid var(--md-sys-color-outline-variant);
    padding-top: 1.5rem;
  }
  .release-assets h4 {
    margin: 0 0 1rem 0;
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-muted);
  }

  .asset-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 0.75rem;
  }
  .asset-link {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: rgba(255, 255, 255, 0.02);
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    text-decoration: none;
    color: #fff;
    font-size: 0.9rem;
    border: 1px solid transparent;
    transition: all 0.2s;
  }
  .asset-link:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: var(--primary-color);
  }
  .asset-link span {
    font-size: 1.25rem;
    color: var(--primary-color);
  }
  .asset-link .size {
    color: var(--text-muted);
    font-size: 0.8rem;
  }

  .empty-state {
    padding: 4rem;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }
  .empty-state span {
    font-size: 4rem;
    opacity: 0.3;
    color: var(--primary-color);
  }
  .empty-state h3 {
    margin: 0;
  }
  .empty-state p {
    color: var(--text-muted);
    max-width: 400px;
  }

  .skeleton-item {
    height: 200px;
    background: var(--surface-variant);
    border-radius: 1rem;
    animation: pulse 1.5s infinite;
  }
  @keyframes pulse {
    0% {
      opacity: 0.3;
    }
    50% {
      opacity: 0.5;
    }
    100% {
      opacity: 0.3;
    }
  }
  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  @media (max-width: 800px) {
    .release-card {
      grid-template-columns: 1fr;
    }
    .release-meta {
      border-right: none;
      border-bottom: 1px solid var(--md-sys-color-outline-variant);
      flex-direction: row;
      justify-content: space-between;
      align-items: center;
    }
  }
</style>
