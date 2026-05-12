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
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import { Settings, RefreshCw, Upload, Trash2, FileText, Bug, GitPullRequest, AlertTriangle } from "lucide-svelte";

  import PageHeader from "./ui/PageHeader.svelte";

  const repo: any = getContext("repo");
  let settings = $state<any>({});
  let branches = $state<any[]>([]);
  let loading = $state(true);
  let saving = $state(false);
  let avatarUploading = $state(false);
  let deleteConfirm = $state("");
  let errors = $state<Record<string, string>>({});

  async function fetchBranches() {
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/branches`);
      if (res.ok) {
        const d = await res.json();
        const b = [];
        if (d.default) b.push(d.default);
        if (d.active) b.push(...d.active);
        if (d.stale) b.push(...d.stale);
        branches = b;
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function fetchSettings() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings`,
      );
      if (res.ok) settings = await res.json();
      await fetchBranches();
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
    const result = RepositorySchema.safeParse(settings);
    if (!result.success) {
      result.error.issues.forEach(err => {
        const key = err.path[0];
        if (typeof key === 'string') {
          errors[key] = err.message;
        }
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
          body: JSON.stringify({
            ...result.data,
            project_folder: settings.project_folder,
            branch: settings.default_branch,
            default_branch: settings.default_branch
          }),
        },
      );
      if (res.ok) {
        toast.success("Repository updated successfully!");
        fetchSettings();
      }
    } catch (e) { } finally {
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

<Card padding="none" class="settings-premium-card">
  <PageHeader 
    title="Ajustes del Repositorio"
    subtitle="Configuración general, visibilidad y metadatos básicos de tu proyecto."
    icon={Settings}
  />

  <div class="settings-content">
    {#if loading}
      <div class="skeleton-form"></div>
    {:else}
      <form class="settings-form" onsubmit={(e) => { e.preventDefault(); saveSettings(); }}>
        <div class="avatar-settings-section island">
          <div class="avatar-preview-wrap">
            <Avatar 
              src={settings.avatar_url} 
              username={settings.name} 
              seed={settings.id?.toString()} 
              size={110} 
              class="repo-avatar-img"
            />
            {#if avatarUploading}
              <div class="avatar-overlay">
                <RefreshCw size={24} class="spin" />
              </div>
            {/if}
          </div>
          <div class="avatar-actions">
            <header>
              <strong>Imagen del Repositorio</strong>
              <p>Cambia el icono que identifica a este proyecto en la plataforma.</p>
            </header>
            <div class="action-buttons">
              <Button variant="outline" size="sm" icon={Upload} onclick={() => document.getElementById('repo-avatar-input')?.click()}>
                Subir Imagen
              </Button>
              {#if settings.use_custom_avatar}
                <Button variant="ghost" size="sm" icon={Trash2} class="danger-text" onclick={deleteAvatar}>
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
            label="Nombre del Repositorio" 
            bind:value={settings.name} 
            error={errors.name}
          />
          <Input 
            label="Sitio Web" 
            bind:value={settings.website} 
            error={errors.website}
            placeholder="https://example.com" 
          />
          <Input 
            label="Carpeta de Proyecto" 
            bind:value={settings.project_folder} 
            error={errors.project_folder}
            placeholder="p.ej. Backend, Frontend, Personal" 
          />

          <div class="select-field full-width">
            <label for="default_branch">Rama por Defecto</label>
            <div class="select-wrapper">
              <select id="default_branch" bind:value={settings.default_branch}>
                {#each branches as branch}
                  <option value={branch.name}>{branch.name}</option>
                {/each}
              </select>
            </div>
            <span class="select-helper-text">La rama por defecto es la que se muestra predeterminadamente al ver el código del proyecto.</span>
          </div>
          
          <Input 
            label="Descripción" 
            bind:value={settings.description} 
            error={errors.description}
            placeholder="Una breve descripción de tu proyecto..."
            multiline
            class="full-width"
          />
        </div>

        <div class="options-grid">
          <Checkbox 
            label="Hacer Privado" 
            description="Solo visible para los colaboradores autorizados." 
            bind:checked={settings.private} 
          />
          <Checkbox 
            label="No Listado" 
            description="El repositorio no aparecerá en los resultados de búsqueda pública." 
            bind:checked={settings.unlisted} 
          />
          <Checkbox 
            label="Archivar Repositorio" 
            description="Convertir en solo lectura. Bloquea commits, issues y pull requests." 
            bind:checked={settings.archived}
            class="archived-item"
          />
        </div>

        <div class="advanced-stack">
          <div class="advanced-section tonal">
            <header class="adv-header">
              <FileText size={24} />
              <strong>Documentación y Wiki</strong>
            </header>
            <div class="adv-content">
              <Checkbox 
                label="Habilitar Wiki" 
                description="Habilitar sistema de documentación integrado." 
                bind:checked={settings.enable_wiki} 
              />
              {#if settings.enable_wiki}
                <div class="nested-box">
                  <Input 
                    label="Nombre de Carpeta Wiki" 
                    bind:value={settings.wiki_folder_name} 
                    placeholder="docs" 
                  />
                  <Input 
                    label="Nombre de Carpeta ADR" 
                    bind:value={settings.adr_folder_name} 
                    placeholder="ADR" 
                  />
                  <Checkbox 
                    label="Acceso Público" 
                    description="Permitir que cualquiera vea la wiki incluso en repositorios privados." 
                    bind:checked={settings.allow_public_wiki} 
                  />
                  <Checkbox 
                    label="Wiki Externa" 
                    description="Redirigir la pestaña Wiki a una URL externa." 
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

          <div class="advanced-section tonal">
            <header class="adv-header">
              <Bug size={24} />
              <strong>Issues y Seguimiento</strong>
            </header>
            <div class="adv-content">
              <Checkbox 
                label="Habilitar Issues" 
                description="Seguimiento de errores y tareas integrado." 
                bind:checked={settings.enable_issues} 
              />
              {#if settings.enable_issues}
                <div class="nested-box">
                  <Checkbox 
                    label="Issues Públicos" 
                    description="Permitir visualización pública de issues en repositorios privados." 
                    bind:checked={settings.allow_public_issues} 
                  />
                  <Checkbox 
                    label="Tracker Externo" 
                    description="Vincular a un sistema de gestión externo." 
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

          <div class="advanced-section tonal">
            <header class="adv-header">
              <GitPullRequest size={24} />
              <strong>Pull Requests</strong>
            </header>
            <div class="adv-content">
              <Checkbox 
                label="Habilitar Pull Requests" 
                description="Revisiones de código y flujo de trabajo colaborativo." 
                bind:checked={settings.enable_pulls} 
              />
              {#if settings.enable_pulls}
                <div class="nested-box">
                  <Checkbox 
                    label="Ignorar Espacios" 
                    description="No mostrar cambios de espacios en los diffs de PR." 
                    bind:checked={settings.pulls_ignore_whitespace} 
                  />
                  <Checkbox 
                    label="Permitir Rebase" 
                    description="Habilitar la opción de rebase al fusionar." 
                    bind:checked={settings.pulls_allow_rebase} 
                  />
                  <Checkbox 
                    label="Fusión Manual" 
                    description="Permitir marcar PRs como fusionados manualmente." 
                    bind:checked={settings.pulls_allow_manual_merge} 
                  />
                </div>
              {/if}
            </div>
          </div>
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

<style>
  .settings-content { padding: 0; }
  @media (max-width: 1300px) { .settings-content { padding: 0; } }
  @media (max-width: 600px) { .settings-content { padding: 0; } }

  .settings-form { display: flex; flex-direction: column; gap: 3rem; }

  .avatar-settings-section { display: flex; align-items: center; gap: 3rem; padding: 2rem; border-radius: 2rem; }
  .avatar-preview-wrap { position: relative; border-radius: 2.5rem; overflow: hidden; background: var(--md-sys-color-surface-container-highest); border: 2px solid var(--md-sys-color-outline-variant); flex-shrink: 0; box-shadow: var(--shadow-md); }
  .avatar-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; color: white; }
  :global(.spin) { animation: spin 1s linear infinite; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

  .avatar-actions { display: flex; flex-direction: column; gap: 1rem; }
  .avatar-actions strong { font-size: 1.25rem; font-weight: 900; color: var(--text-color); }
  .avatar-actions p { margin: 0; color: var(--text-muted); font-size: 0.95rem; line-height: 1.4; }
  .action-buttons { display: flex; gap: 0.75rem; align-items: center; }

  @media (max-width: 1300px) {
    .avatar-settings-section { gap: 1.5rem; padding: 1.5rem; border-radius: 1.5rem; }
    .avatar-preview-wrap { border-radius: 1.75rem; }
    :global(.repo-avatar-img) { width: 90px !important; height: 90px !important; }
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
  .danger-zone-container { margin-top: 4rem; }
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
