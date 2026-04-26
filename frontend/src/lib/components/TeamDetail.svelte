<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { dialog } from '$lib/dialog.svelte';
  import { page } from "$app/state";
  import Avatar from "$lib/components/Avatar.svelte";
  import { goto } from "$app/navigation";

  let { orgName, teamID } = $props();

  let team = $state<any>(null);
  let loading = $state(true);
  let newMember = $state("");
  let newRepo = $state("");
  let editing = $state(false);
  let editForm = $state({ Name: "", Description: "", Permission: "read" });

  async function fetchTeam() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/orgs/teams/${teamID}`);
      if (res.ok) {
        team = await res.json();
        editForm = {
          Name: team.Name,
          Description: team.Description,
          Permission:
            team.Authorize === 1
              ? "read"
              : team.Authorize === 2
                ? "write"
                : "admin",
        };
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function updateTeam() {
    try {
      const res = await apiFetch(`/api/v2/orgs/teams/${teamID}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(editForm),
      });
      if (res.ok) {
        editing = false;
        fetchTeam();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function deleteTeam() {
    if (
      !await dialog.confirm(
        "Are you sure you want to delete this team? This action is permanent.",
      )
    )
      return;
    try {
      const res = await apiFetch(`/api/v2/orgs/teams/${teamID}`, {
        method: "DELETE",
      });
      if (res.ok) {
        goto(`/org/${orgName}/teams`);
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function addMember() {
    if (!newMember) return;
    try {
      const res = await apiFetch(
        `/api/v2/orgs/teams/${teamID}/members/${newMember}`,
        { method: "POST" },
      );
      if (res.ok) {
        newMember = "";
        fetchTeam();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function removeMember(username: string) {
    try {
      const res = await apiFetch(
        `/api/v2/orgs/teams/${teamID}/members/${username}`,
        { method: "DELETE" },
      );
      if (res.ok) fetchTeam();
    } catch (e) {
      console.error(e);
    }
  }

  async function addRepo() {
    if (!newRepo) return;
    try {
      const res = await apiFetch(
        `/api/v2/orgs/teams/${teamID}/repos/${newRepo}`,
        { method: "POST" },
      );
      if (res.ok) {
        newRepo = "";
        fetchTeam();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function removeRepo(repoID: number) {
    try {
      const res = await apiFetch(
        `/api/v2/orgs/teams/${teamID}/repos/${repoID}`,
        { method: "DELETE" },
      );
      if (res.ok) fetchTeam();
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchTeam);
</script>

{#if loading}
  <div class="skeleton-detail tonal"></div>
{:else if team}
  <div class="team-detail-view">
    <header class="detail-header tonal">
      <div class="title-group">
        <a href="/org/{orgName}/teams" class="back-link">
          <span class="material-symbols-rounded">arrow_back</span>
          Back to Teams
        </a>
        {#if editing}
          <input type="text" class="edit-title" bind:value={editForm.Name} />
          <textarea class="edit-desc" bind:value={editForm.Description}
          ></textarea>
        {:else}
          <h1>{team.Name}</h1>
          <p>{team.Description || "No description provided."}</p>
        {/if}
      </div>
      <div class="header-actions">
        {#if editing}
          <button class="btn btn-outline" onclick={() => (editing = false)}
            >Cancel</button
          >
          <button class="btn btn-primary" onclick={updateTeam}
            >Save Changes</button
          >
        {:else}
          {#if team.Name !== "Owners"}
            <button
              class="btn btn-icon-text btn-danger-text"
              onclick={deleteTeam}
            >
              <span class="material-symbols-rounded">delete</span>
              Delete Team
            </button>
          {/if}
          <button class="btn btn-primary" onclick={() => (editing = true)}>
            <span class="material-symbols-rounded">edit</span>
            Edit Team
          </button>
        {/if}
      </div>
    </header>

    {#if editing}
      <div class="edit-settings">
        <div class="field">
          <label for="perm">Team Permission Level</label>
          <select id="perm" bind:value={editForm.Permission}>
            <option value="read">Read: Can view and clone repositories</option>
            <option value="write">Write: Can push to repositories</option>
            <option value="admin">Admin: Full control over repositories</option>
          </select>
          <p class="help">
            This permission applies to all repositories assigned to this team.
          </p>
        </div>
      </div>
    {/if}

    <div class="management-grid">
      <!-- Members Section -->
      <section class="management-card">
        <header class="card-header">
          <div class="header-title">
            <span class="material-symbols-rounded">group</span>
            <h3>Members</h3>
          </div>
          <span class="count">{team.NumMembers}</span>
        </header>
        <div class="card-body">
          <form
            class="add-form"
            onsubmit={(e) => {
              e.preventDefault();
              addMember();
            }}
          >
            <input
              type="text"
              bind:value={newMember}
              placeholder="Add member by username..."
            />
            <button
              type="submit"
              class="btn btn-primary btn-sm"
              disabled={!newMember}>Add</button
            >
          </form>

          <div class="item-list">
            {#if team.Members && team.Members.length > 0}
              {#each team.Members as member}
                <div class="list-item">
                  <Avatar
                    src={member.avatar_url}
                    username={member.username}
                    size={32}
                    class="avatar-xs"
                  />
                  <div class="user-meta">
                    <span class="name">{member.username}</span>
                  </div>
                  {#if team.Name !== "Owners" || team.NumMembers > 1}
                    <button
                      class="btn btn-icon btn-danger-text"
                      title="Remove member"
                      onclick={() => removeMember(member.username)}
                    >
                      <span class="material-symbols-rounded">person_remove</span
                      >
                    </button>
                  {/if}
                </div>
              {/each}
            {:else}
              <p class="empty-list">No members in this team.</p>
            {/if}
          </div>
        </div>
      </section>

      <!-- Repositories Section -->
      <section class="management-card">
        <header class="card-header">
          <div class="header-title">
            <span class="material-symbols-rounded">book</span>
            <h3>Repositories</h3>
          </div>
          <span class="count">{team.NumRepos}</span>
        </header>
        <div class="card-body">
          <form
            class="add-form"
            onsubmit={(e) => {
              e.preventDefault();
              addRepo();
            }}
          >
            <input
              type="text"
              bind:value={newRepo}
              placeholder="Add repository by name..."
            />
            <button
              type="submit"
              class="btn btn-primary btn-sm"
              disabled={!newRepo}>Add</button
            >
          </form>

          <div class="item-list">
            {#if team.Repos && team.Repos.length > 0}
              {#each team.Repos as repo}
                <div class="list-item">
                  <span class="material-symbols-rounded repo-icon">book</span>
                  <div class="repo-meta">
                    <span class="name">{repo.name}</span>
                  </div>
                  <button
                    class="btn btn-icon btn-danger-text"
                    title="Remove repository"
                    onclick={() => removeRepo(repo.id)}
                  >
                    <span class="material-symbols-rounded">close</span>
                  </button>
                </div>
              {/each}
            {:else}
              <p class="empty-list">No repositories assigned to this team.</p>
            {/if}
          </div>
        </div>
      </section>
    </div>
  </div>
{/if}

<style>
  .team-detail-view {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    animation: fadeIn 0.4s ease-out;
  }

  .detail-header {
    padding: 2.5rem;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }
  .back-link {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    text-decoration: none;
    color: var(--text-muted);
    font-weight: 700;
    margin-bottom: 1rem;
    font-size: 0.9rem;
  }
  .back-link:hover {
    color: var(--primary-color);
  }

  .title-group h1 {
    margin: 0;
    font-size: 2.25rem;
    font-weight: 800;
  }
  .title-group p {
    color: var(--text-muted);
    margin: 0.5rem 0 0 0;
    font-size: 1.1rem;
  }

  .edit-title {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid var(--primary-color);
    font-size: 2.25rem;
    font-weight: 800;
    color: #fff;
    width: 100%;
    border-radius: 0.5rem;
    padding: 0.5rem;
  }
  .edit-desc {
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid var(--md-sys-color-outline-variant);
    font-size: 1.1rem;
    color: var(--text-muted);
    width: 100%;
    border-radius: 0.5rem;
    padding: 0.5rem;
    margin-top: 0.5rem;
    height: 80px;
  }

  .header-actions {
    display: flex;
    gap: 1rem;
  }

  .edit-settings {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    margin-top: -1rem;
  }
  .edit-settings .help {
    font-size: 0.85rem;
    color: var(--text-muted);
    margin: 0.25rem 0 0 0;
  }

  .management-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }
  @media (max-width: 800px) {
    .management-grid {
      grid-template-columns: 1fr;
    }
  }

  .management-card {
    padding: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  .card-header {
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    background: rgba(255, 255, 255, 0.02);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .header-title {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }
  .header-title span {
    color: var(--primary-color);
  }
  .card-header h3 {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 700;
  }
  .count {
    font-size: 0.75rem;
    font-weight: 800;
    color: var(--text-muted);
    background: rgba(255, 255, 255, 0.05);
    padding: 0.2rem 0.6rem;
    border-radius: 1rem;
  }

  .card-body {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .add-form {
    display: flex;
    gap: 1rem;
  }
  .add-form input {
    flex: 1;
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 0.6rem 1rem;
    border-radius: 0.75rem;
    color: #fff;
    font-family: inherit;
  }

  .item-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .list-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem 1rem;
    background: rgba(255, 255, 255, 0.02);
    border-radius: 0.75rem;
    transition: all 0.2s;
  }
  .list-item:hover {
    background: rgba(255, 255, 255, 0.04);
  }

  .repo-icon {
    font-size: 1.25rem;
    color: var(--primary-color);
    opacity: 0.6;
  }

  .user-meta,
  .repo-meta {
    flex: 1;
    display: flex;
    flex-direction: column;
  }
  .name {
    font-weight: 600;
    font-size: 0.95rem;
  }

  .empty-list {
    text-align: center;
    color: var(--text-muted);
    font-style: italic;
    padding: 1rem;
  }

  .btn-danger-text {
    color: #ef4444;
  }
  .btn-danger-text:hover {
    background: rgba(239, 68, 68, 0.1);
  }

  .skeleton-detail {
    height: 400px;
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 0.3;
    }
    50% {
      opacity: 0.5;
    }
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

  .field select {
    background: rgba(0, 0, 0, 0.2);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    color: #fff;
    width: 100%;
  }

  /* Responsive adjustment */
  @media (max-width: 900px) {
    .management-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
