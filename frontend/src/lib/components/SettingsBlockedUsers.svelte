<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  import Avatar from "./Avatar.svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import Input from "./ui/Input.svelte";
  import { UserMinus, UserPlus, ShieldAlert } from "lucide-svelte";

  let blockedUsers = $state<any[]>([]);
  let loading = $state(true);
  let newUsername = $state("");
  let submitting = $state(false);

  async function fetchBlocks() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/blocks");
      if (res.ok) blockedUsers = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function blockUser() {
    if (!newUsername.trim()) return;
    submitting = true;
    try {
      const res = await apiFetch(`/api/v2/user/settings/blocks/${newUsername}`, {
        method: "POST"
      });
      if (res.ok) {
        toast.success(`Usuario ${newUsername} bloqueado.`);
        newUsername = "";
        fetchBlocks();
      } else {
        const data = await res.json();
        toast.error(data.error || "No se pudo bloquear al usuario.");
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  async function unblockUser(username: string) {
    if (!await dialog.confirm(`¿Desbloquear a ${username}? Volverá a tener permiso para interactuar contigo.`)) return;
    
    try {
      const res = await apiFetch(`/api/v2/user/settings/blocks/${username}`, {
        method: "DELETE"
      });
      if (res.ok) {
        toast.success("Usuario desbloqueado.");
        fetchBlocks();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchBlocks);
</script>

<div class="blocks-settings">
  <div class="settings-blocked-container">


    <div class="add-block-form tonal">
      <div class="form-row">
        <Input 
          label="Bloquear Usuario" 
          bind:value={newUsername} 
          placeholder="Escribe un nombre de usuario..." 
          class="grow"
        />
        <Button 
          variant="primary" 
          onclick={blockUser} 
          disabled={submitting || !newUsername.trim()}
          icon="person_add_disabled"
        >
          {submitting ? 'Bloqueando...' : 'Bloquear'}
        </Button>
      </div>
    </div>

    <div class="blocks-list">
      {#if loading}
        <div class="skeleton-list"></div>
      {:else if blockedUsers.length === 0}
        <div class="empty-state">
           <ShieldAlert size={48} />
           <p>No tienes a nadie bloqueado. ¡Genial!</p>
        </div>
      {:else}
        <div class="list-container">
          {#each blockedUsers as user}
            <div class="block-item tonal">
              <div class="user-info">
                <Avatar src={user.avatar_url} username={user.username} size={40} />
                <div class="details">
                  <strong>{user.full_name || user.username}</strong>
                  <span>@{user.username}</span>
                </div>
              </div>
              <Button variant="ghost" icon="person_remove" onclick={() => unblockUser(user.username)}>
                Desbloquear
              </Button>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .blocks-settings { display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header { margin-bottom: 2.5rem; }
  .title-with-icon { display: flex; align-items: center; gap: 1.5rem; }
  .hero-icon { font-size: 2.5rem; color: var(--error); background: rgba(var(--error-rgb), 0.1); padding: 1rem; border-radius: 1.5rem; }
  
  h2 { margin: 0; font-size: 1.75rem; font-weight: 900; }
  .subtitle { margin: 0.25rem 0 0 0; color: var(--text-muted); font-size: 1rem; }

  .add-block-form { padding: 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); margin-bottom: 2rem; }
  .form-row { display: flex; gap: 1rem; align-items: flex-end; }
  :global(.grow) { flex: 1; }

  .list-container { display: flex; flex-direction: column; gap: 0.75rem; }
  .block-item { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.5rem; border-radius: 1.25rem; border: 1px solid var(--md-sys-color-outline-variant); }
  
  .user-info { display: flex; align-items: center; gap: 1.25rem; }
  .details { display: flex; flex-direction: column; }
  .details strong { font-size: 1rem; }
  .details span { font-size: 0.8rem; color: var(--text-muted); }

  .empty-state { padding: 4rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }

  .skeleton-list { height: 150px; background: var(--surface-variant); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }
</style>
