<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { goto } from "$app/navigation";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  import { RepositorySchema } from "$lib/schemas";
  import Checkbox from "./ui/Checkbox.svelte";
  import Input from "./ui/Input.svelte";
  import Avatar from "./Avatar.svelte";

  const repo: any = getContext("repo");
  let settings = $state<any>({});
  let loading = $state(true);
  let saving = $state(false);
  let avatarUploading = $state(false);
  let deleteConfirm = $state("");
  let errors = $state<Record<string, string>>({});

  async function fetchSettings() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`,
      );
      if (res.ok) settings = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function updateAvatar(e: Event) {
    const target = e.target as HTMLInputElement;
    if (!target.files?.length) return;

    avatarUploading = true;
    const formData = new FormData();
    formData.append("avatar", target.files[0]);

    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/avatar`,
        {
          method: "POST",
          body: formData,
        },
      );
      if (res.ok) {
        toast.success("Avatar updated successfully!");
        fetchSettings();
      }
    } catch (e) {
      console.error(e);
    } finally {
      avatarUploading = false;
    }
  }

  async function deleteAvatar() {
    if (!await dialog.confirm("¿Seguro que quieres eliminar el avatar del repositorio?")) return;
    
    avatarUploading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/avatar`,
        { method: "DELETE" },
      );
      if (res.ok) {
        toast.success("Avatar deleted.");
        fetchSettings();
      }
    } catch (e) {
      console.error(e);
    } finally {
      avatarUploading = false;
    }
  }

  async function saveSettings() {
    errors = {};
    
    // Validación con Zod
    const result = RepositorySchema.safeParse(settings);
    if (!result.success) {
      result.error.issues.forEach(err => {
        errors[err.path[0]] = err.message;
      });
      toast.error("Por favor, corrige los errores en el formulario.");
      return;
    }

    saving = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`,
        {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(result.data),
        },
      );
      if (res.ok) {
        toast.success("Repository updated successfully!");
        fetchSettings();
      }
    } catch (e) {
      // toast.error ya se gestiona en apiFetch
    } finally {
      saving = false;
    }
  }

  async function deleteRepo() {
    if (deleteConfirm !== settings.name) return;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`,
        {
          method: "DELETE",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name: deleteConfirm }),
        },
      );
      if (res.ok) {
        toast.success("Repository deleted.");
        goto("/dashboard");
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function transferRepo() {
    const newOwner = await dialog.prompt("Introduce el nombre de usuario o organización del nuevo propietario:");
    if (!newOwner) return;
    
    saving = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/transfer`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ new_owner: newOwner }),
        },
      );
      if (res.ok) {
        toast.success("Repository transfer initiated!");
        goto("/dashboard");
      }
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
    }
  }

  onMount(fetchSettings);
</script>

