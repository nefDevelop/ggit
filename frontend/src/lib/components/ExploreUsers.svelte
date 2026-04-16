<script lang="ts">
  import { onMount } from 'svelte';
  
  let { q = '' } = $props();
  
  let users = $state<any[]>([]);
  let loading = $state(true);
  let total = $state(0);
  let page = $state(1);

  async function fetchUsers() {
    loading = true;
    try {
      const res = await fetch(`/api/v2/explore/users?q=${encodeURIComponent(q)}&page=${page}`);
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
        <div class="skeleton-row island"></div>
      {/each}
    </div>
  {:else if users.length > 0}
    <div class="users-list">
      {#each users as user}
        <a href="/{user.username}" class="user-card island tonal hover-lift">
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
  .users-list { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1rem; }
  .user-card { display: flex; align-items: center; gap: 1.25rem; padding: 1.25rem; text-decoration: none; color: inherit; border: 1px solid var(--border-color); }
  
  .avatar { width: 48px; height: 48px; border-radius: 50%; background: #222; }
  
  .user-info { flex: 1; display: flex; flex-direction: column; gap: 0.1rem; }
  .username { font-size: 1.1rem; font-weight: 700; color: var(--primary-color); }
  .full-name { font-size: 0.85rem; color: var(--text-muted); }

  .user-meta { display: flex; align-items: center; gap: 0.4rem; font-size: 0.75rem; color: var(--text-muted); opacity: 0.7; }
  .user-meta span { font-size: 0.9rem; }

  .empty-state { text-align: center; padding: 4rem; color: var(--text-muted); }
  .empty-state span { font-size: 4rem; opacity: 0.2; }
  
  .skeleton-list { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1rem; }
  .skeleton-row { height: 80px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
