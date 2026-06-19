<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import Button from "./ui/Button.svelte";
  

import UserCog from "lucide-svelte/icons/user-cog";
import UserPlus from "lucide-svelte/icons/user-plus";
import AlertCircle from "lucide-svelte/icons/alert-circle";
import RefreshCw from "lucide-svelte/icons/refresh-cw";
import Save from "lucide-svelte/icons/save";



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
    if (submitting) return;
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

<div class="user-editor">
  <header class="editor-header">
    <div class="title-row">
      {#if userId}
        <UserCog size={56} class="icon-hero" />
      {:else}
        <UserPlus size={56} class="icon-hero" />
      {/if}
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
          <AlertCircle size={24} />
          {error}
        </div>
      {/if}

      <div class="form-actions">
        <Button variant="ghost" onclick={onCancel} disabled={submitting}>
          {i18n.t('cancel')}
        </Button>
        <Button type="submit" variant="primary" size="lg" disabled={submitting} loading={submitting} icon={Save}>
          {i18n.t('admin.save_user')}
        </Button>
      </div>
    </form>
  {/if}
</div>

<style>
  .user-editor {
    padding: 3.5rem;
    border-radius: 2.5rem;
    background: var(--md-sys-color-surface);
    border: 1px solid var(--md-sys-color-outline-variant);
    max-width: 1100px;
    margin: 0 auto;
    box-shadow: 0 32px 64px -12px rgba(0, 0, 0, 0.4);
    position: relative;
    overflow: hidden;
  }

  .user-editor::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 8px;
    background: linear-gradient(90deg, var(--primary-color), var(--secondary), var(--tertiary));
    opacity: 0.8;
  }

  .editor-header { 
    margin-bottom: 4rem;
    animation: slideDown 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  @keyframes slideDown {
    from { opacity: 0; transform: translateY(-20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .title-row { display: flex; align-items: center; gap: 2.5rem; }
  
  :global(.icon-hero) { 
    color: var(--primary-color); 
    background: rgba(var(--primary-rgb), 0.08); 
    padding: 1.5rem; 
    border-radius: 1.5rem;
    box-shadow: 0 8px 16px rgba(var(--primary-rgb), 0.15);
    border: 1px solid rgba(var(--primary-rgb), 0.1);
  }
  
  h3 { margin: 0; font-size: 2.75rem; font-weight: 950; letter-spacing: -0.05em; color: var(--text-color); }
  .subtitle { margin: 0.75rem 0 0 0; color: var(--text-muted); font-size: 1.15rem; font-weight: 500; }

  .form-sections { display: flex; flex-direction: column; gap: 4.5rem; }
  
  .form-section {
    animation: fadeInScale 0.7s cubic-bezier(0.2, 0.8, 0.2, 1) backwards;
  }
  
  .form-section:nth-child(1) { animation-delay: 0.1s; }
  .form-section:nth-child(2) { animation-delay: 0.2s; }
  .form-section:nth-child(3) { animation-delay: 0.3s; }

  @keyframes fadeInScale {
    from { opacity: 0; transform: scale(0.98) translateY(10px); }
    to { opacity: 1; transform: scale(1) translateY(0); }
  }

  .section-label { 
    font-size: 0.85rem; 
    font-weight: 900; 
    text-transform: uppercase; 
    color: var(--primary-color); 
    letter-spacing: 0.2em; 
    margin-bottom: 2.5rem; 
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }

  .section-label::after {
    content: '';
    height: 1px;
    flex: 1;
    background: linear-gradient(90deg, var(--md-sys-color-outline-variant), transparent);
  }

  .form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 2.5rem; }
  @media (max-width: 768px) { .form-grid { grid-template-columns: 1fr; } }

  .field { display: flex; flex-direction: column; gap: 1rem; }
  .field label { font-size: 0.95rem; font-weight: 800; color: var(--text-color); margin-left: 0.25rem; }
  
  .field input { 
    background: var(--md-sys-color-surface-container); 
    border: 2px solid transparent; 
    padding: 1.15rem 1.75rem; 
    border-radius: 1.25rem; 
    color: var(--text-color); 
    font-size: 1.1rem; 
    font-weight: 500;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); 
    box-shadow: inset 0 2px 4px rgba(0,0,0,0.05);
  }

  .field input:focus { 
    border-color: var(--primary-color); 
    outline: none; 
    background: var(--md-sys-color-surface-container-high); 
    transform: translateY(-2px);
    box-shadow: 0 12px 24px -8px rgba(var(--primary-rgb), 0.25);
  }

  .hint { margin: 0.5rem 0 0 0.5rem; font-size: 0.85rem; color: var(--text-muted); opacity: 0.8; font-style: italic; }

  .checkbox-group-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1.5rem; }

  .premium-checkbox {
    display: flex;
    align-items: flex-start;
    gap: 1.5rem;
    padding: 1.75rem;
    background: var(--md-sys-color-surface-container-low);
    border: 2px solid transparent;
    border-radius: 1.5rem;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
  }
  
  .premium-checkbox input { 
    position: absolute;
    opacity: 0;
  }

  .premium-checkbox::before {
    content: '';
    position: absolute;
    inset: 0;
    background: var(--primary-color);
    opacity: 0;
    transition: opacity 0.3s;
  }

  .premium-checkbox:hover { 
    background: var(--md-sys-color-surface-container); 
    transform: translateY(-4px) scale(1.02); 
    box-shadow: 0 12px 24px -4px rgba(0,0,0,0.1);
  }

  .premium-checkbox.checked { 
    border-color: var(--primary-color); 
    background: rgba(var(--primary-rgb), 0.05);
  }

  .premium-checkbox.checked .label { color: var(--primary-color); }

  .premium-checkbox.danger.checked { 
    border-color: var(--error); 
    background: rgba(var(--error-rgb), 0.05); 
  }
  .premium-checkbox.danger.checked .label { color: var(--error); }
  
  .check-info { position: relative; z-index: 1; display: flex; flex-direction: column; gap: 0.5rem; }
  .check-info .label { font-weight: 900; font-size: 1.15rem; color: var(--text-color); transition: color 0.3s; }
  .check-info .desc { font-size: 0.9rem; color: var(--text-muted); line-height: 1.5; font-weight: 500; }

  .error-banner { 
    background: var(--md-sys-color-error-container); 
    color: var(--md-sys-color-on-error-container); 
    padding: 1.5rem 2.5rem; 
    border-radius: 1.5rem; 
    display: flex; 
    align-items: center; 
    gap: 1.25rem; 
    font-weight: 800; 
    margin-top: 4rem; 
    border: 1px solid var(--md-sys-color-error);
    animation: shake 0.4s cubic-bezier(.36,.07,.19,.97) both;
  }

  @keyframes shake {
    10%, 90% { transform: translate3d(-1px, 0, 0); }
    20%, 80% { transform: translate3d(2px, 0, 0); }
    30%, 50%, 70% { transform: translate3d(-4px, 0, 0); }
    40%, 60% { transform: translate3d(4px, 0, 0); }
  }

  .form-actions { 
    display: flex; 
    justify-content: flex-end; 
    gap: 2rem; 
    margin-top: 5rem; 
    padding-top: 3.5rem; 
    border-top: 1px solid var(--md-sys-color-outline-variant); 
  }

  .loading-state { height: 600px; display: flex; align-items: center; justify-content: center; }
  .loading-spinner { width: 64px; height: 64px; border: 5px solid var(--md-sys-color-primary-container); border-top-color: var(--primary-color); border-radius: 50%; animation: spin 0.8s cubic-bezier(0.4, 0, 0.2, 1) infinite; }
  
  @keyframes spin { to { transform: rotate(360deg); } }
</style>
