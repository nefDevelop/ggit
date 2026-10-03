<script lang="ts">
  import { authState } from "$lib/auth.svelte";
  import { toast } from "$lib/toast.svelte";
  import { apiFetch } from "$lib/api";
  import Avatar from "./Avatar.svelte";
  import LogOut from "lucide-svelte/icons/log-out";

  async function handleLogout() {
    await apiFetch("/api/v2/user/logout", { method: "POST" });
    window.location.href = "/user/login";
  }
</script>

{#if authState.isAuthenticated}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="user-profile-summary clickable-box"
    onclick={() => {
      toast.showNotices = !toast.showNotices;
      toast.clearUnread();
    }}
  >
    <div class="avatar-ring">
      <Avatar
        src={authState.user?.avatar_url}
        username={authState.user?.username}
        size={40}
      />
      {#if toast.hasUnread}
        <span class="unread-dot"></span>
      {/if}
    </div>
    <div class="user-details">
      <span class="full-name"
        >{authState.user?.full_name || authState.user?.username}</span
      >
      <span class="user-handle">@{authState.user?.username}</span>
    </div>
    <button class="minimal-logout" onclick={(e) => { e.stopPropagation(); handleLogout(); }} title="Logout" aria-label="Cerrar sesión">
      <LogOut size={20} />
    </button>
  </div>
{/if}

<style>
  .user-profile-summary {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    margin-bottom: 2rem;
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface);
    border-radius: 0.75rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.3s ease;
  }
  .user-profile-summary:hover {
    background: rgba(255, 255, 255, 0.06);
    transform: translateY(-2px);
  }
  @media (min-width: 1025px) {
    .clickable-box { cursor: pointer; }
  }
  .avatar-ring {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--primary-color), var(--tertiary));
    flex-shrink: 0;
  }
  .unread-dot {
    position: absolute;
    top: -2px;
    right: -2px;
    width: 12px;
    height: 12px;
    background: var(--md-sys-color-error);
    border-radius: 50%;
    border: 2px solid var(--md-sys-color-surface-container-highest);
    animation: pulse 1.5s infinite ease-in-out;
    z-index: 2;
  }
  @keyframes pulse {
    0% { transform: scale(0.9); opacity: 0.7; }
    50% { transform: scale(1.1); opacity: 1; }
    100% { transform: scale(0.9); opacity: 0.7; }
  }
  .user-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }
  .full-name {
    font-size: 0.9rem;
    font-weight: 800;
    color: inherit;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .user-handle {
    font-size: 0.75rem;
    color: inherit;
    opacity: 0.8;
    font-weight: 600;
  }
  .minimal-logout {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    border-radius: 0.75rem;
    transition: all 0.2s;
  }
  .minimal-logout:hover {
    color: var(--error);
    background: color-mix(in srgb, var(--status-danger) 10%, transparent);
  }
  @media (max-width: 1024px) {
    .user-profile-summary { display: none; }
  }
</style>
