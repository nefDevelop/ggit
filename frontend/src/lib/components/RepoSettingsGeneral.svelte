<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { goto } from "$app/navigation";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  import { RepositorySchema } from "$lib/schemas";
  import Checkbox from "./ui/Checkbox.svelte";
  import Input from "./ui/Input.svelte";

  const repo: any = getContext("repo");
  let settings = $state<any>({});
  let loading = $state(true);
  let saving = $state(false);
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
          label="Issues" 
          description="Lightweight issue tracking." 
          bind:checked={settings.enable_issues} 
        />
        <Checkbox 
          label="Pull Requests" 
          description="Code reviews and merges." 
          bind:checked={settings.enable_pulls} 
        />
        <Checkbox 
          label="Wiki" 
          description="Project documentation." 
          bind:checked={settings.enable_wiki} 
        />
        <Checkbox 
          label="Archivar Repositorio" 
          description="Convertir en solo lectura. Útil para proyectos terminados." 
          bind:checked={settings.archived}
          class="archived-item"
        />
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

  .danger { margin-top: 3rem; border-color: var(--md-sys-color-error); background: rgba(var(--error-rgb), 0.02); }
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
