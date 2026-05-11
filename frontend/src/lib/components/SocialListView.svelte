<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { page } from '$app/state';
  import UserRow from '$lib/components/UserRow.svelte';
  import RepoRow from '$lib/components/RepoRow.svelte';

  let { type } = $props<{ type: 'watchers' | 'stars' | 'forks' }>();
  
  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);

  let data = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');

  const titles: Record<'watchers' | 'stars' | 'forks', string> = {
    watchers: 'Observadores',
    stars: 'Estrellas',
    forks: 'Forks'
  };

  async function fetchData() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${username}/${reponame}/${type}`);
      if (res.ok) {
        const json = await res.json();
        data = Array.isArray(json) ? json : [];
      } else {
        error = `No se pudieron cargar los ${type}.`;
        data = [];
      }
    } catch (e) {
      error = 'Error de conexión.';
    } finally {
      loading = false;
    }
  }

  onMount(fetchData);
</script>

<div class="social-list-view island">
  <header class="view-header">
    <div class="title-group">
        <span class="material-symbols-rounded">
            {type === 'watchers' ? 'visibility' : type === 'stars' ? 'star' : 'fork_left'}
        </span>
        <h2>{titles[type as 'watchers' | 'stars' | 'forks']}</h2>
    </div>
    <div class="meta">
        {data.length} {data.length === 1 ? 'elemento' : 'elementos'}
    </div>
  </header>

  <div class="content">
    {#if loading}
      <div class="skeleton-list">
        {#each Array(5) as _}
          <div class="skeleton-item"></div>
        {/each}
      </div>
    {:else if error}
      <div class="error-msg">
        <span class="material-symbols-rounded">error</span>
        <p>{error}</p>
      </div>
    {:else if data.length === 0}
      <div class="empty-state">
        <span class="material-symbols-rounded">group_off</span>
        <p>Aún no hay {type} en este repositorio.</p>
      </div>
    {:else}
      <div class="list-grid">
        {#each data as item}
          {#if type === 'forks'}
            <RepoRow repo={item} />
          {:else}
            <UserRow user={item} />
          {/if}
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .social-list-view {
    padding: 0;
    overflow: hidden;
    animation: fadeIn 0.4s ease-out;
  }

  .view-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem 2rem;
    background: var(--surface-variant);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .title-group {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .title-group span {
    color: var(--primary);
    font-size: 1.5rem;
  }

  .view-header h2 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 700;
  }

  .meta {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text-muted);
    background: var(--md-sys-color-surface-variant);
    padding: 0.2rem 0.75rem;
    border-radius: 20px;
  }

  .content {
    padding: 1.5rem;
  }

  .list-grid {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .empty-state {
    padding: 4rem;
    text-align: center;
    color: var(--text-muted);
  }

  .empty-state span {
    font-size: 3rem;
    opacity: 0.2;
    margin-bottom: 1rem;
    display: block;
  }

  .skeleton-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .skeleton-item {
    height: 60px;
    background: var(--surface-variant-hover);
    border-radius: 0.5rem;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse { 0% { opacity: 0.5; } 50% { opacity: 0.8; } 100% { opacity: 0.5; } }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
