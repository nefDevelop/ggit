<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';

  let {
    userId = null,
    onSave,
    onCancel,
  } = $props<{
    userId?: number | null;
    onSave: () => void;
    onCancel: () => void;
  }>();

  let user = $state<any>({
    username: "",
    email: "",
    password: "",
    full_name: "",
    website: "",
    location: "",
    max_repo_creation: -1,
    activated: true,
    admin: false,
    prohibit_login: false,
    allow_git_hook: false,
    allow_import_local: false,
  });

  let loading = $state(false);
  let submitting = $state(false);
  let error = $state("");

  async function fetchUser() {
    if (!userId) return;
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/users/${userId}`);
      if (res.ok) {
        const data = await res.json();
        user = {
          username: data.username,
          email: data.email,
          full_name: data.full_name || "",
          website: data.website || "",
          location: data.location || "",
          max_repo_creation: data.max_repo_creation ?? -1,
          activated: data.is_activated,
          admin: data.is_admin,
          prohibit_login: data.prohibit_login || false,
          allow_git_hook: data.allow_git_hook || false,
          allow_import_local: data.allow_import_local || false,
          password: "", 
        };
      }
    } catch (e) {
      error = "Failed to load user data.";
    } finally {
      loading = false;
    }
  }

  async function save() {
    submitting = true;
    error = "";
    const method = userId ? "PATCH" : "POST";
    const url = userId
      ? `/api/v2/admin/users/${userId}`
      : "/api/v2/admin/users";

    const payload = { ...user };
    if (userId && !payload.password) delete payload.password;
    if (userId) delete payload.username;

    try {
      const res = await apiFetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (res.ok) {
        onSave();
      } else {
        const data = await res.json();
        error = data.error || "Failed to save user.";
      }
    } catch (e) {
      error = "A network error occurred.";
    } finally {
      submitting = false;
    }
  }

  onMount(fetchUser);
</script>

<div class="user-editor tonal">
  <header class="editor-header">
    <div class="title-row">
      <span class="material-symbols-rounded icon-hero">{userId ? 'person_edit' : 'person_add'}</span>
      <div>
        <h3>{userId ? i18n.t('admin.edit_user') : i18n.t('admin.create_user')}</h3>
        <p class="subtitle">
          {userId
            ? i18n.t('admin.modifying_account', user.username)
            : i18n.t('admin.add_new_member')}
        </p>
      </div>
    </div>
  </header>

  {#if loading}
    <div class="loading-state">
      <div class="loading-spinner"></div>
    </div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); save(); }}>
      <div class="form-sections">
        <!-- Main Info -->
        <div class="form-section">
          <h4 class="section-label">Información Básica</h4>
          <div class="form-grid">
            {#if !userId}
              <div class="field">
                <label for="username">{i18n.t('username')}</label>
                <input id="username" type="text" bind:value={user.username} placeholder="johndoe" required />
              </div>
            {/if}
            <div class="field">
              <label for="email">{i18n.t('email')}</label>
              <input id="email" type="email" bind:value={user.email} placeholder="email@example.com" required />
            </div>
            <div class="field">
              <label for="password">
                {userId ? "Nueva Contraseña" : i18n.t('password')}
              </label>
              <input id="password" type="password" bind:value={user.password} required={!userId} placeholder="••••••••" />
              {#if userId}<p class="hint">Dejar en blanco para mantener la actual</p>{/if}
            </div>
            <div class="field">
              <label for="full_name">{i18n.t('settings.full_name')}</label>
              <input id="full_name" type="text" bind:value={user.full_name} placeholder="John Doe" />
            </div>
          </div>
        </div>

        <!-- Details -->
        <div class="form-section">
          <h4 class="section-label">Detalles y Límites</h4>
          <div class="form-grid">
            <div class="field">
              <label for="website">{i18n.t('website')}</label>
              <input id="website" type="url" bind:value={user.website} placeholder="https://..." />
            </div>
            <div class="field">
              <label for="location">{i18n.t('settings.location')}</label>
              <input id="location" type="text" bind:value={user.location} placeholder="City, Country" />
            </div>
            <div class="field">
              <label for="max_repo">{i18n.t('admin.max_repos')}</label>
              <input id="max_repo" type="number" bind:value={user.max_repo_creation} min="-1" />
              <p class="hint">{i18n.t('admin.max_repos_helper')}</p>
            </div>
          </div>
        </div>

        <!-- Permissions -->
        <div class="form-section">
          <h4 class="section-label">Estado y Permisos</h4>
          <div class="checkbox-group-grid">
            <label class="premium-checkbox" class:checked={user.activated}>
              <input type="checkbox" bind:checked={user.activated} />
              <div class="check-info">
                <span class="label">{i18n.t('admin.account_activated')}</span>
                <span class="desc">Permitir al usuario usar la plataforma.</span>
              </div>
            </label>

            <label class="premium-checkbox" class:checked={user.admin}>
              <input type="checkbox" bind:checked={user.admin} />
              <div class="check-info">
                <span class="label">{i18n.t('admin.admin_privileges')}</span>
                <span class="desc">Acceso total al panel de administración.</span>
              </div>
            </label>

            <label class="premium-checkbox danger" class:checked={user.prohibit_login}>
              <input type="checkbox" bind:checked={user.prohibit_login} />
              <div class="check-info">
                <span class="label">{i18n.t('admin.prohibit_login')}</span>
                <span class="desc">Bloquear el acceso a este usuario.</span>
              </div>
            </label>

            <label class="premium-checkbox" class:checked={user.allow_git_hook}>
              <input type="checkbox" bind:checked={user.allow_git_hook} />
              <div class="check-info">
                <span class="label">Permitir Git Hooks</span>
                <span class="desc">El usuario podrá gestionar scripts de servidor.</span>
              </div>
            </label>

            <label class="premium-checkbox" class:checked={user.allow_import_local}>
              <input type="checkbox" bind:checked={user.allow_import_local} />
              <div class="check-info">
                <span class="label">Importación Local</span>
                <span class="desc">Permitir migrar repositorios desde el disco local.</span>
              </div>
            </label>
          </div>
        </div>
      </div>

      {#if error}
        <div class="error-banner">
          <span class="material-symbols-rounded">error</span>
          {error}
        </div>
      {/if}

      <div class="form-actions">
        <button type="button" class="btn btn-ghost" onclick={onCancel}>
          {i18n.t('cancel')}
        </button>
        <button type="submit" class="btn btn-primary btn-lg" disabled={submitting}>
          <span class="material-symbols-rounded">{submitting ? 'sync' : 'save'}</span>
          {submitting ? i18n.t('settings.saving') : i18n.t('admin.save_user')}
        </button>
      </div>
    </form>
  {/if}
</div>

<style>
  .user-editor {
    padding: 3rem;
    border-radius: 3rem;
    background: rgba(var(--surface-rgb), 0.5);
    border: 1px solid var(--md-sys-color-outline-variant);
    backdrop-filter: blur(24px);
    max-width: 1000px;
    margin: 0 auto;
    box-shadow: var(--shadow-xl);
  }

  .editor-header { margin-bottom: 4rem; }
  .title-row { display: flex; align-items: center; gap: 2rem; }
  .icon-hero { font-size: 3.5rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  
  h3 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .form-sections { display: flex; flex-direction: column; gap: 4rem; }
  .section-label { font-size: 0.8rem; font-weight: 900; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; margin-bottom: 2rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; opacity: 0.8; }

  .form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 2rem; }
  @media (max-width: 768px) { .form-grid { grid-template-columns: 1fr; } }

  .field { display: flex; flex-direction: column; gap: 0.75rem; }
  .field label { font-size: 0.9rem; font-weight: 700; color: var(--text-color); margin-left: 0.5rem; }
  .field input { background: var(--surface-variant); border: 2px solid transparent; padding: 1rem 1.5rem; border-radius: 1.25rem; color: var(--text-color); font-size: 1.05rem; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
  .field input:focus { border-color: var(--primary-color); outline: none; background: rgba(var(--primary-rgb), 0.05); box-shadow: 0 0 0 5px rgba(var(--primary-rgb), 0.15); }
  .hint { margin: 0.25rem 0 0 0.5rem; font-size: 0.8rem; color: var(--text-muted); opacity: 0.7; font-style: italic; }

  .checkbox-group-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 1.25rem; }

  .premium-checkbox {
    display: flex;
    align-items: flex-start;
    gap: 1.25rem;
    padding: 1.5rem;
    background: var(--surface-variant);
    border: 2px solid transparent;
    border-radius: 1.75rem;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }
  .premium-checkbox input { display: none; }
  .premium-checkbox:hover { background: rgba(255,255,255,0.05); transform: translateY(-2px); }
  .premium-checkbox.checked { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.08); box-shadow: var(--shadow-md); }
  .premium-checkbox.danger.checked { border-color: var(--error); background: rgba(var(--error-rgb), 0.08); }
  
  .check-info { display: flex; flex-direction: column; gap: 0.4rem; }
  .check-info .label { font-weight: 800; font-size: 1.05rem; color: var(--text-color); }
  .check-info .desc { font-size: 0.85rem; color: var(--text-muted); line-height: 1.4; }

  .error-banner { background: rgba(var(--error-rgb), 0.15); color: var(--error); padding: 1.25rem 2rem; border-radius: 1.5rem; display: flex; align-items: center; gap: 1rem; font-weight: 700; margin-top: 3rem; border-left: 4px solid var(--error); }

  .form-actions { display: flex; justify-content: flex-end; gap: 1.5rem; margin-top: 4rem; padding-top: 2.5rem; border-top: 1px solid var(--md-sys-color-outline-variant); }
  
  .btn-lg { padding: 1.25rem 2.5rem; border-radius: 1.5rem; font-weight: 800; gap: 1rem; font-size: 1.1rem; }

  .loading-state { height: 500px; display: flex; align-items: center; justify-content: center; }
  .loading-spinner { width: 48px; height: 48px; border: 4px solid rgba(var(--primary-rgb), 0.1); border-top-color: var(--primary-color); border-radius: 50%; animation: spin 1s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
</style>
