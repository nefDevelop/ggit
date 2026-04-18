<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let { q = '' } = $props();
  
  let repos = $state<any[]>([]);
  let loading = $state(true);
  let selectedLang = $state('');
  let selectedTopic = $state('');
  let popularTopics = $state<any[]>([]);

  const languages = [
    { name: 'All', value: '' },
    { name: 'Go', value: 'Go' },
    { name: 'JavaScript', value: 'JavaScript' },
    { name: 'TypeScript', value: 'TypeScript' },
    { name: 'Python', value: 'Python' },
    { name: 'Rust', value: 'Rust' },
    { name: 'C++', value: 'C++' },
    { name: 'Java', value: 'Java' },
    { name: 'HTML', value: 'HTML' },
    { name: 'CSS', value: 'CSS' },
    { name: 'Zig', value: 'Zig' },
  ];

  async function fetchPopularTopics() {
    try {
      const res = await apiFetch('/api/v2/explore/topics');
      if (res.ok) {
        const data = await res.json();
        popularTopics = data.topics || [];
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchRepos() {
    loading = true;
    try {
      let url = `/api/v2/explore/repos?q=${encodeURIComponent(q)}&lang=${encodeURIComponent(selectedLang)}`;
      if (selectedTopic) {
        url += `&topic=${encodeURIComponent(selectedTopic)}`;
      }
      const res = await apiFetch(url);
      if (res.ok) {
        const data = await res.json();
        repos = data.repos || [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    fetchRepos();
  });

  onMount(() => {
    fetchRepos();
    fetchPopularTopics();
  });
</script>

<div class="explore-repos-layout">
  <!-- SIDEBAR -->
  <aside class="lang-sidebar island tonal">
    <h3 class="sidebar-title">Languages</h3>
    <div class="lang-list">
        {#each languages as lang}
            <button 
                class="lang-item" 
                class:active={selectedLang === lang.value}
                onclick={() => selectedLang = lang.value}
            >
                <span class="lang-dot" style="opacity: {selectedLang === lang.value ? 1 : 0.4}"></span>
                {lang.name}
            </button>
        {/each}
    </div>

    {#if popularTopics.length > 0}
      <h3 class="sidebar-title" style="margin-top: 2.5rem;">Topics</h3>
      <div class="lang-list">
          <button 
              class="lang-item" 
              class:active={selectedTopic === ''}
              onclick={() => selectedTopic = ''}
          >
              <span class="material-symbols-rounded" style="font-size: 1.1rem; opacity: 0.5;">tag</span>
              Clear Tópicos
          </button>
          {#each popularTopics as topic}
              <button 
                  class="lang-item" 
                  class:active={selectedTopic === topic.Name}
                  onclick={() => selectedTopic = topic.Name}
              >
                  <span class="material-symbols-rounded" style="font-size: 1.1rem; opacity: {selectedTopic === topic.Name ? 1 : 0.4}">tag</span>
                  {topic.Name}
              </button>
          {/each}
      </div>
    {/if}
  </aside>

  <!-- MAIN LIST -->
  <div class="explore-repos-main">
    {#if loading && repos.length === 0}
      <div class="skeleton-grid">
        {#each Array(6) as _}
          <div class="skeleton-card island"></div>
        {/each}
      </div>
    {:else if repos.length > 0}
      <div class="repos-grid">
        {#each repos as repo}
          <a href="/{repo.owner?.username || 'unknown'}/{repo.name}" class="repo-card island tonal hover-lift">
            <div class="repo-header">
              <span class="material-symbols-rounded repo-icon">
                {repo.private ? 'lock' : repo.fork ? 'fork_right' : 'book'}
              </span>
              <div class="repo-names">
                <span class="owner">{repo.owner?.username || 'unknown'} /</span>
                <strong class="name">{repo.name}</strong>
              </div>
            </div>
            <p class="repo-desc">{repo.description || 'No description provided.'}</p>
            <div class="repo-stats">
              <div class="stat">
                <span class="material-symbols-rounded">star</span>
                {repo.stars_count}
              </div>
              <div class="stat">
                <span class="material-symbols-rounded">fork_right</span>
                {repo.forks_count}
              </div>
              {#if repo.language}
                <div class="stat">
                  <span class="lang-dot"></span>
                  {repo.language}
                </div>
              {/if}
              <span class="updated">Updated {new Date(repo.updated_at).toLocaleDateString()}</span>
            </div>
            {#if repo.topics && repo.topics.length > 0}
              <div class="repo-topics">
                {#each repo.topics as topic}
                  <span class="topic-badge">{topic}</span>
                {/each}
              </div>
            {/if}
          </a>
        {/each}
      </div>
    {:else}
      <div class="empty-state island tonal">
        <span class="material-symbols-rounded">search_off</span>
        <h3>No repositories found</h3>
        <p>Try searching for something else or changing filters.</p>
      </div>
    {/if}
  </div>
</div>

<style>
  .explore-repos-layout { display: grid; grid-template-columns: 240px 1fr; gap: 2rem; align-items: start; }

  /* SIDEBAR */
  .lang-sidebar { padding: 1.5rem; position: sticky; top: 2rem; }
  .sidebar-title { margin: 0 0 1.5rem 0; font-size: 1.1rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-muted); }
  .lang-list { display: flex; flex-direction: column; gap: 0.25rem; }
  .lang-item { background: transparent; border: none; padding: 0.75rem 1rem; border-radius: 0.75rem; text-align: left; color: var(--text-muted); font-weight: 700; font-size: 0.95rem; cursor: pointer; transition: all 0.2s; display: flex; align-items: center; gap: 0.75rem; }
  .lang-item:hover { background: rgba(255, 255, 255, 0.03); color: var(--text-color); }
  .lang-item.active { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); }

  /* MAIN */
  .repos-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(400px, 1fr)); gap: 1.5rem; }
  .repo-card { display: flex; flex-direction: column; gap: 1rem; padding: 2rem; text-decoration: none; color: inherit; height: 100%; border: 1px solid var(--border-color); transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
  .repo-card:hover { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .repo-header { display: flex; align-items: center; gap: 1rem; }
  .repo-icon { font-size: 1.75rem; color: var(--primary-color); opacity: 0.8; }
  .repo-names { font-size: 1.25rem; font-weight: 500; }
  .owner { opacity: 0.6; font-weight: 400; }
  .name { font-weight: 800; color: #fff; }

  .repo-desc { font-size: 0.95rem; color: var(--text-muted); margin: 0; line-height: 1.6; flex: 1; display: -webkit-box; -webkit-line-clamp: 3; line-clamp: 3; -webkit-box-orient: vertical; overflow: hidden; }

  .repo-stats { display: flex; align-items: center; gap: 1.5rem; font-size: 0.85rem; color: var(--text-muted); margin-top: 1.5rem; padding-top: 1rem; border-top: 1px solid var(--border-color); }
  .stat { display: flex; align-items: center; gap: 0.4rem; font-weight: 700; }
  .stat span { font-size: 1.1rem; color: var(--primary-color); }
  .lang-dot { width: 10px; height: 10px; border-radius: 50%; background: var(--primary-color); box-shadow: 0 0 8px var(--primary-color); }
  .updated { margin-left: auto; opacity: 0.4; font-size: 0.75rem; }

  .empty-state { padding: 5rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state span { font-size: 4rem; color: var(--primary-color); opacity: 0.3; }
  .empty-state h3 { margin: 0; }
  .empty-state p { color: var(--text-muted); }

  .repo-topics { display: flex; flex-wrap: wrap; gap: 0.5rem; margin-top: 1rem; }
  .topic-badge { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.2rem 0.6rem; border-radius: 2rem; font-size: 0.75rem; font-weight: 800; border: 1px solid rgba(var(--primary-rgb), 0.2); }
  .topic-badge:hover { background: rgba(var(--primary-rgb), 0.2); }

  .skeleton-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(400px, 1fr)); gap: 1.5rem; }
  .skeleton-card { height: 180px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.5; } 100% { opacity: 0.3; } }

  @media (max-width: 900px) {
    .explore-repos-layout { grid-template-columns: 1fr; }
    .lang-sidebar { position: static; overflow-x: auto; }
    .lang-list { flex-direction: row; }
  }
</style>
