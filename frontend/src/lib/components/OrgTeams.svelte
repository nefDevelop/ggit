<script lang="ts">
  import { onMount } from 'svelte';
  
  let { orgName } = $props();
  
  let teams = $state<any[]>([]);
  let loading = $state(true);
  let showModal = $state(false);
  let newTeam = $state({ name: '', description: '', permission: 'read' });

  async function fetchTeams() {
    loading = true;
    try {
      const res = await fetch(`/api/v2/orgs/${orgName}/teams`);
      if (res.ok) teams = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function createTeam() {
    if (!newTeam.name) return;
    try {
      const res = await fetch(`/api/v2/orgs/${orgName}/teams`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newTeam)
      });
      if (res.ok) {
        showModal = false;
        newTeam = { name: '', description: '', permission: 'read' };
        fetchTeams();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchTeams);
</script>

<div class="org-teams">
  <header class="section-actions">
     <button class="btn btn-primary" onclick={() => showModal = true}>
       <span class="material-symbols-rounded">add</span>
       New Team
     </button>
  </header>

  {#if showModal}
    <div class="modal-backdrop" onclick={() => showModal = false}>
      <div class="modal-content island tonal" onclick={(e) => e.stopPropagation()}>
        <header>
          <h3>Create New Team</h3>
          <p>Teams allow you to group members and assign permissions to repositories.</p>
        </header>
        <form onsubmit={(e) => { e.preventDefault(); createTeam(); }}>
          <div class="field">
            <label for="name">Team Name</label>
            <input id="name" type="text" bind:value={newTeam.name} placeholder="e.g. Frontend Developers" required />
          </div>
          <div class="field">
            <label for="desc">Description</label>
            <textarea id="desc" bind:value={newTeam.description} placeholder="What is this team for?"></textarea>
          </div>
          <div class="field">
            <label for="perm">Default Permission</label>
            <select id="perm" bind:value={newTeam.permission}>
              <option value="read">Read Access</option>
              <option value="write">Write Access</option>
              <option value="admin">Admin Access</option>
            </select>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-outline" onclick={() => showModal = false}>Cancel</button>
            <button type="submit" class="btn btn-primary" disabled={!newTeam.name}>Create Team</button>
          </div>
        </form>
      </div>
    </div>
  {/if}

  <div class="teams-grid">
    {#each teams as team}
      <div class="team-card island tonal hover-lift">
        <header class="team-header">
           <div class="team-title">
             <strong>{team.Name}</strong>
             <span class="badge" class:owner={team.Name === 'Owners'}>
               {team.Authorize === 1 ? 'Read' : team.Authorize === 2 ? 'Write' : 'Admin'}
             </span>
           </div>
           <p class="team-desc">{team.Description || 'No description provided.'}</p>
        </header>

        <div class="team-footer">
           <div class="team-stats">
              <div class="item">
                <span class="material-symbols-rounded">person</span>
                {team.NumMembers} Members
              </div>
              <div class="item">
                <span class="material-symbols-rounded">book</span>
                {team.NumRepos} Repos
              </div>
           </div>
           <a href="/org/{orgName}/teams/{team.ID}" class="btn btn-outline small">Manage</a>
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .org-teams { display: flex; flex-direction: column; gap: 1.5rem; }
  .section-actions { display: flex; justify-content: flex-end; }
  
  .teams-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.5rem; }
  
  .team-card { padding: 1.5rem; display: flex; flex-direction: column; gap: 1.5rem; border: 1px solid var(--border-color); }
  
  .team-title { display: flex; align-items: center; justify-content: space-between; }
  .team-title strong { font-size: 1.1rem; }
  
  .badge { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; background: rgba(255,255,255,0.05); padding: 0.2rem 0.6rem; border-radius: 4px; }
  .badge.owner { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); }

  .team-desc { font-size: 0.9rem; color: var(--text-muted); margin: 0.5rem 0 0 0; line-height: 1.4; }

  .team-footer { display: flex; align-items: center; justify-content: space-between; padding-top: 1rem; border-top: 1px solid var(--border-color); }
  
  .team-stats { display: flex; gap: 1rem; }
  .team-stats .item { display: flex; align-items: center; gap: 0.4rem; font-size: 0.8rem; color: var(--text-muted); }
  .team-stats .item span { font-size: 1.1rem; opacity: 0.7; }

  .small { padding: 0.4rem 0.8rem; font-size: 0.8rem; }

  /* Modal Styles */
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0,0,0,0.6);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal-content {
    width: 500px;
    max-width: 90vw;
    padding: 2.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    box-shadow: var(--shadow-xl);
  }

  .modal-content header h3 { margin: 0; font-size: 1.5rem; }
  .modal-content header p { color: var(--text-muted); font-size: 0.9rem; margin-top: 0.5rem; }

  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.85rem; font-weight: 700; color: var(--text-muted); }
  .field input, .field textarea, .field select {
    background: rgba(0,0,0,0.2);
    border: 1px solid var(--border-color);
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    color: #fff;
    font-family: inherit;
  }
  .field textarea { height: 100px; resize: none; }

  .modal-footer { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
</style>
