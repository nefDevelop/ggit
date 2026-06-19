<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  import { RepositorySchema } from "$lib/schemas";
  import Checkbox from "./ui/Checkbox.svelte";
  import Input from "./ui/Input.svelte";
  import Avatar from "./Avatar.svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  import { i18n } from "$lib/i18n.svelte";
  
import Settings from "lucide-svelte/icons/settings";
import RefreshCw from "lucide-svelte/icons/refresh-cw";
import Upload from "lucide-svelte/icons/upload";
import Trash2 from "lucide-svelte/icons/trash-2";
import FileText from "lucide-svelte/icons/file-text";
import Bug from "lucide-svelte/icons/bug";
import GitPullRequest from "lucide-svelte/icons/git-pull-request";
import AlertTriangle from "lucide-svelte/icons/alert-triangle";


  const repo: any = getContext("repo");
  let settings = $state<any>(null);
  let loading = $state(true);
  let saving = $state(false);
  let deleteConfirm = $state("");
  let avatarInput = $state<HTMLInputElement>();
  let avatarUploading = $state(false);
  let avatarCacheBuster = $state(Date.now());

  async function fetchSettings() {
    if (!repo?.data?.owner?.username || !repo?.data?.name) return;
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`);
      if (res.ok) {
        settings = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function saveSettings() {
    if (!settings) return;
    saving = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`, {
        method: "PATCH",
        body: JSON.stringify(settings),
      });
      if (res.ok) {
        toast.success(i18n.t("settings.saved_success") || "Ajustes guardados correctamente");
        if (settings.name !== repo.data.name) {
          goto(`/${repo.data.owner.username}/${settings.name}/settings`);
        }
      }
    } catch (e) {
      console.error(e);
      toast.error(i18n.t("settings.saved_error") || "Error al guardar los ajustes");
    } finally {
      saving = false;
    }
  }

  async function uploadAvatar(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file || !settings) return;
    avatarUploading = true;
    try {
      const formData = new FormData();
      formData.append("avatar", file);
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/avatar`, {
        method: "POST",
        body: formData,
      });
      if (res.ok) {
        avatarCacheBuster = Date.now();
        fetchSettings();
        toast.success("Avatar actualizado");
      }
    } catch (e) {
      console.error(e);
    } finally {
      avatarUploading = false;
    }
  }

  async function deleteAvatar() {
    if (!settings) return;
    if (!await dialog.confirm("¿Estás seguro de eliminar el logo del repositorio?")) return;
    avatarUploading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/avatar`, {
        method: "DELETE",
      });
      if (res.ok) {
        fetchSettings();
        toast.success("Avatar eliminado");
      }
    } catch (e) {
      console.error(e);
    } finally {
      avatarUploading = false;
    }
  }

  async function deleteRepo() {
    if (!settings) return;
    if (deleteConfirm !== settings.name) return;
    if (!await dialog.confirm(`¿ESTÁS SEGURO? Esta acción es IRREVERSIBLE. Se eliminará el repositorio "${settings.name}" permanentemente.`)) return;
    saving = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`, {
        method: "DELETE",
      });
      if (res.ok) {
        goto("/dashboard");
      }
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
    }
  }

  async function transferRepo() {
    if (!settings) return;
    const newOwner = await dialog.prompt("Introduce el nombre del nuevo propietario (usuario u organización):");
    if (!newOwner) return;
    saving = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/transfer`, {
        method: "POST",
        body: JSON.stringify({ new_owner: newOwner }),
      });
      if (res.ok) {
        toast.success("Transferencia iniciada.");
        goto("/dashboard");
      }
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
    }
  }

  // Use an effect to watch for repo changes so it fetches settings
  $effect(() => {
    if (repo?.data) {
        fetchSettings();
    }
  });

</script>

