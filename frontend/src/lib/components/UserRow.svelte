<script lang="ts">
  import Avatar from '$lib/components/Avatar.svelte';
  let { user } = $props();
  
  // Format dates or other info if needed
</script>

<div class="user-row island">
  <Avatar src={user.avatar_url} username={user.username} size={48} class="avatar" />
  <div class="user-info">
    <a href="/{user.username}" class="name">{user.full_name || user.username}</a>
    <span class="username">@{user.username}</span>
  </div>
  <div class="actions">
    <a href="/{user.username}" class="btn btn-secondary btn-sm">Ver Perfil</a>
  </div>
</div>

<style>
  .user-row {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 1rem 1.5rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    border: 1px solid transparent;
  }

  .user-row:hover {
    background: var(--surface-variant-hover);
    border-color: var(--primary-light);
    transform: translateX(4px);
  }

  .user-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
  }

  .name {
    font-weight: 700;
    color: var(--text-color);
    text-decoration: none;
    font-size: 1.05rem;
  }

  .name:hover {
    color: var(--primary);
  }

  .username {
    font-size: 0.85rem;
    color: var(--text-muted);
    font-family: monospace;
  }

  .actions {
    display: flex;
    gap: 0.5rem;
  }

  .btn-sm {
    padding: 0.4rem 0.8rem;
    font-size: 0.8rem;
  }
</style>
