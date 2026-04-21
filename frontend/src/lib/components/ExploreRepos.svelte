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
  <!-- MAIN LIST -->
  <div class="explore-repos-main">
    {#if loading && repos.length === 0}
      <div class="skeleton-grid">
        {#each Array(6) as _}
          <div class="skeleton-card"></div>
        {/each}
      </div>
    {:else if repos.length > 0}
      <div class="repos-grid">
        {#each repos as repo}
          <a href="/{repo.owner?.username || 'unknown'}/{repo.name}" class="repo-card tonal hover-lift">
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
      <div class="empty-state tonal">
        <span class="material-symbols-rounded">search_off</span>
        <h3>No repositories found</h3>
        <p>Try searching for something else or changing filters.</p>
      </div>
    {/if}
  </div>

  <!-- SIDEBAR (RIGHT) -->
  <aside class="lang-sidebar tonal">
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
</div>

<style>
  .explore-repos-layout { display: grid; grid-template-columns: 1fr 220px; gap: 3rem; align-items: start; }

  /* SIDEBAR */
  .lang-sidebar { padding: 0.5rem; position: sticky; top: 2rem; }
  .sidebar-title { 
    margin: 0 0 1rem 0.75rem; 
    font-size: 0.7rem; 
    font-weight: 900; 
    text-transform: uppercase; 
    letter-spacing: 0.15em; 
    color: var(--text-muted); 
    opacity: 0.5;
  }
  .lang-list { display: flex; flex-direction: column; gap: 0.1rem; }
  .lang-item { 
    background: transparent; 
    border: none; 
    padding: 0.6rem 0.75rem; 
    border-radius: 0.5rem; 
    text-align: left; 
    color: var(--text-muted); 
    font-weight: 600; 
    font-size: 0.85rem; 
    cursor: pointer; 
    transition: all 0.2s; 
    display: flex; 
    align-items: center; 
    gap: 0.75rem; 
  }
  .lang-item:hover { background: var(--md-sys-color-surface-container-low); color: var(--text-color); }
  .lang-item.active { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); font-weight: 800; }

  /* MAIN */
  .repos-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 2.5rem; grid-auto-rows: min-content; }
  .repo-card { 
    display: flex; 
    flex-direction: column; 
    gap: 1rem; 
    padding: 1.75rem; 
    text-decoration: none; 
    color: inherit; 
    border: 1px solid var(--md-sys-color-outline-variant); 
    border-radius: 2rem;
    transition: all 0.5s cubic-bezier(0.16, 1, 0.3, 1);
    background: var(--md-sys-color-surface-container-low);
    position: relative;
  }
  .repo-card:hover { 
    border-color: var(--primary-color); 
    background: var(--md-sys-color-surface-container);
    transform: translateY(-4px);
    box-shadow: 0 12px 24px -10px rgba(0, 0, 0, 0.2);
  }

  .repo-header { display: flex; align-items: center; gap: 0.75rem; }
  .repo-icon { font-size: 1.5rem; color: var(--primary-color); opacity: 0.8; }
  .repo-names { font-size: 1.1rem; font-weight: 500; }
  .owner { opacity: 0.6; font-weight: 400; font-size: 0.9em; }
  .name { font-weight: 700; color: var(--md-sys-color-on-surface); }

  .repo-desc { 
    font-size: 0.85rem; 
    color: var(--text-muted); 
    margin: 0; 
    line-height: 1.5; 
    flex: 1; 
    display: -webkit-box; 
    -webkit-line-clamp: 2; 
    -webkit-box-orient: vertical; 
    overflow: hidden; 
    opacity: 0.8;
  }

  .repo-stats { 
    display: flex; 
    align-items: center; 
    gap: 1.25rem; 
    font-size: 0.75rem; 
    color: var(--text-muted); 
    margin-top: 0.75rem; 
    padding-top: 1rem; 
    border-top: 1px solid var(--md-sys-color-outline-variant); 
  }
  .stat { display: flex; align-items: center; gap: 0.3rem; font-weight: 700; opacity: 0.7; }
  .stat span { font-size: 1rem; color: var(--primary-color); }
  .lang-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--primary-color); }
  .updated { margin-left: auto; opacity: 0.4; font-size: 0.7rem; font-weight: 600; }

  .repo-topics { display: flex; flex-wrap: wrap; gap: 0.4rem; margin-top: 0.5rem; }
  .topic-badge { 
    background: var(--md-sys-color-surface-container-high);
    color: var(--text-muted); 
    padding: 0.15rem 0.5rem; 
    border-radius: 0.5rem; 
    font-size: 0.65rem; 
    font-weight: 700; 
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.2s;
  }
  .repo-card:hover .topic-badge { border-color: rgba(var(--primary-rgb), 0.3); color: var(--primary-color); }

  .empty-state { padding: 5rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state span { font-size: 4rem; color: var(--primary-color); opacity: 0.3; }
  
  .skeleton-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(360px, 1fr)); gap: 1rem; }
  .skeleton-card { height: 160px; border-radius: 1rem; animation: pulse 1.5s infinite; background: var(--md-sys-color-surface-container); }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.5; } 100% { opacity: 0.3; } }

  @media (max-width: 900px) {
    .explore-repos-layout { grid-template-columns: 1fr; gap: 1.5rem; }
    .lang-sidebar { position: static; }
    .lang-list { flex-direction: row; flex-wrap: wrap; }
  }
</style>
