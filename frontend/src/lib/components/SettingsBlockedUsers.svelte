<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  import Avatar from "./Avatar.svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import Input from "./ui/Input.svelte";
  

import UserMinus from "lucide-svelte/icons/user-minus";
import UserPlus from "lucide-svelte/icons/user-plus";
import ShieldAlert from "lucide-svelte/icons/shield-alert";



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

<div class="modular-stack">
  <!-- ADD BLOCK FORM -->
  <div class="config-module">
    <div class="module-tag">
      <UserPlus size={18} />
      Bloquear Usuario
    </div>
    <div class="form-box">
      <form onsubmit={(e) => { e.preventDefault(); blockUser(); }} class="form-content">
        <div class="form-row">
          <Input bind:value={newUsername} placeholder="Escribe un nombre de usuario..." class="grow" />
          <Button type="submit" variant="primary" disabled={submitting || !newUsername.trim()}>
            {submitting ? 'Bloqueando...' : 'Bloquear'}
          </Button>
        </div>
      </form>
    </div>
  </div>

  <!-- BLOCKED USERS LIST -->
  <div class="config-module">
    <div class="module-tag">
      <UserMinus size={18} />
      Usuarios Bloqueados
    </div>
    <div class="tonal-list">
      {#if loading}
        <div class="skeleton-list"></div>
      {:else if blockedUsers.length === 0}
        <div class="empty-state">
           <ShieldAlert size={48} />
           <p>No tienes a nadie bloqueado. ¡Genial!</p>
        </div>
      {:else}
        {#each blockedUsers as user}
          <div class="block-item">
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
      {/if}
    </div>
  </div>
</div>

<style>
  .modular-stack { display: flex; flex-direction: column; gap: 3rem; }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }

  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }

  .form-box { padding: 3rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); background: rgba(255,255,255,0.01); }
  .form-content { display: flex; flex-direction: column; gap: 2rem; }
  .form-row { display: flex; gap: 1rem; align-items: flex-end; }
  :global(.grow) { flex: 1; }

  .block-item { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  .block-item:last-child { border-bottom: none; }
  .block-item:hover { background: rgba(255,255,255,0.02); }

  .user-info { display: flex; align-items: center; gap: 1.25rem; }
  .details { display: flex; flex-direction: column; }
  .details strong { font-size: 1rem; }
  .details span { font-size: 0.8rem; color: var(--text-muted); }

  .empty-state { padding: 4rem 2rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }

  .skeleton-list { height: 150px; background: var(--surface-variant); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }
</style>
