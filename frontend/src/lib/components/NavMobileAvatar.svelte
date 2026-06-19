<script lang="ts">
  import { authState } from "$lib/auth.svelte";
  import { toast } from "$lib/toast.svelte";
  import Avatar from "./Avatar.svelte";
</script>

{#if authState.isAuthenticated}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="nav-item mobile-avatar-item"
    onclick={() => {
      toast.showNotices = !toast.showNotices;
      toast.clearUnread();
    }}
  >
    <div class="avatar-ring-small">
      <Avatar src={authState.user?.avatar_url} username={authState.user?.username} size={28} />
      {#if toast.hasUnread}
        <span class="unread-dot-small"></span>
      {/if}
    </div>
  </div>
{/if}

<style>
  .mobile-avatar-item { display: none; }
  @media (max-width: 1024px) {
    .mobile-avatar-item {
      display: flex !important;
      align-items: center; justify-content: center;
      cursor: pointer;
      padding: 0.5rem;
    }
    .avatar-ring-small {
      position: relative; display: flex; align-items: center; justify-content: center;
    }
    .unread-dot-small {
      position: absolute; top: -1px; right: -1px;
      width: 8px; height: 8px;
      background: var(--md-sys-color-error); border-radius: 50%;
      border: 1.5px solid var(--md-sys-color-surface-container-highest);
      animation: pulse 1.5s infinite ease-in-out; z-index: 2;
    }
    @keyframes pulse {
      0% { transform: scale(0.9); opacity: 0.7; }
      50% { transform: scale(1.1); opacity: 1; }
      100% { transform: scale(0.9); opacity: 0.7; }
    }
  }
</style>
