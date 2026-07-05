<script lang="ts">
  import { onMount } from "svelte";
  import { apiFetch } from "$lib/api";
  import { i18n } from '$lib/i18n.svelte';
  import { formatDateTime } from "$lib/utils/date";
  import { dialog } from '$lib/dialog.svelte';
  import { toast } from '$lib/toast.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let sessions = $state<any[]>([]);
  let loading = $state(true);

  async function fetchSessions() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/sessions");
      if (res.ok) {
        sessions = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function revokeSession(id: number) {
    if (!await dialog.confirm(i18n.t("settings.sessions_revoke_confirm"))) return;
    
    try {
      const res = await apiFetch(`/api/v2/user/settings/sessions/${id}`, { method: "DELETE" });
      if (res.ok) {
        toast.success(i18n.t("settings.sessions_revoke_success"));
        fetchSessions();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function revokeOtherSessions() {
    if (!await dialog.confirm(i18n.t("settings.sessions_revoke_all_confirm"))) return;
    
    try {
      const res = await apiFetch("/api/v2/user/settings/sessions", { method: "DELETE" });
      if (res.ok) {
        toast.success(i18n.t("settings.sessions_revoke_all_success"));
        fetchSessions();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchSessions);
</script>

<div class="settings-sessions-container">
  {#if sessions.length > 1}
    <div class="revoke-all-bar">
      <Button variant="ghost" class="danger-btn" icon="logout" onclick={revokeOtherSessions}>
        {i18n.t("settings.sessions_revoke_all")}
      </Button>
    </div>
  {/if}

  <div class="modular-stack">
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">devices</span>
        Sesiones Activas
      </div>
      <div class="tonal-list">
        {#if loading}
          <div class="skeleton-list"></div>
        {:else if sessions.length === 0}
          <div class="empty-state">
            <p>{i18n.t("settings.sessions_empty")}</p>
          </div>
        {:else}
          {#each sessions as session}
            <div class="session-item" class:is-current={session.current}>
              <div class="session-info">
                <span class="material-symbols-rounded device-icon">
                  {session.user_agent.includes('Mobi') ? 'smartphone' : 'desktop_windows'}
                </span>
                <div class="details">
                  <div class="ua-line">
                    <strong>{session.user_agent}</strong>
                    {#if session.current}
                      <span class="current-badge">{i18n.t("settings.sessions_current")}</span>
                    {/if}
                  </div>
                  <div class="meta-line">
                    <span class="ip">{session.remote_addr}</span>
                    <span class="dot">•</span>
                    <span class="time">{i18n.t("settings.sessions_last_activity")}: {formatDateTime(session.updated_at)}</span>
                  </div>
                </div>
              </div>
              {#if !session.current}
                <Button variant="ghost" icon="cancel" onclick={() => revokeSession(session.id)}>{i18n.t("settings.sessions_revoke")}</Button>
              {/if}
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</div>

<style>
  .revoke-all-bar { display: flex; justify-content: flex-end; margin-bottom: 2rem; }

  .modular-stack { display: flex; flex-direction: column; gap: 3rem; }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }

  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }

  .tonal-list { display: flex; flex-direction: column; background: rgba(0,0,0,0.15); border-radius: 1.5rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }

  .session-item { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  .session-item:last-child { border-bottom: none; }
  .session-item:hover { background: rgba(255,255,255,0.02); }
  .session-item.is-current { background: rgba(var(--primary-rgb), 0.05); }

  .session-info { display: flex; align-items: center; gap: 1.5rem; flex: 1; min-width: 0; }
  .device-icon { font-size: 2rem; color: var(--text-muted); opacity: 0.7; }
  .is-current .device-icon { color: var(--primary-color); opacity: 1; }

  .details { display: flex; flex-direction: column; gap: 0.25rem; min-width: 0; }
  .ua-line { display: flex; align-items: center; gap: 0.75rem; }
  .ua-line strong { font-size: 1.05rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  
  .current-badge { font-size: 0.7rem; font-weight: 900; background: var(--primary-color); color: white; padding: 0.2rem 0.6rem; border-radius: 2rem; text-transform: uppercase; letter-spacing: 0.05em; }

  .meta-line { display: flex; align-items: center; gap: 0.5rem; color: var(--text-muted); font-size: 0.9rem; }
  .dot { opacity: 0.5; }

  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .empty-state { padding: 4rem 2rem; text-align: center; color: var(--text-muted); }

  .skeleton-list { height: 200px; background: var(--surface-variant); border-radius: 1.5rem; animation: pulse 1.5s infinite; opacity: 0.3; }
  @keyframes pulse { 50% { opacity: 0.1; } }

  @media (max-width: 600px) {
    .session-item { flex-direction: column; align-items: stretch; gap: 1rem; padding: 1.5rem; }
    .session-info { gap: 1rem; }
  }
</style>
