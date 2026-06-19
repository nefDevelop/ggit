<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { dialog } from "$lib/dialog.svelte";
  import Avatar from "$lib/components/Avatar.svelte";
  import ReleaseEditor from "./ReleaseEditor.svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import Badge from "./ui/Badge.svelte";
  
import { Tag, Rocket, Plus, Edit, Trash2, ChevronsLeft, ChevronRight, Calendar, GitMerge, Download, FileArchive } from "lucide-svelte";


  let { username, reponame } = $props();

  let releases = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  let hasNext = $state(false);
  let nextAfter = $state("");
  let hasPrevious = $state(false);

  let showEditor = $state(false);
  let editingRelease = $state<any>(null);

  async function fetchReleases(after = "") {
    loading = true;
    error = "";
    try {
      const query = after ? `?after=${after}` : "";
      const res = await apiFetch(
        `/api/v2/${username}/${reponame}/releases${query}`,
      );
      if (res.ok) {
        const data = await res.json();
        if (Array.isArray(data)) {
          releases = data;
          hasNext = false;
        } else {
          releases = data.releases || [];
          hasNext = data.has_next;
          nextAfter = data.next_after;
          hasPrevious = data.has_previous;
        }
      } else {
        error = "Failed to load releases.";
      }
    } catch (e) {
      console.error(e);
      error = "Network error. Please try again.";
    } finally {
      loading = false;
    }
  }

  function handleAdd() {
    editingRelease = null;
    showEditor = true;
  }

  function handleEdit(rel: any) {
    editingRelease = rel;
    showEditor = true;
  }

  async function handleDelete(id: number) {
    if (
      !(await dialog.confirm(
        "Are you sure you want to delete this release? This will also remove the associated Git tag if it was created through Ggit.",
      ))
    )
      return;

    try {
      const res = await apiFetch(
        `/api/v2/${username}/${reponame}/releases/${id}`,
        {
          method: "DELETE",
        },
      );
      if (res.ok) {
        releases = releases.filter((r) => r.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  function onEditorSaved(newRel: any) {
    showEditor = false;
    fetchReleases();
  }

  onMount(() => fetchReleases());
</script>

<div class="repo-releases">
  {#if showEditor}
    <div class="editor-overlay">
      <ReleaseEditor
        {username}
        {reponame}
        release={editingRelease}
        onSaved={onEditorSaved}
        onCancel={() => (showEditor = false)}
      />
    </div>
  {/if}

  {#if releases.length > 0}
  <header class="section-actions">
    <Button variant="primary" onclick={handleAdd} icon={Plus}>
      New Release
    </Button>
  </header>
  {/if}

  <!-- Header actions removed, moved to timeline start -->

  {#if loading && releases.length === 0}
    <div class="skeleton-list">
      {#each Array(3) as _}
        <div class="skeleton-item"></div>
      {/each}
    </div>
  {:else if error}
    <div class="error-msg">{error}</div>
  {:else if releases.length === 0}
    <div class="empty-state">
      <div class="empty-icon-ring">
        <Rocket size={40} />
      </div>
      <h2>No releases yet</h2>
      <p>Start tagging your milestones by creating the first release.</p>
      <Button onclick={handleAdd} icon={Plus}>Create your first release</Button>
    </div>
  {:else}
    <div class="flat-list-container">
      {#each releases as release}
        <div class="list-item-row" class:is-draft={release.is_draft}>
          <div class="list-item-icon">
            <Rocket size={18} style="color: var(--text-muted); opacity: 0.5;" />
          </div>

          <div class="item-main-content">
            <div class="item-title-bar">
              <span class="item-title-link">
                {release.title || release.tag_name}
              </span>
              {#if release.is_draft}
                <span class="badge-mini" style="background: var(--text-muted)">Draft</span>
              {:else if release.is_prerelease}
                <span class="badge-mini" style="background: #f59e0b">Pre-release</span>
              {/if}
              {#if release === releases.find((r) => !r.is_draft && !r.is_prerelease)}
                <span class="badge-mini" style="background: var(--success)">Latest</span>
              {/if}
            </div>

            <div class="item-meta-data">
              {#if release.tag_name}
                <span style="display: inline-flex; align-items: center; gap: 0.2rem; margin-right: 0.75rem;">
                  <Tag size={12} />
                  <a href={`/${username}/${reponame}/src/tag/${release.tag_name}`} style="color: inherit; text-decoration: none;">{release.tag_name}</a>
                </span>
              {/if}
              {#if release.target}
                <span style="display: inline-flex; align-items: center; gap: 0.2rem; margin-right: 0.75rem;">
                  <GitMerge size={12} /> {release.target}
                </span>
              {/if}
              <span style="display: inline-flex; align-items: center; gap: 0.2rem;">
                <Calendar size={12} />
                {new Date(release.created_at * 1000).toLocaleDateString(undefined, { dateStyle: "medium" })}
              </span>
            </div>

            {#if release.note}
              <div class="release-note-preview" style="font-size: 0.9rem; color: var(--text-color); margin-top: 0.5rem; opacity: 0.9; max-height: 4.5em; overflow: hidden; text-overflow: ellipsis; display: -webkit-box; -webkit-line-clamp: 3; -webkit-box-orient: vertical;">
                {release.note}
              </div>
            {/if}

            <div class="asset-links-inline" style="display: flex; align-items: center; gap: 0.5rem; margin-top: 0.75rem;">
              {#if release.attachments && release.attachments.length > 0}
                {#each release.attachments as asset}
                  <a href={asset.download_url} class="btn-icon-minimal" title={asset.name} download>
                    <Download size={16} />
                  </a>
                {/each}
              {/if}
              <a href={`/api/v2/${username}/${reponame}/archive/${release.tag_name}.zip`} class="btn-icon-minimal" title="Source code (zip)" download>
                <FileArchive size={16} />
              </a>
              <a href={`/api/v2/${username}/${reponame}/archive/${release.tag_name}.tar.gz`} class="btn-icon-minimal" title="Source code (tar.gz)" download>
                <FileArchive size={16} />
              </a>
            </div>
          </div>

          {#if release.id}
            <div class="item-actions-right">
              <button class="btn-icon-minimal" onclick={() => handleEdit(release)} title="Edit"><Edit size={16} /></button>
              <button class="btn-icon-minimal delete" onclick={() => handleDelete(release.id)} title="Delete"><Trash2 size={16} /></button>
            </div>
          {/if}
        </div>
      {/each}
    </div>

    {#if hasNext || hasPrevious}
      <div class="pagination">
        <Button
          variant="secondary"
          onclick={() => fetchReleases("")}
          disabled={loading || !hasPrevious}
          icon={ChevronsLeft}
        >
          Newest
        </Button>
        <Button
          onclick={() => fetchReleases(nextAfter)}
          disabled={loading || !hasNext}
          icon={ChevronRight}
        >
          Older releases
        </Button>
      </div>
    {/if}
  {/if}
</div>

<style>
  .repo-releases {
    display: flex;
    flex-direction: column;
    gap: 3rem;
    min-height: 400px;
    padding: 0;
  }

  .section-actions {
    display: flex;
    justify-content: flex-end;
    padding: 0.5rem 0;
  }

  .editor-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(12px);
    z-index: 1000;
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding: 2rem;
    overflow-y: auto;
  }

  @media (max-width: 1100px) {
    .editor-overlay {
      padding: 1rem;
    }
  }

  .empty-state {
    padding: 6rem 2rem;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.25rem;
    color: var(--text-muted);
  }

  .empty-state h2 {
    margin: 0;
    font-size: 1.75rem;
    font-weight: 950;
    color: var(--text-color);
  }

  .empty-state p {
    max-width: 400px;
    margin: 0 0 1rem 0;
    font-size: 1.1rem;
    line-height: 1.5;
  }

  .empty-icon-ring {
    width: 80px;
    height: 80px;
    border: 2px dashed var(--md-sys-color-outline-variant);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 1rem;
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }
</style>
