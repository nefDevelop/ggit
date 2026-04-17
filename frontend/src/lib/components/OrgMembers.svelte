<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api';
  
  let { orgName } = $props();
  
  let members = $state<any[]>([]);
  let loading = $state(true);
  let error = $state('');

  async function fetchMembers() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}/members`);
      if (res.ok) {
        members = await res.json();
      } else {
        error = 'Failed to load members.';
      }
    } catch (e) {
      console.error(e);
      error = 'Network error.';
    } finally {
      loading = false;
    }
  }

  async function removeMember(username: string) {
    if (!confirm(`Are you sure you want to remove ${username} from the organization?`)) return;
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}/members/${username}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        await fetchMembers();
      } else {
        const data = await res.json();
        alert(data.message || 'Failed to remove member');
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchMembers);
</script>

<div class="org-members card-surface">
  <div class="header-actions">
    <h3 class="section-label">Organization Members</h3>
    <button class="btn btn-primary btn-sm">
      <span class="material-symbols-rounded">person_add</span> Invite Member
    </button>
  </div>

  <div class="members-list">
    {#if loading}
      <div class="skeletons">
        {#each Array(4) as _}<div class="skeleton-org"></div>{/each}
      </div>
    {:else if error}
      <div class="error-state island tonal-error">
        <span class="material-symbols-rounded">error</span>
        <p>{error}</p>
      </div>
    {:else if members.length === 0}
      <div class="empty-state island tonal">
        <span class="material-symbols-rounded">person_off</span>
        <p>This organization has no members yet.</p>
      </div>
    {:else}
      <div class="members-grid">
        {#each members as member}
          <div class="member-card island tonal hover-lift">
             <div class="member-main">
                <img src={member.avatar_url} alt="" class="member-avatar" />
                <div class="member-info">
                   <a href="/{member.username}" class="member-name">{member.full_name || member.username}</a>
                   <span class="member-user">@{member.username}</span>
                </div>
             </div>
             <div class="member-actions">
                <button class="btn btn-ghost error btn-sm" onclick={() => removeMember(member.username)}>
                   <span class="material-symbols-rounded">person_remove</span> Remove
                </button>
             </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .org-members { display: flex; flex-direction: column; gap: 1.5rem; }
  .header-actions { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--border-color); padding-bottom: 1rem; }
  
  .section-label { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); margin: 0; }

  .members-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1.5rem; }
  .member-card { display: flex; align-items: center; justify-content: space-between; padding: 1.25rem 1.5rem; }
  
  .member-main { display: flex; align-items: center; gap: 1rem; }
  .member-avatar { width: 44px; height: 44px; border-radius: 12px; }
  
  .member-info { display: flex; flex-direction: column; }
  .member-name { font-weight: 700; color: #fff; text-decoration: none; font-size: 1.05rem; }
  .member-name:hover { color: var(--primary-color); }
  .member-user { font-size: 0.85rem; color: var(--text-muted); }

  .skeleton-org { height: 72px; background: var(--md3-surface-3); border-radius: 1rem; margin-bottom: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }

  .empty-state, .error-state { text-align: center; padding: 4rem; display: flex; flex-direction: column; align-items: center; gap: 1rem; border-radius: 1rem; }
  .empty-state span, .error-state span { font-size: 3rem; opacity: 0.2; }

  .btn-ghost.error { color: #ef4444; }
  .btn-ghost.error:hover { background: rgba(239, 68, 68, 0.1); }
</style>
