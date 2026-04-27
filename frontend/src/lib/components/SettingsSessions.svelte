<script lang="ts">
  import { onMount } from "svelte";
  import { apiFetch } from "$lib/api";
  import { i18n } from '$lib/i18n.svelte';
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

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">devices</span>
      <div>
        <h2>{i18n.t("settings.sessions")}</h2>
        <p class="subtitle">{i18n.t("settings.sessions_desc")}</p>
      </div>
    </div>
    {#if sessions.length > 1}
      <Button variant="ghost" class="danger-btn" icon="logout" onclick={revokeOtherSessions}>
        {i18n.t("settings.sessions_revoke_all")}
      </Button>
    {/if}
  </header>

  <div class="sessions-list">
    {#if loading}
      <div class="skeleton-list"></div>
    {:else if sessions.length === 0}
      <p class="empty-msg">{i18n.t("settings.sessions_empty")}</p>
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
                <span class="time">{i18n.t("settings.sessions_last_activity")}: {new Date(session.updated_at).toLocaleString()}</span>
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
</Card>

<style>
  .section-header { margin-bottom: 3rem; display: flex; justify-content: space-between; align-items: flex-start; gap: 2rem; }
  @media (max-width: 800px) { .section-header { flex-direction: column; align-items: stretch; } }

  .title-with-icon { display: flex; align-items: center; gap: 1.5rem; }
  .hero-icon { font-size: 2.5rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1rem; border-radius: 1.5rem; }
  h2 { margin: 0; font-size: 1.75rem; font-weight: 900; }
  .subtitle { margin: 0.25rem 0 0 0; color: var(--text-muted); font-size: 1rem; }

  .sessions-list { display: flex; flex-direction: column; gap: 1rem; }
  
  .session-item { display: flex; align-items: center; justify-content: space-between; padding: 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); background: rgba(var(--primary-rgb), 0.02); transition: all 0.2s; }
  .session-item.is-current { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  
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

  .skeleton-list { height: 200px; background: var(--surface-variant); border-radius: 1.5rem; animation: pulse 1.5s infinite; opacity: 0.3; }
  @keyframes pulse { 50% { opacity: 0.1; } }

  .empty-msg { text-align: center; padding: 3rem; color: var(--text-muted); font-style: italic; }

  @media (max-width: 600px) {
    .session-item { flex-direction: column; align-items: stretch; gap: 1rem; }
    .session-info { gap: 1rem; }
  }
</style>
