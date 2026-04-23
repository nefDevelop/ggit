<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import RepoRow from "$lib/components/RepoRow.svelte";
  import UserRow from "$lib/components/UserRow.svelte";

  const username = page.params.username;
  let user = $state<any>(null);
  let activeTab = $state("repos"); // 'repos', 'stars', 'followers', 'following'
  let data = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  async function fetchProfile() {
    try {
      const res = await apiFetch(`/api/v2/users/${username}`);
      if (res.ok) {
        user = await res.json();
      } else {
        error = "Usuario no encontrado";
      }
    } catch (e) {
      error = "Error de conexión";
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
      error = "Error al cargar datos";
    } finally {
      loading = false;
    }
  }

  onMount(async () => {
    await fetchProfile();
    if (user) await fetchTabData("repos");
  });
</script>

<div class="profile-page">
  {#if error && !user}
    <div class="error-msg tonal">{error}</div>
  {:else if user}
    <header class="profile-header tonal">
      <img src={user.avatar_url} alt={user.username} class="avatar" />
      <div class="info">
        <div class="title-row">
          <h1>{user.full_name || user.username}</h1>
          <span class="badge">Usuario</span>
        </div>
        <p class="username">@{user.username}</p>
        <p class="email">{user.email}</p>
      </div>
    </header>

    <div class="profile-tabs">
      <button
        class:active={activeTab === "repos"}
        onclick={() => fetchTabData("repos")}
      >
        <span class="material-symbols-rounded">folder</span>
        Repositorios
      </button>
      <button
        class:active={activeTab === "stars"}
        onclick={() => fetchTabData("stars")}
      >
        <span class="material-symbols-rounded">star</span>
        Estrellas
      </button>
      <button
        class:active={activeTab === "followers"}
        onclick={() => fetchTabData("followers")}
      >
        <span class="material-symbols-rounded">group</span>
        Seguidores
      </button>
      <button
        class:active={activeTab === "following"}
        onclick={() => fetchTabData("following")}
      >
        <span class="material-symbols-rounded">person_add</span>
        Siguiendo
      </button>
    </div>

    <section class="content-section">
      {#if loading}
        <div class="loading-grid">
          {#each Array(3) as _}
            <div class="skeleton tonal"></div>
          {/each}
        </div>
      {:else if data.length === 0}
        <div class="empty-state tonal">
          <span class="material-symbols-rounded">sentiment_dissatisfied</span>
          <p>No hay elementos para mostrar en esta sección.</p>
        </div>
      {:else}
        <div class="data-grid">
          {#each data as item}
            {#if activeTab === "repos" || activeTab === "stars"}
              <RepoRow repo={item} />
            {:else}
              <UserRow user={item} />
            {/if}
          {/each}
        </div>
      {/if}
    </section>
  {/if}
</div>

<style>
  .profile-page {
    max-width: 1000px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding: 1rem;
    animation: fadeIn 0.4s ease-out;
  }

  .profile-header {
    display: flex;
    align-items: center;
    gap: 2rem;
    padding: 2.5rem;
    background: var(--surface-variant);
  }

  .avatar {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    box-shadow: var(--shadow-lg);
    border: 4px solid var(--surface);
  }

  .title-row {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
  .badge {
    background: var(--primary-light);
    color: var(--primary);
    font-size: 0.75rem;
    font-weight: 800;
    padding: 0.2rem 0.6rem;
    border-radius: 20px;
    text-transform: uppercase;
  }
  .info h1 {
    margin: 0;
    font-size: 2.25rem;
    font-weight: 800;
  }
  .username {
    color: var(--primary);
    font-family: monospace;
    font-size: 1.1rem;
    margin-top: 0.25rem;
  }
  .email {
    color: var(--text-muted);
    font-size: 0.9rem;
    margin-top: 0.5rem;
  }

  .profile-tabs {
    display: flex;
    gap: 0.5rem;
    background: var(--surface-variant);
    padding: 0.4rem;
    border-radius: 1rem;
    align-self: flex-start;
  }

  .profile-tabs button {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.6rem 1.2rem;
    border: none;
    background: transparent;
    color: var(--text-muted);
    font-weight: 700;
    cursor: pointer;
    border-radius: 0.75rem;
    transition: all 0.2s;
  }

  .profile-tabs button:hover {
    background: rgba(255, 255, 255, 0.05);
  }
  .profile-tabs button.active {
    background: var(--surface);
    color: var(--primary);
    box-shadow: var(--shadow-sm);
  }

  .content-section {
    min-height: 300px;
  }

  .data-grid,
  .loading-grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .empty-state {
    padding: 4rem;
    text-align: center;
    color: var(--text-muted);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .empty-state span {
    font-size: 3rem;
    opacity: 0.3;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
