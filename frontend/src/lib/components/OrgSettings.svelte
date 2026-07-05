<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  

import Building2 from "lucide-svelte/icons/building-2";
import RefreshCw from "lucide-svelte/icons/refresh-cw";
import Upload from "lucide-svelte/icons/upload";
import Trash2 from "lucide-svelte/icons/trash-2";
import HeartCrack from "lucide-svelte/icons/heart-crack";


  import { i18n } from "$lib/i18n.svelte";
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";
  import Input from "./ui/Input.svelte";
  import Checkbox from "./ui/Checkbox.svelte";
  import Avatar from "./Avatar.svelte";
  import PageHeader from "./ui/PageHeader.svelte";

  let { orgName } = $props();

  let org = $state<any>({});
  let loading = $state(true);
  let saving = $state(false);
  let avatarUploading = $state(false);

  async function fetchOrg() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}`);
      if (res.ok) org = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function saveOrg() {
    saving = true;
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(org),
      });
      if (res.ok) {
        toast.success("Organization updated successfully!");
      }
    } catch (e) {
      console.error(e);
      toast.error("Failed to update organization.");
    } finally {
      saving = false;
    }
  }

  async function updateAvatar(e: Event) {
    const target = e.target as HTMLInputElement;
    if (!target.files?.length) return;

    avatarUploading = true;
    const formData = new FormData();
    formData.append("avatar", target.files[0]);

    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}/avatar`, {
        method: "POST",
        body: formData,
      });
      if (res.ok) {
        toast.success("Avatar updated successfully!");
        fetchOrg();
      }
    } catch (e) {
      console.error(e);
    } finally {
      avatarUploading = false;
    }
  }

  async function deleteAvatar() {
    if (!await dialog.confirm("¿Seguro que quieres eliminar el avatar de la organización?")) return;
    
    avatarUploading = true;
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}/avatar`, { method: "DELETE" });
      if (res.ok) {
        toast.success("Avatar deleted.");
        fetchOrg();
      }
    } catch (e) {
      console.error(e);
    } finally {
      avatarUploading = false;
    }
  }

  async function deleteOrg() {
    if (!await dialog.confirm(i18n.t("settings.org_delete_confirm") || "Are you absolutely sure you want to delete this organization? This action is permanent.")) return;
    
    saving = true;
    try {
      const res = await apiFetch(`/api/v2/orgs/${orgName}`, { method: "DELETE" });
      if (res.ok) {
        window.location.href = "/";
      } else {
        toast.error("Failed to delete organization.");
      }
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
    }
  }

  onMount(fetchOrg);
</script>

<div class="org-settings-wrapper">
  <Card padding="none" class="settings-premium-card">
    <PageHeader
      title="Ajustes de la Organización"
      subtitle="Gestiona la identidad y visibilidad de tu organización"
      icon={Building2}
    />

    {#if loading}
      <div class="skeleton-form" style="padding: 2.5rem;"></div>
    {:else}
      <form class="settings-form" onsubmit={(e) => { e.preventDefault(); saveOrg(); }}>
        
        <div class="avatar-settings-section">
          <div class="avatar-preview-wrap">
            <Avatar 
              src={org.avatar_url} 
              username={org.username} 
              seed={org.id?.toString()} 
              size={120} 
            />
            {#if avatarUploading}
              <div class="avatar-overlay">
                <RefreshCw size={24} class="spin" />
              </div>
            {/if}
          </div>
          <div class="avatar-actions">
            <header>
              <h4>Avatar de la Organización</h4>
              <p>Sube una imagen personalizada para tu organización.</p>
            </header>
            <div class="action-buttons">
              <Button variant="outline" icon={Upload} onclick={() => document.getElementById('org-avatar-input')?.click()}>
                Subir Imagen
              </Button>
              {#if org.use_custom_avatar}
                <Button variant="ghost" icon={Trash2} class="danger-text" onclick={deleteAvatar}>
                  Eliminar
                </Button>
              {/if}
              <input 
                type="file" 
                id="org-avatar-input" 
                accept="image/*" 
                hidden 
                onchange={updateAvatar} 
              />
            </div>
          </div>
        </div>

        <div class="form-grid">
          <Input label="Nombre Público" bind:value={org.full_name} placeholder="Ggit Core Team" />
          <div class="field">
            <label for="desc">Descripción</label>
            <textarea id="desc" bind:value={org.description} placeholder="Una breve descripción..."></textarea>
          </div>
          <Input label="Sitio Web" type="url" bind:value={org.website} placeholder="https://ggit.io" />
          <Input label="Ubicación" bind:value={org.location} placeholder="Madrid, España" />
          
          <div class="field full-width">
            <Checkbox 
              label="Organización Privada" 
              description="Las organizaciones privadas solo son visibles para sus miembros." 
              bind:checked={org.private} 
            />
          </div>
        </div>

        <div class="form-footer">
          <Button type="submit" variant="primary" disabled={saving}>
            {saving ? 'Guardando...' : 'Actualizar Ajustes'}
          </Button>
        </div>
      </form>
    {/if}
  </Card>

  <div style="margin-top: 3rem;">
  <Card padding="none" class="settings-premium-card danger-card">
    <PageHeader
      title="Zona de Peligro"
      subtitle="Acciones irreversibles sobre la organización"
      icon={HeartCrack}
    />
    
    <div class="danger-zone-content">
      <div class="danger-item">
        <div class="danger-info">
          <strong>Eliminar esta organización</strong>
          <p>Esta acción eliminará permanentemente la organización, todos sus repositorios y equipos asociados. No se puede deshacer.</p>
        </div>
        <Button variant="primary" class="danger-btn" onclick={deleteOrg} disabled={saving}>
          {saving ? 'Eliminando...' : 'Eliminar Organización'}
        </Button>
      </div>
    </div>
  </Card>
  </div>
</div>

<style>
  .settings-form { padding: 2.5rem; display: flex; flex-direction: column; gap: 3rem; }
  @media (max-width: 600px) {
    .settings-form { padding: 1.5rem; gap: 2rem; }
  }
  
  .avatar-settings-section { display: flex; align-items: center; gap: 3rem; padding-bottom: 3rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); margin-bottom: 1rem; }
  .avatar-preview-wrap { position: relative; border-radius: 50%; overflow: hidden; background: var(--md-sys-color-surface-container-highest); border: 2px solid var(--md-sys-color-outline-variant); flex-shrink: 0; }
  .avatar-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; color: white; }
  :global(.spin) { animation: spin 1s linear infinite; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

  .avatar-actions { display: flex; flex-direction: column; gap: 1.5rem; }
  .avatar-actions h4 { margin: 0; font-size: 1.25rem; font-weight: 900; }
  .avatar-actions p { margin: 0; color: var(--text-muted); font-size: 1rem; }
  .action-buttons { display: flex; gap: 1rem; align-items: center; }

  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; }
  .field { display: flex; flex-direction: column; gap: 0.75rem; }
  .field label { font-size: 0.9rem; font-weight: 700; color: var(--text-color); margin-left: 0.5rem; }
  .field textarea { background: var(--md-sys-color-surface-container-high); border: 1px solid var(--md-sys-color-outline-variant); padding: 0.75rem 1rem; border-radius: 0.75rem; color: var(--text-color); font-size: 1rem; min-height: 100px; font-family: inherit; transition: all 0.2s; }
  .field textarea:focus { border-color: var(--md-sys-color-primary); outline: none; background: var(--md-sys-color-surface-container-highest); box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary), transparent 85%); }

  .full-width { grid-column: span 2; }

  .form-footer { display: flex; justify-content: flex-end; padding-top: 2rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  :global(.danger-card) { border: 2px solid var(--md-sys-color-error) !important; }
  .danger-zone-content { padding: 2.5rem; }
  .danger-item { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  .danger-info strong { font-size: 1.1rem; color: var(--error); display: block; margin-bottom: 0.25rem; }
  .danger-info p { color: var(--text-muted); font-size: 0.95rem; margin: 0; }

  :global(.danger-btn) { background: var(--error) !important; color: white !important; }
  :global(.danger-text) { color: var(--error) !important; }

  .skeleton-form { height: 400px; background: var(--surface-variant); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.6; } }

  @media (max-width: 800px) {
    .form-grid { grid-template-columns: 1fr; }
    .full-width { grid-column: auto; }
    .avatar-settings-section { flex-direction: column; text-align: center; }
    .danger-item { flex-direction: column; text-align: center; }
  }
</style>
