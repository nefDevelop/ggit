<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Button from './ui/Button.svelte';

  let { q = '' } = $props();
  
  let topics = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');

  async function fetchTopics() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/explore/topics?q=${q}`);
      if (res.ok) {
        topics = await res.json();
      }
    } catch (e) {
      error = 'Failed to load topics.';
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    fetchTopics();
  });
</script>

<div class="explore-topics">
  {#if loading}
    <div class="loading-grid">
      {#each Array(8) as _}
        <div class="skeleton-topic tonal"></div>
      {/each}
    </div>
  {:else if topics.length === 0}
    <div class="empty-state tonal">
      <span class="material-symbols-rounded">tag_off</span>
      <p>No topics found matching your search.</p>
    </div>
  {:else}
    <div class="topics-grid">
      {#each topics as topic}
        <a href="/explore/repos?q=topic:{topic.Name}" class="topic-card tonal hover-lift">
          <div class="topic-icon">
            <span class="material-symbols-rounded">tag</span>
          </div>
          <div class="topic-info">
            <strong>{topic.Name}</strong>
            <span>{topic.RepoCount} Repositorios</span>
          </div>
        </a>
      {/each}
    </div>
  {/if}
</div>

<style>
  .topics-grid { 
    display: grid; 
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); 
    gap: 1.5rem; 
    animation: fadeIn 0.4s ease-out;
  }

  .topic-card { 
    display: flex; 
    align-items: center; 
    gap: 1.25rem; 
    padding: 1.5rem; 
    text-decoration: none; 
    color: inherit; 
    border-radius: 1.5rem;
    background: var(--surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }
  .topic-card:hover { border-color: var(--primary-color); transform: translateY(-4px); box-shadow: var(--shadow-lg); }

  .topic-icon { 
    width: 48px; 
    height: 48px; 
    display: flex; 
    align-items: center; 
    justify-content: center; 
    background: rgba(var(--primary-rgb), 0.1); 
    color: var(--primary-color); 
    border-radius: 12px;
  }
  .topic-icon span { font-size: 1.5rem; }

  .topic-info { display: flex; flex-direction: column; }
  .topic-info strong { font-size: 1.1rem; font-weight: 800; letter-spacing: -0.01em; }
  .topic-info span { font-size: 0.8rem; color: var(--text-muted); font-weight: 600; }

  .loading-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1.5rem; }
  .skeleton-topic { height: 80px; border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }

  .empty-state { padding: 4rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state span { font-size: 3rem; opacity: 0.2; }

  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