<div class="general-settings-wrapper">
  <Card padding="none" class="settings-premium-card">
    <PageHeader 
      title="Ajustes del Repositorio"
      subtitle="Configuración general, visibilidad y metadatos básicos de tu proyecto."
      icon={Settings}
    />

    <div class="settings-content">
      {#if loading || !settings}
        <div class="skeleton-form"></div>
      {:else}
        <form class="settings-form" onsubmit={(e) => { e.preventDefault(); saveSettings(); }}>
          <div class="avatar-settings-section island">
            <div class="avatar-preview-wrap">
              <Avatar 
                src={settings.avatar_url ? `${settings.avatar_url}${settings.avatar_url.includes('?') ? '&' : '?'}t=${avatarCacheBuster}` : ''} 
                username={settings.name} 
                size={110} 
                type="repo"
                class="repo-avatar-img"
              />
              {#if avatarUploading}
                <div class="avatar-overlay">
                  <RefreshCw class="spin" />
                </div>
              {/if}
            </div>
            <div class="avatar-actions">
              <strong>Logo del Repositorio</strong>
              <p>Personaliza la identidad visual de tu proyecto con una imagen cuadrada.</p>
              <div class="action-buttons">
                <Button variant="tonal" size="sm" icon={Upload} onclick={() => avatarInput?.click()} loading={avatarUploading}>Subir Nueva</Button>
                {#if settings.avatar_url}
                  <Button variant="ghost" size="sm" icon={Trash2} class="danger-text" onclick={deleteAvatar} loading={avatarUploading}>Eliminar</Button>
                {/if}
              </div>
              <input type="file" bind:this={avatarInput} accept="image/*" style="display: none;" onchange={uploadAvatar} />
            </div>
          </div>

          <div class="form-grid">
            <Input label="Nombre del Repositorio" bind:value={settings.name} placeholder="mi-gran-proyecto" required />
            <Input label="Propietario" value={settings.owner?.username || repo.data?.owner?.username || ''} disabled />
            <Input label="Descripción" bind:value={settings.description} placeholder="Una breve descripción sobre qué hace este proyecto..." class="full-width" />
            <Input label="Página Web" bind:value={settings.website} placeholder="https://mi-proyecto.com" class="full-width" />
          </div>

          <div class="options-grid">
            <Checkbox 
              label="Hacer Repositorio Privado" 
              description="Solo tú y las personas que autorices explícitamente podrán ver el contenido de este repositorio."
              bind:checked={settings.private} 
            />
            
            <div class="select-field">
              <label for="visibility">Visibilidad por defecto</label>
              <div class="select-wrapper">
                <select id="visibility">
                  <option value="public">Público</option>
                  <option value="private">Privado</option>
                </select>
              </div>
              <p class="select-helper-text">Define cómo aparecerá el repositorio en el buscador.</p>
            </div>
          </div>

          <div class="advanced-stack island">
            <header class="adv-header">
              <Settings size={20} />
              <strong>Módulos Activos</strong>
            </header>
            
            <div class="options-grid">
              <Checkbox label="Wiki" description="Documentación colaborativa." bind:checked={settings.enable_wiki} />
              <Checkbox label="Pull Requests" description="Aceptar contribuciones externas." bind:checked={settings.enable_pull_requests} />
              <Checkbox label="Issues" description="Seguimiento de tareas y errores." bind:checked={settings.enable_issues} />
              <Checkbox label="External Wiki" description="Usar una wiki alojada en otro lugar." bind:checked={settings.enable_external_wiki} />
            </div>

            {#if settings.enable_external_wiki}
              <div class="nested-box tonal">
                <Input label="URL de Wiki Externa" bind:value={settings.external_wiki_url} placeholder="https://docs.mi-proyecto.com" />
              </div>
            {/if}
          </div>

          <div class="advanced-stack island" class:archived-item={settings.archived}>
            <header class="adv-header">
              <FileText size={20} />
              <strong>Estado del Repositorio</strong>
            </header>
            <Checkbox 
              label="Archivar Repositorio" 
              description="El repositorio pasará a modo solo lectura. Útil para proyectos terminados."
              bind:checked={settings.archived} 
            />
          </div>

          <div class="form-actions-premium">
            <Button type="submit" variant="primary" loading={saving} disabled={saving}>
              Guardar Cambios Generales
            </Button>
          </div>
        </form>
      {/if}
    </div>
  </Card>

  {#if settings}
  <div class="danger-zone-container">
    <Card class="tonal" padding="none">
      <header class="danger-header">
        <AlertTriangle size={32} />
        <h3>Zona de Peligro</h3>
      </header>
      
      <div class="danger-content">
        <div class="danger-item">
          <div class="danger-info">
            <strong>Eliminar este repositorio</strong>
            <p>Una vez eliminado, no hay vuelta atrás. Por favor, asegúrate de lo que haces.</p>
          </div>
          <div class="danger-action-row">
            <Input bind:value={deleteConfirm} placeholder="Escribe el nombre del repositorio para confirmar" />
            <Button
              variant="primary"
              class="delete-btn"
              disabled={deleteConfirm !== settings.name}
              onclick={deleteRepo}
            >
              Eliminar Repositorio
            </Button>
          </div>
        </div>

        <div class="danger-item">
          <div class="danger-info">
            <strong>Transferir Propiedad</strong>
            <p>Mueve este repositorio a otro usuario o a una organización de la que seas miembro.</p>
          </div>
          <Button variant="outline" class="danger-outline" onclick={transferRepo} disabled={saving}>
            Transferir repositorio
          </Button>
        </div>
      </div>
    </Card>
  </div>
  {/if}
</div>

<style>
  .settings-content { padding: 2.5rem; }
  @media (max-width: 1300px) { .settings-content { padding: 1.5rem; } }
  @media (max-width: 600px) { .settings-content { padding: 1.5rem; } }

  .settings-form { display: flex; flex-direction: column; gap: 3rem; }

  .avatar-settings-section { display: flex; align-items: center; gap: 3rem; padding: 2rem; border-radius: 2rem; }
  .avatar-preview-wrap { 
    position: relative; 
    border-radius: 1rem; 
    overflow: hidden; 
    background: var(--md-sys-color-surface-container-highest); 
    border: 2px solid var(--md-sys-color-outline-variant); 
    flex-shrink: 0; 
    box-shadow: var(--shadow-md);
    width: 110px;
    height: 110px;
  }
  .avatar-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; color: white; }
  :global(.spin) { animation: spin 1s linear infinite; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

  .avatar-actions { display: flex; flex-direction: column; gap: 1rem; }
  .avatar-actions strong { font-size: 1.25rem; font-weight: 900; color: var(--text-color); }
  .avatar-actions p { margin: 0; color: var(--text-muted); font-size: 0.95rem; line-height: 1.4; }
  .action-buttons { display: flex; gap: 0.75rem; align-items: center; }

  @media (max-width: 1300px) {
    .avatar-settings-section { gap: 1.5rem; padding: 1.5rem; border-radius: 1.5rem; }
    .avatar-preview-wrap { border-radius: 0.75rem; width: 90px; height: 90px; }
    :global(.repo-avatar-img) { width: 100% !important; height: 100% !important; }
  }

  @media (max-width: 650px) {
    .avatar-settings-section { flex-direction: column; text-align: center; padding: 1.5rem; }
    .action-buttons { justify-content: center; }
  }

  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  :global(.full-width) { grid-column: span 2; }

  .select-field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    width: 100%;
  }
  .select-field label {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface-variant);
    margin-left: 0.25rem;
  }
  .select-wrapper {
    position: relative;
    width: 100%;
    display: flex;
    align-items: center;
  }
  .select-wrapper select {
    width: 100%;
    background: var(--md-sys-color-surface-container-high);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface);
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    font-size: 1rem;
    font-family: inherit;
    transition: all 0.2s;
    outline: none;
    box-sizing: border-box;
    cursor: pointer;
  }
  .select-wrapper select:focus {
    border-color: var(--md-sys-color-primary);
    background: var(--md-sys-color-surface-container-highest);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary), transparent 85%);
  }
  .select-helper-text {
    font-size: 0.8rem;
    color: var(--text-muted);
    font-weight: 500;
    margin-left: 0.25rem;
    opacity: 0.8;
  }

  @media (max-width: 1300px) {
    .form-grid { gap: 1rem; }
    .options-grid { grid-template-columns: 1fr; }
  }

  @media (max-width: 800px) {
    .form-grid { grid-template-columns: 1fr; }
    :global(.full-width) { grid-column: span 1; }
  }

  .options-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1.25rem; }
  :global(.archived-item) { border: 2px solid rgba(var(--primary-rgb), 0.2) !important; background: rgba(var(--primary-rgb), 0.05) !important; }

  .advanced-stack { display: flex; flex-direction: column; gap: 1.25rem; }
  .advanced-section { padding: 1.75rem; border-radius: 1.75rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .adv-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; color: var(--primary-color); }
  .adv-header strong { font-weight: 950; font-size: 0.8rem; text-transform: uppercase; letter-spacing: 0.1em; }
  
  .nested-box { margin-top: 1.25rem; padding-left: 2rem; border-left: 2px solid var(--md-sys-color-outline-variant); display: flex; flex-direction: column; gap: 1.25rem; animation: fadeIn 0.3s ease-out; }

  .form-actions-premium { display: flex; justify-content: flex-end; padding-top: 2rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  /* DANGER ZONE */
  .danger-zone-container { }
  .danger-zone-container :global(.card) { border: 2px solid var(--md-sys-color-error) !important; }
  .danger-header { padding: 2rem 2.5rem; border-bottom: 1px solid rgba(var(--error-rgb), 0.1); display: flex; align-items: center; gap: 1rem; color: var(--error); }
  .danger-header h3 { margin: 0; font-size: 1.5rem; font-weight: 900; }

  .danger-content { padding: 2.5rem; display: flex; flex-direction: column; gap: 2.5rem; }
  .danger-item { display: flex; flex-direction: column; gap: 1.5rem; }
  .danger-info strong { color: var(--error); font-size: 1.25rem; font-weight: 900; }
  .danger-info p { color: var(--text-muted); margin: 0.5rem 0 0; font-size: 1rem; }

  .danger-action-row { display: flex; gap: 1rem; align-items: flex-end; }
  .danger-action-row :global(.input-container) { flex: 1; }
  :global(.delete-btn) { background: var(--error) !important; color: var(--on-error) !important; }

  .danger-item:not(:first-child) { padding-top: 2.5rem; border-top: 1px solid rgba(var(--error-rgb), 0.1); }

  @media (max-width: 800px) {
    .danger-action-row { flex-direction: column; align-items: stretch; }
    .danger-header { padding: 1.5rem; }
    .danger-content { padding: 1.5rem; }
  }

  .skeleton-form { height: 400px; background: var(--md-sys-color-surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
</style>
