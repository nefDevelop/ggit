<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import RepoRow from "$lib/components/RepoRow.svelte";
  import UserRow from "$lib/components/UserRow.svelte";
  import Avatar from "$lib/components/Avatar.svelte";

  const username = $derived(page.params.username);
  let user = $state<any>(null);
  let orgs = $state<any[]>([]);
  let activeTab = $state("repos"); // 'repos', 'stars'
  let data = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  async function fetchProfile() {
    try {
      const res = await apiFetch(`/api/v2/users/${username}`);
      if (res.ok) {
        user = await res.json();
      } else {
        error = "User not found";
      }
    } catch (e) {
      error = "Connection error";
    }
  }

  async function fetchTabData(tab: string) {
    loading = true;
    activeTab = tab;
    try {
      const res = await apiFetch(`/api/v2/users/${username}/${tab}`);
      if (res.ok) {
        data = await res.json();
      } else {
        data = [];
      }
    } catch (e) {
      error = "Failed to load data";
    } finally {
      loading = false;
    }
  }

  onMount(async () => {
    await fetchProfile();
    if (user) {
      await fetchTabData("repos");
      // Fetch public orgs for this user
      try {
        const res = await apiFetch(`/api/v2/users/${username}/orgs`);
        if (res.ok) orgs = await res.json();
      } catch (e) {
        console.error(e);
      }
    }
  });

  // Re-fetch when user changes
  $effect(() => {
    if (username) {
        fetchProfile().then(() => fetchTabData("repos"));
    }
  });
</script>