<div class="settings-card">
  {#if loading}
    <div class="skeleton-form"></div>
  {:else}
    <form class="settings-form" onsubmit={(e) => { e.preventDefault(); saveSettings(); }}>
      <div class="avatar-settings-section">
        <div class="avatar-preview-wrap">
          <Avatar 
            src={settings.avatar_url} 
            username={settings.name} 
            seed={settings.id?.toString()} 
            size={120} 
          />
          {#if avatarUploading}
            <div class="avatar-overlay">
              <span class="material-symbols-rounded spin">sync</span>
            </div>
          {/if}
        </div>
        <div class="avatar-actions">
          <header>
            <h4>Avatar del Repositorio</h4>
            <p>Sube una imagen personalizada para este repositorio.</p>
          </header>
          <div class="action-buttons">
            <Button variant="tonal" icon="upload" onclick={() => document.getElementById('repo-avatar-input')?.click()}>
              Subir Imagen
            </Button>
            {#if settings.use_custom_avatar}
              <Button variant="ghost" icon="delete" class="danger-text" onclick={deleteAvatar}>
                Eliminar
              </Button>
            {/if}
            <input 
              type="file" 
              id="repo-avatar-input" 
              accept="image/*" 
              hidden 
              onchange={updateAvatar} 
            />
          </div>
        </div>
      </div>

      <div class="form-grid">
        <Input 
          label="Repository Name" 
          bind:value={settings.name} 
          error={errors.name}
          required 
        />
        <Input 
          label="Website" 
          bind:value={settings.website} 
          error={errors.website}
          placeholder="https://example.com" 
        />
        
        <Input 
          label="Description" 
          bind:value={settings.description} 
          error={errors.description}
          placeholder="Short description..."
          multiline
          class="full-width"
        />

        <Input 
          label="Wiki Folder" 
          bind:value={settings.wiki_folder_name} 
          error={errors.wiki_folder_name}
          placeholder="docs" 
        />
      </div>

      <div class="options-grid">
        <Checkbox 
          label="Make Private" 
          description="Only visible to collaborators." 
          bind:checked={settings.private} 
        />
        <Checkbox 
          label="Unlisted" 
          description="No aparecerá en la exploración pública." 
          bind:checked={settings.unlisted} 
        />
        <Checkbox 
          label="Archivar Repositorio" 
          description="Convertir en solo lectura. Útil para proyectos terminados." 
          bind:checked={settings.archived}
          class="archived-item"
        />
      </div>

      <div class="advanced-section">
        <header class="section-header-inline">
          <span class="material-symbols-rounded">description</span>
          <h4>Configuración de Wiki</h4>
        </header>
        <div class="section-content">
          <Checkbox 
            label="Habilitar Wiki" 
            description="Activar documentación del proyecto." 
            bind:checked={settings.enable_wiki} 
          />
          {#if settings.enable_wiki}
            <div class="nested-options">
              <Input 
                label="Nombre de Carpeta Wiki" 
                bind:value={settings.wiki_folder_name} 
                placeholder="docs" 
              />
              <Checkbox 
                label="Permitir acceso público" 
                description="Cualquiera puede ver la wiki, incluso si el repositorio es privado." 
                bind:checked={settings.allow_public_wiki} 
              />
              <Checkbox 
                label="Wiki Externa" 
                description="Redirigir a una URL de documentación externa." 
                bind:checked={settings.enable_external_wiki} 
              />
              {#if settings.enable_external_wiki}
                <Input 
                  label="URL de Wiki Externa" 
                  bind:value={settings.external_wiki_url} 
                  placeholder="https://docs.example.com" 
                />
              {/if}
            </div>
          {/if}
        </div>
      </div>

      <div class="advanced-section">
        <header class="section-header-inline">
          <span class="material-symbols-rounded">bug_report</span>
          <h4>Gestión de Issues</h4>
        </header>
        <div class="section-content">
          <Checkbox 
            label="Habilitar Issues" 
            description="Seguimiento de tareas integrado." 
            bind:checked={settings.enable_issues} 
          />
          {#if settings.enable_issues}
            <div class="nested-options">
              <Checkbox 
                label="Permitir acceso público" 
                description="Cualquiera puede ver los issues, incluso si el repositorio es privado." 
                bind:checked={settings.allow_public_issues} 
              />
              <Checkbox 
                label="Tracker Externo" 
                description="Usar un sistema de seguimiento de errores externo (Jira, GitHub, etc)." 
                bind:checked={settings.enable_external_tracker} 
              />
              {#if settings.enable_external_tracker}
                <Input 
                  label="URL de Tracker Externo" 
                  bind:value={settings.external_tracker_url} 
                  placeholder="https://jira.com/issues" 
                />
              {/if}
            </div>
          {/if}
        </div>
      </div>

      <div class="advanced-section">
        <header class="section-header-inline">
          <span class="material-symbols-rounded">merge</span>
          <h4>Pull Requests</h4>
        </header>
        <div class="section-content">
          <Checkbox 
            label="Habilitar Pull Requests" 
            description="Revisiones de código y fusiones." 
            bind:checked={settings.enable_pulls} 
          />
          {#if settings.enable_pulls}
            <div class="nested-options">
              <Checkbox 
                label="Ignorar Espacios en Blanco" 
                description="Ocultar cambios de espaciado en los diffs." 
                bind:checked={settings.pulls_ignore_whitespace} 
              />
              <Checkbox 
                label="Permitir Rebase" 
                description="Habilitar rebase para commits de PR." 
                bind:checked={settings.pulls_allow_rebase} 
              />
              <Checkbox 
                label="Permitir Fusión Manual" 
                description="Habilitar la capacidad de marcar PRs como fusionados manualmente." 
                bind:checked={settings.pulls_allow_manual_merge} 
              />
            </div>
          {/if}
        </div>
      </div>

      <div class="form-actions">
        <button type="submit" class="btn btn-primary" disabled={saving}>
          {saving ? "Saving..." : "Save Settings"}
        </button>
      </div>
    </form>
  {/if}
</div>

<div class="settings-card danger">
  <header class="card-header">
    <h3 class="card-title red">Danger Zone</h3>
  </header>
  <div class="card-body danger-body">
    <div class="danger-item">
      <div class="danger-text">
        <strong>Delete this repository</strong>
        <p>Once you delete a repository, there is no going back. Please be certain.</p>
      </div>
      <div class="delete-action">
        <Input bind:value={deleteConfirm} placeholder="Type repository name to confirm" />
        <button
          class="btn btn-danger"
          disabled={deleteConfirm !== settings.name}
          onclick={deleteRepo}
        >
          Delete Repository
        </button>
      </div>
    </div>

    <div class="danger-item">
      <div class="danger-text">
        <strong>Transferir Propiedad</strong>
        <p>Transferir este repositorio a otro usuario o a una organización.</p>
      </div>
      <div class="transfer-action">
        <button class="btn btn-tonal" onclick={transferRepo} disabled={saving}>
          Transferir repositorio
        </button>
      </div>
    </div>
  </div>
</div>

<style>
  .settings-card { background: var(--md-sys-color-surface-container); border-radius: 2rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 800; }
  .card-title.red { color: var(--md-sys-color-error); }

  .settings-form { padding: 2.5rem; display: flex; flex-direction: column; gap: 2.5rem; }
  
  .avatar-settings-section { display: flex; align-items: center; gap: 3rem; padding-bottom: 3rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); margin-bottom: 1rem; }
  .avatar-preview-wrap { position: relative; border-radius: 50%; overflow: hidden; background: var(--md-sys-color-surface-container-highest); border: 2px solid var(--md-sys-color-outline-variant); flex-shrink: 0; }
  .avatar-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; color: white; }
  .spin { animation: spin 1s linear infinite; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

  .avatar-actions { display: flex; flex-direction: column; gap: 1.5rem; }
  .avatar-actions h4 { margin: 0; font-size: 1.25rem; font-weight: 900; }
  .avatar-actions p { margin: 0; color: var(--text-muted); font-size: 1rem; }
  .action-buttons { display: flex; gap: 1rem; align-items: center; }

  @media (max-width: 600px) {
    .avatar-settings-section { flex-direction: column; text-align: center; gap: 2rem; }
    .action-buttons { justify-content: center; }
  }

  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  :global(.full-width) { grid-column: span 2; }

  @media (max-width: 800px) {
    .form-grid { grid-template-columns: 1fr; }
    :global(.full-width) { grid-column: span 1; }
  }

  .options-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  @media (max-width: 800px) { .options-grid { grid-template-columns: 1fr; } }

  :global(.archived-item) { border: 1px solid rgba(var(--primary-rgb), 0.3) !important; background: rgba(var(--primary-rgb), 0.05) !important; }

  .form-actions { display: flex; justify-content: flex-end; padding-top: 2rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  .advanced-section { border-top: 1px solid var(--md-sys-color-outline-variant); padding-top: 2rem; }
  .section-header-inline { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.5rem; color: var(--primary-color); }
  .section-header-inline span { font-size: 1.5rem; opacity: 0.8; }
  .section-header-inline h4 { margin: 0; font-size: 1.1rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.05em; }
  
  .section-content { display: flex; flex-direction: column; gap: 1rem; }
  .nested-options { margin-left: 2rem; padding-left: 1.5rem; border-left: 2px solid var(--md-sys-color-outline-variant); display: flex; flex-direction: column; gap: 1rem; }

  .danger { margin-top: 4rem; border-color: var(--md-sys-color-error); background: rgba(var(--error-rgb), 0.02); }
  .danger-body { padding: 2.5rem; display: flex; flex-direction: column; gap: 2rem; }
  
  .danger-item { display: flex; flex-direction: column; gap: 1.5rem; }
  .danger-text strong { color: var(--md-sys-color-error); font-size: 1.25rem; font-weight: 850; }
  .danger-text p { color: var(--text-muted); margin: 0.5rem 0 0; }

  .delete-action { display: flex; gap: 1rem; align-items: flex-end; }
  .delete-action :global(.field) { flex: 1; }

  .danger-item:not(:first-child) { padding-top: 2rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  .skeleton-form { height: 400px; background: var(--md-sys-color-surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
