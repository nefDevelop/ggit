<script lang="ts">
  import GitPullRequest from "lucide-svelte/icons/git-pull-request";
  import AlertCircle from "lucide-svelte/icons/alert-circle";
  import CheckCircle2 from "lucide-svelte/icons/check-circle-2";

  let { inboxItems = [], inboxLoading = false }: {
    inboxItems: any[];
    inboxLoading: boolean;
  } = $props();
</script>

<div class="inbox-stream">
  {#if inboxLoading && inboxItems.length === 0}
    <div class="skeleton-list">
      {#each Array(3) as _}<div class="skeleton-row"></div>{/each}
    </div>
  {:else if inboxItems.length === 0}
    <div class="empty-dashboard">
      <CheckCircle2 size={48} />
      <h3>You're all caught up!</h3>
      <p>No issues or pull requests assigned to you.</p>
    </div>
  {:else}
    {#each inboxItems as item}
      <a href={`/${item.repo_owner}/${item.repo_name}/${item.is_pull ? 'pulls' : 'issues'}/${item.index}`} class="inbox-card tonal">
        <div class="inbox-icon {item.is_pull ? 'pull' : 'issue'}">
          {#if item.is_pull}
            <GitPullRequest size={20} />
          {:else}
            <AlertCircle size={20} />
          {/if}
        </div>
        <div class="inbox-content">
          <div class="inbox-repo-line">{item.repo_owner} / {item.repo_name}</div>
          <div class="inbox-title-line">
            <span class="index">#{item.index}</span>
            <span class="title">{item.title}</span>
          </div>
          <div class="inbox-meta-line">
            <span>Opened by {item.poster_name} • {new Date(item.updated_at * 1000).toLocaleDateString()}</span>
          </div>
        </div>
      </a>
    {/each}
  {/if}
</div>

<style>
  .inbox-stream { display: flex; flex-direction: column; gap: 0.5rem; }
  .empty-dashboard {
    display: flex; flex-direction: column; align-items: center;
    justify-content: center; padding: 4rem 2rem; gap: 1rem;
    opacity: 0.7; text-align: center;
  }
  .inbox-card {
    display: flex; align-items: flex-start; gap: 1rem;
    padding: 1rem; border-radius: 1rem;
    text-decoration: none; color: inherit;
    transition: background 0.2s;
  }
  .inbox-card:hover { background: var(--md-sys-color-surface-container-high); }
  .inbox-icon {
    width: 40px; height: 40px; border-radius: 0.75rem;
    display: flex; align-items: center; justify-content: center; flex-shrink: 0;
  }
  .inbox-icon.pull { background: color-mix(in srgb, var(--status-success) 10%, transparent); color: var(--status-success); }
  .inbox-icon.issue { background: rgba(var(--error-rgb), 0.1); color: var(--error); }
  .inbox-content { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 0.25rem; }
  .inbox-repo-line { font-size: 0.75rem; font-weight: 700; color: var(--text-muted); }
  .inbox-title-line { display: flex; gap: 0.5rem; align-items: baseline; }
  .index { font-weight: 800; color: var(--primary-color); flex-shrink: 0; }
  .title { font-weight: 600; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .inbox-meta-line { font-size: 0.75rem; color: var(--text-muted); }
  .skeleton-list { display: flex; flex-direction: column; gap: 0.5rem; }
  .skeleton-row { height: 80px; background: var(--md-sys-color-surface-container-high); border-radius: 1rem; }
  h3 { margin: 0; font-size: 1.25rem; }
  p { margin: 0; font-size: 0.9rem; }
</style>
