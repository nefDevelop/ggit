<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import Badge from '$lib/components/ui/Badge.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Avatar from '$lib/components/Avatar.svelte';
  import { getLanguageColor } from '$lib/utils/languages';
  
import Lock from "lucide-svelte/icons/lock";
import GitFork from "lucide-svelte/icons/git-fork";
import BookText from "lucide-svelte/icons/book-text";
import Star from "lucide-svelte/icons/star";
import SearchX from "lucide-svelte/icons/search-x";
import Tag from "lucide-svelte/icons/tag";

  
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
          <a href="/{repo.owner?.username || 'unknown'}/{repo.name}" class="card-link">
            <Card hover variant="low" padding="md" class="repo-card">
              <div class="repo-header">
                <Avatar 
                  src={repo.avatar_url} 
                  username={repo.name} 
                  seed={repo.id?.toString()} 
                  size={32} 
                  type="repo" 
                />
                <div class="repo-names">
                  <span class="owner">{repo.owner?.username || 'unknown'} /</span>
                  <strong class="name">{repo.name}</strong>
                </div>
                <div class="repo-type-icon">
                  {#if repo.private}
                    <Lock size={14} />
                  {:else if repo.fork}
                    <GitFork size={14} />
                  {/if}
                </div>
              </div>
              <p class="repo-desc line-clamp-2">
                {repo.description || 'No description provided.'}
              </p>
              
              <div class="repo-stats">
                <div class="stat">
                  <Star size={16} />
                  {repo.stars_count}
                </div>
                <div class="stat">
                  <GitFork size={16} />
                  {repo.forks_count}
                </div>
                {#if repo.language}
                  <div class="stat">
                    <span class="lang-dot" style="background: {getLanguageColor(repo.language)}"></span>
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
        <SearchX size={48} />
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
                    <span class="lang-dot" style="background: {lang.value ? getLanguageColor(lang.value) : 'var(--md-sys-color-outline)'}; opacity: {selectedLang === lang.value ? 1 : 0.5}"></span>
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
                    <Tag size={14} style="opacity: {selectedTopic === topic.Name ? 1 : 0.4}" />
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
  .lang-item.active { background: var(--md-sys-color-secondary-container); color: var(--md-sys-color-on-secondary-container); font-weight: 800; }

  /* MAIN */
  .repos-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem; }
  :global(.repo-card) { display: flex !important; flex-direction: column !important; gap: 1rem !important; }

  .repo-header { display: flex; align-items: center; gap: 0.75rem; }
  .repo-type-icon { margin-left: auto; color: var(--md-sys-color-primary); opacity: 0.6; display: flex; }
  .repo-names { font-size: 1.1rem; }
  .owner { opacity: 0.6; font-weight: 400; font-size: 0.9em; }
  .name { font-weight: 700; color: var(--md-sys-color-on-surface); }

  .repo-desc { font-size: 0.85rem; color: var(--md-sys-color-on-surface-variant); margin: 0; line-height: 1.4; flex: 1; opacity: 0.8; }

  .repo-stats { display: flex; align-items: center; gap: 1.25rem; font-size: 0.75rem; color: var(--md-sys-color-on-surface-variant); margin-top: auto; padding-top: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); }
  .stat { display: flex; align-items: center; gap: 0.3rem; font-weight: 700; opacity: 0.7; }
  .stat span { font-size: 1rem; color: var(--tertiary); }
  .lang-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--tertiary); }

  .repo-topics { display: flex; flex-wrap: wrap; gap: 0.4rem; margin-top: 0.5rem; }

  .empty-state { padding: 5rem; text-align: center; border-radius: 2rem; background: var(--md-sys-color-surface-container-low); }
  
  .skeleton-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1rem; }
  :global(.skeleton-card) { height: 180px; animation: pulse 1.5s infinite; background: var(--md-sys-color-surface-container-low); }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  @media (max-width: 900px) {
    .explore-repos-layout { display: flex; flex-direction: column; gap: 1.5rem; align-items: stretch; }
    .lang-sidebar { order: -1; position: static; border-radius: 1.5rem; padding: 1.25rem; }
    .sidebar-title { display: none; }
    .lang-list { flex-direction: row; flex-wrap: wrap; }
  }
</style>
