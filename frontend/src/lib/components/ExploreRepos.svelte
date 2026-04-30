<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  
  let { q = '' } = $props();
  
  let repos = $state<any[]>([]);
  let loading = $state(true);
  let selectedLang = $state('');
  let selectedTopic = $state('');
  let popularTopics = $state<any[]>([]);
  let languages = $state<{name: string, value: string}[]>([{ name: 'All', value: '' }]);

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

  async function fetchLanguages() {
    try {
      const res = await apiFetch('/api/v2/explore/languages');
      if (res.ok) {
        const data = await res.json();
        const serverLangs = (data.languages || []).map((l: string) => ({ name: l, value: l }));
        languages = [{ name: 'All', value: '' }, ...serverLangs];
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchRepos() {
    loading = true;
    try {
      let url = `/api/v2/explore/repos?q=${encodeURIComponent(q)}&lang=${encodeURIComponent(selectedLang)}`;
      if (selectedTopic) url += `&topic=${encodeURIComponent(selectedTopic)}`;
      
      const res = await apiFetch(url);
      if (res.ok) {
        const data = await res.json();
        repos = (data.repos || []).filter((r: any) => r.owner?.username);
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
    fetchLanguages();
  });
</script>

<div class="explore-repos-layout" class:no-sidebar={languages.length <= 1 && popularTopics.length === 0}>
  <div class="explore-repos-main">
    {#if loading && repos.length === 0}
      <div class="skeleton-grid">
        {#each Array(6) as _}
          <Card variant="low" class="skeleton-card" />
        {/each}
      </div>
    {:else if repos.length > 0}
      <div class="repos-grid">
        {#each repos as repo}
          <a href="/{repo.owner.username}/{repo.name}" class="card-link">
            <Card hover variant="low" padding="md" class="repo-card">
              <div class="repo-header">
                <span class="material-symbols-rounded repo-icon">
                  {repo.private ? 'lock' : repo.fork ? 'fork_right' : 'book'}
                </span>
                <div class="repo-names">
                  <span class="owner">{repo.owner.username} /</span>
                  <strong class="name">{repo.name}</strong>
                </div>
              </div>
              <p class="repo-desc line-clamp-2">
                {repo.description || 'No description provided.'}
              </p>
              
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
              </div>

              {#if repo.topics && repo.topics.length > 0}
                <div class="repo-topics">
                  {#each repo.topics.slice(0, 3) as topic}
                    <Badge variant="ghost" size="sm">#{topic}</Badge>
                  {/each}
                </div>
              {/if}
            </Card>
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

  {#if languages.length > 1 || popularTopics.length > 0}
    <aside class="lang-sidebar tonal">
      {#if languages.length > 1}
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
      {/if}

      {#if popularTopics.length > 0}
        <h3 class="sidebar-title" style="margin-top: {languages.length > 1 ? '2.5rem' : '1rem'};">Topics</h3>
        <div class="lang-list">
            {#each popularTopics.slice(0, 10) as topic}
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
  {/if}
</div>

<style>
  .explore-repos-layout { display: grid; grid-template-columns: 1fr 220px; gap: 3rem; align-items: start; }
  .explore-repos-layout.no-sidebar { grid-template-columns: 1fr; }
  .card-link { text-decoration: none; color: inherit; display: block; }

  /* SIDEBAR */
  .lang-sidebar { padding: 1.5rem 0.75rem; position: sticky; top: 2rem; background: var(--md-sys-color-surface-container-low); border-radius: 2rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .sidebar-title { margin: 0 0 1rem 0.75rem; font-size: 0.7rem; font-weight: 900; text-transform: uppercase; letter-spacing: 0.15em; color: var(--md-sys-color-on-surface-variant); opacity: 0.5; }
  .lang-list { display: flex; flex-direction: column; gap: 0.2rem; }
  .lang-item { background: transparent; border: none; padding: 0.6rem 0.75rem; border-radius: 0.8rem; text-align: left; color: var(--md-sys-color-on-surface-variant); font-weight: 600; font-size: 0.85rem; cursor: pointer; transition: all 0.2s; display: flex; align-items: center; gap: 0.75rem; }
  .lang-item:hover { background: var(--md-sys-color-surface-variant); }
  .lang-item.active { background: var(--md-sys-color-primary-container); color: var(--md-sys-color-on-primary-container); font-weight: 800; }

  /* MAIN */
  .repos-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem; }
  :global(.repo-card) { display: flex !important; flex-direction: column !important; gap: 1rem !important; }

  .repo-header { display: flex; align-items: center; gap: 0.75rem; }
  .repo-icon { font-size: 1.5rem; color: var(--md-sys-color-primary); opacity: 0.8; }
  .repo-names { font-size: 1.1rem; }
  .owner { opacity: 0.6; font-weight: 400; font-size: 0.9em; }
  .name { font-weight: 700; color: var(--md-sys-color-on-surface); }

  .repo-desc { font-size: 0.85rem; color: var(--md-sys-color-on-surface-variant); margin: 0; line-height: 1.4; flex: 1; opacity: 0.8; }

  .repo-stats { display: flex; align-items: center; gap: 1.25rem; font-size: 0.75rem; color: var(--md-sys-color-on-surface-variant); margin-top: auto; padding-top: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); }
  .stat { display: flex; align-items: center; gap: 0.3rem; font-weight: 700; opacity: 0.7; }
  .stat span { font-size: 1rem; color: var(--md-sys-color-primary); }
  .lang-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--md-sys-color-primary); }

  .repo-topics { display: flex; flex-wrap: wrap; gap: 0.4rem; margin-top: 0.5rem; }

  .empty-state { padding: 5rem; text-align: center; border-radius: 2rem; background: var(--md-sys-color-surface-container-low); }
  .empty-state span { font-size: 4rem; color: var(--md-sys-color-primary); opacity: 0.3; }
  
  .skeleton-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1rem; }
  :global(.skeleton-card) { height: 180px; animation: pulse 1.5s infinite; background: var(--md-sys-color-surface-container-low); }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  @media (max-width: 900px) {
    .explore-repos-layout { grid-template-columns: 1fr; gap: 1.5rem; }
    .lang-sidebar { position: static; border-radius: 1.5rem; padding: 1.25rem; }
    .sidebar-title { display: none; }
    .lang-list { flex-direction: row; flex-wrap: wrap; }
  }
</style>
