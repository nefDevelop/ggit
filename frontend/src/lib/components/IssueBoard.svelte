<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { getContext } from 'svelte';
  import Avatar from './Avatar.svelte';

  let { issues = [] } = $props();
  const repo = getContext<RepoContext>('repo');

  // Logic to group issues by columns
  let columns = $derived.by(() => {
    const cols: Record<string, any[]> = {
      'Backlog': []
    };
    
    issues.forEach(issue => {
      let status = 'Backlog';
      const statusLabel = issue.labels?.find((l: any) => l.name.startsWith('status::'));
      if (statusLabel) {
        status = statusLabel.name.split('::')[1];
      } else if (issue.state === 'closed') {
        status = 'Done';
      }
      
      if (!cols[status]) cols[status] = [];
      cols[status].push(issue);
    });

    const sortedKeys = Object.keys(cols).sort((a, b) => {
        const order = ['backlog', 'todo', 'in-progress', 'review', 'done'];
        const ia = order.indexOf(a.toLowerCase());
        const ib = order.indexOf(b.toLowerCase());
        if (ia !== -1 && ib !== -1) return ia - ib;
        if (ia !== -1) return -1;
        if (ib !== -1) return 1;
        return a.localeCompare(b);
    });

    return sortedKeys.map(key => ({
      title: key,
      issues: cols[key]
    }));
  });

  function getStatusColor(title: string) {
      const t = title.toLowerCase();
      if (t === 'done') return '#10b981';
      if (t === 'in-progress' || t === 'doing') return '#3b82f6';
      if (t === 'review') return '#f59e0b';
      if (t === 'todo') return '#6b7280';
      return 'var(--text-muted)';
  }
</script>

<div class="kanban-board">
  {#each columns as column}
    <div class="kanban-column">
      <header class="column-header">
        <div class="status-indicator" style="background: {getStatusColor(column.title)}"></div>
        <h3 class="column-title">{column.title}</h3>
        <span class="column-count">{column.issues.length}</span>
      </header>

      <div class="column-body">
        {#each column.issues as issue}
          <a href="/{repo.data.owner.username}/{repo.data.name}/issues/{issue.number}" class="kanban-card">
            <div class="card-title">{issue.title}</div>
            <div class="card-meta">
              <span class="issue-number">#{issue.number}</span>
              {#if issue.user}
                <Avatar src={issue.user.avatar_url} username={issue.user.username} size={20} />
              {/if}
            </div>
            {#if issue.labels && issue.labels.length > 0}
              <div class="card-labels">
                {#each issue.labels.filter((l: any) => !l.name.startsWith('status::')) as label}
                  <span class="mini-label" style="background: {label.color}"></span>
                {/each}
              </div>
            {/if}
          </a>
        {/each}
      </div>
    </div>
  {/each}
</div>

<style>
  .kanban-board {
    display: flex;
    gap: 1.5rem;
    overflow-x: auto;
    padding-bottom: 1rem;
    min-height: 60vh;
    align-items: flex-start;
  }

  .kanban-column {
    flex: 0 0 300px;
    background: rgba(255, 255, 255, 0.02);
    border-radius: 1rem;
    display: flex;
    flex-direction: column;
    max-height: 100%;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .column-header {
    padding: 1rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .status-indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
  }

  .column-title {
    margin: 0;
    font-size: 0.9rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    flex: 1;
  }

  .column-count {
    font-size: 0.75rem;
    font-weight: 700;
    color: var(--text-muted);
    background: var(--surface-variant);
    padding: 2px 8px;
    border-radius: 10px;
  }

  .column-body {
    padding: 0.75rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    overflow-y: auto;
  }

  .kanban-card {
    background: var(--bg-color);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 0.75rem;
    padding: 1rem;
    text-decoration: none;
    color: inherit;
    transition: all 0.2s;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    box-shadow: var(--shadow-sm);
  }

  .kanban-card:hover {
    transform: translateY(-2px);
    border-color: var(--primary-color);
    box-shadow: var(--shadow-md);
  }

  .card-title {
    font-size: 0.95rem;
    font-weight: 700;
    line-height: 1.4;
  }

  .card-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .issue-number {
    font-size: 0.8rem;
    color: var(--text-muted);
    font-weight: 600;
  }

  .card-labels {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .mini-label {
    width: 24px;
    height: 4px;
    border-radius: 2px;
    opacity: 0.8;
  }
</style>
