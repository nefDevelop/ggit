<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import RepoRow from "$lib/components/RepoRow.svelte";
  import UserRow from "$lib/components/UserRow.svelte";
  import GistRow from "$lib/components/GistRow.svelte";
  import Avatar from "$lib/components/Avatar.svelte";
  import Card from "$lib/components/ui/Card.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";

  const username = $derived(page.params.username);
  let user = $state<any>(null);
  let orgs = $state<any[]>([]);
  let activeTab = $state("repos"); // 'repos', 'stars'
  let data = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  async function fetchProfile() {
    if (!username || username === 'undefined') {
       error = "Invalid user";
       return;
    }
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
    if (!username || username === 'undefined') return;
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
    if (username && username !== 'undefined') {
      try {
        const res = await apiFetch(`/api/v2/users/${username}/orgs`);
        if (res.ok) orgs = await res.json();
      } catch (e) {
        console.error(e);
      }
    }
  });

  $effect(() => {
    if (username && username !== 'undefined') {
      fetchProfile();
    }
  });

  $effect(() => {
    if (username && username !== 'undefined') {
      fetchTabData(activeTab);
    }
  });
</script>

<svelte:head>
  <title>{username} · Ggit</title>
</svelte:head>

<div class="profile-page">
  {#if error && !user}
    <div class="error-container">
      <Card variant="low" padding="lg">
        <div class="error-msg">
          <span class="material-symbols-rounded">error</span>
          <h3>{error}</h3>
          <p>The user you are looking for does not exist or has been removed.</p>
          <Button variant="primary" href="/explore/users">Explore Users</Button>
        </div>
      </Card>
    </div>
  {:else if user}
    <div class="profile-hero">
      <div class="hero-bg"></div>
      <div class="hero-content">
        <header class="profile-header">
          <div class="avatar-container">
            <img src={user.avatar_url} alt={user.username} class="main-avatar" loading="lazy" decoding="async" />
          </div>
          <div class="user-main-info">
            <div class="title-row">
              <h1>{user.full_name || user.username}</h1>
              <Badge variant="primary">User</Badge>
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
                  <a href="/org/{org.username}" class="org-mini" title={org.full_name || org.username}>
                    <Avatar src={org.avatar_url} username={org.username} size={32} />
                  </a>
                {/each}
              </div>
            {/if}
          </div>
        </header>
      </div>
    </div>

    <div class="profile-container">
      <aside class="profile-sidebar">
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
          <button
            class:active={activeTab === "gists"}
            onclick={() => fetchTabData("gists")}
          >
            <span class="material-symbols-rounded">code</span>
            Gists
          </button>
        </nav>
      </aside>

      <main class="profile-content">
        {#if loading}
          <div class="loading-stack">
            {#each Array(4) as _}
              <Card variant="low" class="skeleton-row" />
            {/each}
          </div>
        {:else if data.length === 0}
          <Card variant="low" padding="lg" class="empty-state">
            <div class="empty-icon-circle">
              <span class="material-symbols-rounded">folder_off</span>
            </div>
            <h3>No items found</h3>
            <p>This user hasn't added anything here yet.</p>
          </Card>
        {:else}
          <div class="repo-stack">
            {#each data as item}
              <Card hover variant="low" padding="none" class="row-card">
                {#if activeTab === "repos" || activeTab === "stars"}
                  <RepoRow repo={item} />
                {:else if activeTab === "gists"}
                  <GistRow gist={item} />
                {:else}
                  <UserRow user={item} />
                {/if}
              </Card>
            {/each}
          </div>
        {/if}
      </main>
    </div>
  {/if}
</div>

<style>
  .profile-page { min-height: 100vh; background: var(--bg-color); animation: fadeIn 0.4s ease-out; }
  
  .profile-hero { position: relative; padding: 5rem 2rem 3rem; overflow: hidden; background: var(--md-sys-color-surface-container-low); border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .hero-bg { 
    position: absolute; top: 0; left: 0; right: 0; height: 100%;
    background: linear-gradient(135deg, var(--primary-color) 0%, var(--status-accent) 100%);
    opacity: 0.1; z-index: 0;
    mask-image: linear-gradient(to bottom, black, transparent);
  }
  
  .hero-content { position: relative; z-index: 1; max-width: 1100px; margin: 0 auto; }
  .profile-header { display: flex; align-items: flex-end; gap: 2.5rem; }
  
  .main-avatar { 
    width: 160px; height: 160px; border-radius: 2.5rem; 
    box-shadow: var(--shadow-xl);
    border: 4px solid var(--md-sys-color-surface);
    background: var(--md-sys-color-surface);
    object-fit: cover;
  }
  
  .user-main-info { flex: 1; padding-bottom: 0.5rem; }
  .title-row { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.5rem; }
  .title-row h1 { margin: 0; font-size: 2.5rem; font-weight: 800; color: var(--md-sys-color-on-surface); letter-spacing: -0.02em; }
  
  .meta-row { display: flex; align-items: center; gap: 0.75rem; color: var(--md-sys-color-on-surface-variant); font-size: 1rem; margin-bottom: 1rem; }
  .username { color: var(--primary-color); font-weight: 600; }
  .sep { opacity: 0.3; }
  
  .bio { font-size: 1.1rem; line-height: 1.6; color: var(--md-sys-color-on-surface); opacity: 0.8; max-width: 600px; margin-bottom: 1.5rem; }
  
  .org-strip { display: flex; gap: 0.8rem; }
  .org-mini { transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1); display: flex; border-radius: 8px; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  .org-mini:hover { transform: scale(1.15) translateY(-2px); border-color: var(--primary-color); }

  .profile-container { max-width: 1100px; margin: 2rem auto; padding: 0 2rem; display: flex; gap: 3rem; }
  
  .profile-sidebar { width: 260px; flex-shrink: 0; }
  .profile-nav { 
    display: flex; flex-direction: column; gap: 0.25rem; 
    padding: 0.75rem; border-radius: 1.5rem; height: fit-content; position: sticky; top: 2rem;
    background: var(--md-sys-color-surface-container-low);
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  
  .profile-nav button {
    display: flex; align-items: center; gap: 1rem; padding: 0.85rem 1.25rem;
    border: none; background: transparent; color: var(--md-sys-color-on-surface-variant);
    font-weight: 700; cursor: pointer; border-radius: 1rem; transition: all 0.2s;
    font-size: 0.95rem;
  }
  
  .profile-nav button:hover { background: var(--md-sys-color-surface-variant); color: var(--md-sys-color-on-surface); }
  .profile-nav button.active { background: var(--primary-color); color: white; box-shadow: 0 8px 16px rgba(var(--primary-rgb), 0.3); }

  .profile-content { flex: 1; min-width: 0; }
  .repo-stack, .loading-stack { display: flex; flex-direction: column; gap: 1rem; }
  

  :global(.skeleton-row) { height: 100px; animation: pulse 1.5s infinite; background: var(--md-sys-color-surface-container-low); border-radius: 1.25rem; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  :global(.empty-state) { text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-icon-circle { width: 80px; height: 80px; border-radius: 50%; background: var(--md-sys-color-surface-container-highest); display: flex; align-items: center; justify-content: center; color: var(--md-sys-color-on-surface-variant); opacity: 0.5; margin-bottom: 1rem; }
  .empty-icon-circle span { font-size: 2.5rem; }

  .error-container { max-width: 600px; margin: 10rem auto; }
  .error-msg { text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .error-msg span { font-size: 4rem; color: var(--md-sys-color-error); }
  .error-msg h3 { margin: 0; color: var(--md-sys-color-on-surface); }
  .error-msg p { color: var(--md-sys-color-on-surface-variant); margin-bottom: 1.5rem; }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @media (max-width: 950px) {
    .profile-container { flex-direction: column; padding: 0 1rem; }
    .profile-sidebar { width: 100%; }
    .profile-nav { width: 100%; flex-direction: row; position: static; }
    .profile-nav button { flex: 1; justify-content: center; }
  }
  
  @media (max-width: 600px) {
    .profile-header { flex-direction: column; align-items: center; text-align: center; gap: 1.5rem; }
    .title-row { justify-content: center; }
    .meta-row { flex-direction: column; gap: 0.25rem; }
    .org-strip { justify-content: center; }
    .main-avatar { width: 140px; height: 140px; margin-top: -2rem; }
    .title-row h1 { font-size: 2rem; }
  }
</style>