<div class="profile-page">
  {#if error && !user}
    <div class="error-msg tonal">{error}</div>
  {:else if user}
    <div class="profile-hero">
      <div class="hero-bg"></div>
      <div class="hero-content">
        <header class="profile-header">
          <div class="avatar-container">
            <img src={user.avatar_url} alt={user.username} class="main-avatar" />
          </div>
          <div class="user-main-info">
            <div class="title-row">
              <h1>{user.full_name || user.username}</h1>
              <span class="user-badge">User</span>
            </div>
            <div class="meta-row">
              <span class="username">@{user.username}</span>
              {#if user.email}
                <span class="sep">•</span>
                <span class="email">{user.email}</span>
              {/if}
            </div>
            {#if user.bio}
              <p class="bio">{user.bio}</p>
            {/if}

            {#if orgs.length > 0}
              <div class="org-strip">
                {#each orgs as org}
                  <a href="/{org.username}" class="org-mini" title={org.full_name || org.username}>
                    <Avatar src={org.avatar_url} username={org.username} size={28} />
                  </a>
                {/each}
              </div>
            {/if}
          </div>
        </header>
      </div>
    </div>

    <div class="profile-container">
      <nav class="profile-nav glass">
        <button
          class:active={activeTab === "repos"}
          onclick={() => fetchTabData("repos")}
        >
          <span class="material-symbols-rounded">folder</span>
          Repositories
        </button>
        <button
          class:active={activeTab === "stars"}
          onclick={() => fetchTabData("stars")}
        >
          <span class="material-symbols-rounded">star</span>
          Stars
        </button>
      </nav>

      <main class="profile-content">
        {#if loading}
          <div class="loading-stack">
            {#each Array(4) as _}
              <div class="skeleton-row glass"></div>
            {/each}
          </div>
        {:else if data.length === 0}
          <div class="empty-state glass">
            <div class="empty-icon-circle">
              <span class="material-symbols-rounded">folder_off</span>
            </div>
            <h3>No items found</h3>
            <p>This user hasn't added anything here yet.</p>
          </div>
        {:else}
          <div class="repo-stack">
            {#each data as item}
              {#if activeTab === "repos" || activeTab === "stars"}
                <div class="row-container glass">
                  <RepoRow repo={item} />
                </div>
              {:else}
                <div class="row-container glass">
                  <UserRow user={item} />
                </div>
              {/if}
            {/each}
          </div>
        {/if}
      </main>
    </div>
  {/if}
</div>

<style>
  .profile-page { min-height: 100vh; background: var(--bg-color); animation: fadeIn 0.4s ease-out; }
  
  .profile-hero { position: relative; padding: 4rem 2rem 2rem; overflow: hidden; background: var(--md3-surface-2); border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .hero-bg { 
    position: absolute; top: 0; left: 0; right: 0; height: 100%;
    background: linear-gradient(135deg, var(--primary-color) 0%, #6366f1 100%);
    opacity: 0.1; z-index: 0;
    mask-image: linear-gradient(to bottom, black, transparent);
  }
  
  .hero-content { position: relative; z-index: 1; max-width: 1100px; margin: 0 auto; }
  
  .profile-header { display: flex; align-items: flex-end; gap: 2.5rem; }
  
  .avatar-container { flex-shrink: 0; }
  .main-avatar { 
    width: 160px; height: 160px; border-radius: 2.5rem; 
    box-shadow: 0 20px 40px rgba(0,0,0,0.2); 
    border: 4px solid var(--surface);
    background: var(--surface);
    object-fit: cover;
  }
  
  .user-main-info { flex: 1; padding-bottom: 0.5rem; }
  .title-row { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.5rem; }
  .title-row h1 { margin: 0; font-size: 2.5rem; font-weight: 800; color: var(--text-color); letter-spacing: -0.02em; }
  .user-badge { 
    background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color);
    font-size: 0.7rem; font-weight: 800; padding: 0.25rem 0.75rem; border-radius: 2rem;
    text-transform: uppercase; border: 1px solid rgba(var(--primary-rgb), 0.2);
  }
  
  .meta-row { display: flex; align-items: center; gap: 0.75rem; color: var(--text-muted); font-size: 1rem; margin-bottom: 1rem; }
  .username { color: var(--primary-color); font-weight: 600; }
  .sep { opacity: 0.3; }
  
  .bio { font-size: 1.1rem; line-height: 1.6; color: var(--text-color); opacity: 0.8; max-width: 600px; margin-bottom: 1.5rem; }
  
  .org-strip { display: flex; gap: 0.8rem; }
  .org-mini { transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1); display: flex; }
  .org-mini:hover { transform: scale(1.15) translateY(-2px); }

  .profile-container { max-width: 1100px; margin: 2rem auto; padding: 0 2rem; display: flex; gap: 2.5rem; }
  
  .profile-nav { 
    width: 240px; flex-shrink: 0; display: flex; flex-direction: column; gap: 0.4rem; 
    padding: 0.75rem; border-radius: 1.25rem; height: fit-content; position: sticky; top: 2rem;
  }
  
  .profile-nav button {
    display: flex; align-items: center; gap: 1rem; padding: 0.85rem 1.1rem;
    border: none; background: transparent; color: var(--text-muted);
    font-weight: 700; cursor: pointer; border-radius: 0.85rem; transition: all 0.2s;
    font-size: 0.95rem;
  }
  
  .profile-nav button:hover { background: rgba(255, 255, 255, 0.05); color: var(--text-color); }
  .profile-nav button.active { background: var(--primary-color); color: white; box-shadow: 0 8px 16px rgba(var(--primary-rgb), 0.3); }

  .profile-content { flex: 1; min-width: 0; }
  
  .repo-stack, .loading-stack { display: flex; flex-direction: column; gap: 1rem; }
  
  .row-container { border-radius: 1.25rem; transition: transform 0.2s ease, box-shadow 0.2s; overflow: hidden; background: var(--surface); border: 1px solid var(--md-sys-color-outline-variant); }
  .row-container:hover { transform: translateY(-2px); box-shadow: var(--shadow-md); border-color: var(--primary-color); }

  .skeleton-row { height: 80px; border-radius: 1.25rem; animation: pulse 1.5s infinite; background: rgba(255,255,255,0.05); }
  @keyframes pulse { 0%, 100% { opacity: 0.4; } 50% { opacity: 0.7; } }

  .empty-state { padding: 5rem 2rem; text-align: center; border-radius: 2rem; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-icon-circle { width: 80px; height: 80px; border-radius: 50%; background: rgba(255,255,255,0.03); display: flex; align-items: center; justify-content: center; color: var(--text-muted); opacity: 0.5; margin-bottom: 1rem; }
  .empty-icon-circle span { font-size: 2.5rem; }
  .empty-state h3 { margin: 0; font-size: 1.5rem; font-weight: 700; }
  .empty-state p { color: var(--text-muted); margin: 0; }

  .glass { background: rgba(255, 255, 255, 0.01); backdrop-filter: blur(12px); border: 1px solid rgba(255, 255, 255, 0.04); }

  .error-msg { padding: 2rem; text-align: center; color: var(--md-sys-color-error); }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @media (max-width: 950px) {
    .profile-container { flex-direction: column; padding: 0 1rem; }
    .profile-nav { width: 100%; flex-direction: row; position: static; }
    .profile-nav button { flex: 1; justify-content: center; }
  }
  
  @media (max-width: 600px) {
    .profile-header { flex-direction: column; align-items: center; text-align: center; gap: 1.5rem; }
    .title-row { justify-content: center; }
    .meta-row { flex-direction: column; gap: 0.25rem; }
    .meta-row .sep { display: none; }
    .org-strip { justify-content: center; }
    .avatar-container { margin-top: -2rem; }
    .main-avatar { width: 140px; height: 140px; }
    .title-row h1 { font-size: 2rem; }
  }
</style>
