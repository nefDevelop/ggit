<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  

import Users from "lucide-svelte/icons/users";
import Building2 from "lucide-svelte/icons/building-2";
import LogOut from "lucide-svelte/icons/log-out";


  import { dialog } from '$lib/dialog.svelte';
  import { toast } from '$lib/toast.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import Avatar from './Avatar.svelte';

  let orgs = $state<any[]>([]);
  let loading = $state(true);

  async function fetchOrgs() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/orgs");
      if (res.ok) orgs = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function leaveOrg(id: number, name: string) {
    if (!await dialog.confirm(`¿Estás seguro de que quieres abandonar la organización ${name}?`)) return;
    try {
      const res = await apiFetch(`/api/v2/user/settings/orgs/${id}`, { method: "DELETE" });
      if (res.ok) {
        await fetchOrgs();
      } else {
        const data = await res.json();
        toast.error(data.message || "Error al abandonar la organización");
      }
    } catch (e) { console.error(e); }
  }

  onMount(fetchOrgs);
</script>

<div class="settings-orgs-container">


  <div class="modular-stack">
    <div class="config-module">
      <div class="module-tag">
        <Users size={20} />
        Equipos y Entidades
      </div>
      <div class="island tonal-list">
        {#if loading}
          <div class="skeleton-orgs"></div>
        {:else if orgs.length === 0}
          <div class="empty-orgs-box">
             <Building2 size={64} />
             <p>No perteneces a ninguna organización todavía.</p>
          </div>
        {:else}
          {#each orgs as org}
            <div class="org-modular-row">
              <div class="org-brand">
                 <div class="org-avatar-shield">
                    <Avatar src={org.avatar_url} username={org.username} size={64} />
                 </div>
                 <div class="org-meta">
                    <a href="/org/{org.username}" class="org-full-name">{org.full_name || org.username}</a>
                    <span class="org-alias">@{org.username}</span>
                 </div>
              </div>
              <div class="org-actions">
                <Button 
                  variant="ghost" 
                  size="sm" 
                  icon={LogOut} 
                  class="danger-btn" 
                  onclick={() => leaveOrg(org.id, org.username)}
                >Abandonar</Button>
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  


  .modular-stack { display: flex; flex-direction: column; gap: 3rem; }
  @media (max-width: 600px) {
    .modular-stack { gap: 2rem; }
  }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }

  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }

  
  .org-modular-row { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  @media (max-width: 600px) { .org-modular-row { flex-direction: column; align-items: stretch; padding: 1rem 1.25rem; gap: 1rem; } }
  .org-modular-row:last-child { border-bottom: none; }
  .org-modular-row:hover { background: rgba(255,255,255,0.02); }

  .org-brand { display: flex; align-items: center; gap: 1.5rem; }
  .org-avatar-shield { width: 64px; height: 64px; border-radius: 1.25rem; overflow: hidden; border: 2px solid var(--md-sys-color-outline-variant); }
  @media (max-width: 600px) { .org-avatar-shield { width: 48px; height: 48px; border-radius: 1rem; } }
  
  .org-meta { display: flex; flex-direction: column; gap: 0.2rem; }
  .org-full-name { font-size: 1.25rem; font-weight: 900; color: var(--text-color); text-decoration: none; }
  .org-full-name:hover { color: var(--primary-color); }
  .org-alias { font-size: 0.9rem; color: var(--text-muted); font-weight: 700; opacity: 0.7; }

  .org-actions { flex-shrink: 0; }
  @media (max-width: 600px) {
    .org-actions { display: flex; justify-content: flex-end; }
  }

  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .empty-orgs-box { padding: 5rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1.5rem; }
  .empty-orgs-box :global(svg) { opacity: 0.2; }

  .skeleton-orgs { height: 240px; background: var(--surface-variant); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
</style>
