<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let { q = '' } = $props();
  
  let users = $state<any[]>([]);
  let loading = $state(true);
  let total = $state(0);
  let page = $state(1);

  async function fetchUsers() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/explore/users?q=${encodeURIComponent(q)}&page=${page}`);
      if (res.ok) {
        const data = await res.json();
        users = data.users || [];
        total = data.total || 0;
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  // React to search query changes with debounce
  let timeout: any;
  $effect(() => {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      page = 1;
      fetchUsers();
    }, 300);
  });

  onMount(fetchUsers);
</script>

<div class="explore-users">
  {#if loading && users.length === 0}
    <div class="skeleton-list">
      {#each Array(6) as _}
        <div class="skeleton-row"></div>
      {/each}
    </div>
  {:else if users.length > 0}
    <div class="users-list">
      {#each users as user}
        <a href="/{user.username}" class="user-card tonal hover-lift">
          <img src={user.avatar_url} alt={user.username} class="avatar" />
          <div class="user-info">
            <strong class="username">{user.username}</strong>
            <span class="full-name">{user.full_name || 'No display name'}</span>
          </div>
          {#if user.location}
            <div class="user-meta">
              <span class="material-symbols-rounded">location_on</span>
              {user.location}
            </div>
          {/if}
          <div class="user-meta">
            <span class="material-symbols-rounded">calendar_today</span>
            Joined {new Date(user.created_at).getFullYear()}
          </div>
        </a>
      {/each}
    </div>
  {:else}
    <div class="empty-state">
      <span class="material-symbols-rounded">person_off</span>
      <h3>No users found</h3>
      <p>Try searching for a different username.</p>
    </div>
  {/if}
</div>

<style>
  .explore-users { padding: 1rem 0; }
  .users-list { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1.5rem; grid-auto-rows: min-content; }
  .user-card { 
    display: flex; 
    align-items: center; 
    gap: 1rem; 
    padding: 1.25rem; 
    text-decoration: none; 
    color: inherit; 
    border: 1px solid var(--md-sys-color-outline-variant); 
    border-radius: 2rem;
    background: var(--md-sys-color-surface-container-low);
    transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1);
    position: relative;
  }
  .user-card:hover { 
    transform: translateX(6px);
    border-color: var(--primary-color);
    background: var(--md-sys-color-surface-container);
    box-shadow: 0 8px 16px -8px rgba(0, 0, 0, 0.2);
  }
  
  .avatar { width: 44px; height: 44px; border-radius: 50%; border: 2px solid var(--md-sys-color-outline-variant); transition: all 0.3s; }
  .user-card:hover .avatar { border-color: var(--primary-color); transform: rotate(5deg); }
  
  .user-info { flex: 1; display: flex; flex-direction: column; gap: 0.05rem; }
  .username { font-size: 1rem; font-weight: 800; color: var(--md-sys-color-on-surface); }
  .full-name { font-size: 0.8rem; color: var(--text-muted); font-weight: 500; }

  .user-meta { display: flex; align-items: center; gap: 0.3rem; font-size: 0.7rem; color: var(--text-muted); opacity: 0.6; font-weight: 600; }
  .user-meta span { font-size: 0.85rem; color: var(--primary-color); }

  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); }
  .empty-state span { font-size: 4rem; opacity: 0.2; }
  
  .skeleton-list { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1rem; }
  .skeleton-row { height: 80px; border-radius: 1rem; animation: pulse 1.5s infinite; background: var(--md-sys-color-surface-container); }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.5; } }
</style>
