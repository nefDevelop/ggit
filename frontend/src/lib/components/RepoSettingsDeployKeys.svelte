<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Modal from "$lib/components/Modal.svelte";
  import Input from "$lib/components/ui/Input.svelte";
  import Checkbox from "$lib/components/ui/Checkbox.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Card from "$lib/components/ui/Card.svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  import { i18n } from "$lib/i18n.svelte";
  
import { Key, Plus } from "lucide-svelte";


  const repo: any = getContext("repo");

  let keys = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  let showModal = $state(false);
  let newKey = $state({ title: "", key: "", read_only: true });
  let isSubmitting = $state(false);

  async function fetchKeys() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/keys`,
      );
      if (res.ok) {
        keys = await res.json();
      } else {
        error = "Failed to load deploy keys.";
      }
    } catch (e) {
      error = "An error occurred while fetching deploy keys.";
    } finally {
      loading = false;
    }
  }

  async function handleAdd() {
    isSubmitting = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/keys`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(newKey),
        },
      );

      if (res.ok) {
        const key = await res.json();
        keys = [...keys, key];
        showModal = false;
        newKey = { title: "", key: "", read_only: true };
      } else {
        const data = await res.json();
        toast.error(data.message || "Failed to add deploy key.");
      }
    } catch (e) {
      console.error(e);
      toast.error("An error occurred while adding the key.");
    } finally {
      isSubmitting = false;
    }
  }

  async function handleDelete(id: number) {
    if (!(await dialog.confirm("¿Estás seguro de que deseas eliminar esta clave de despliegue?"))) return;

    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/keys/${id}`,
        {
          method: "DELETE",
        },
      );
      if (res.ok) {
        keys = keys.filter((k) => k.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  }

  onMount(fetchKeys);
</script>

<Card padding="none" class="settings-premium-card">
  <PageHeader 
    title="Claves de Despliegue"
    subtitle="Permite que servidores externos clonen este repositorio vía SSH sin una cuenta de usuario."
    icon={Key}
  >
    {#snippet actions()}
      <Button variant="primary" icon={Plus} onclick={() => (showModal = true)}>Añadir Clave</Button>
    {/snippet}
  </PageHeader>

  <div class="settings-content">
    {#if loading}
      <div class="skeleton-stack">
        {#each Array(2) as _}
          <div class="skeleton-item tonal"></div>
        {/each}
      </div>
    {:else if error}
      <div class="error-banner">
        <span class="material-symbols-rounded">error</span>
        {error}
      </div>
    {:else if keys.length === 0}
      <div class="empty-state island">
        <span class="material-symbols-rounded">key_off</span>
        <p>No hay claves de despliegue configuradas para este repositorio.</p>
        <Button variant="outline" size="sm" onclick={() => (showModal = true)}>Configurar la primera</Button>
      </div>
    {:else}
      <div class="keys-stack">
        {#each keys as key}
          <div class="key-tile tonal">
            <div class="key-main">
              <div class="key-icon-wrap">
                <span class="material-symbols-rounded">key</span>
              </div>
              <div class="key-info">
                <div class="key-header-row">
                  <span class="key-title">{key.title}</span>
                  {#if key.read_only}
                    <span class="badge">Solo Lectura</span>
                  {:else}
                    <span class="badge active">Lectura y Escritura</span>
                  {/if}
                </div>
                <div class="key-fingerprint">
                  <code>{key.key.substring(0, 48)}...</code>
                </div>
                <div class="key-meta">
                  Añadida el {formatDate(key.created_at)}
                </div>
              </div>
            </div>
            <div class="key-actions">
              <Button 
                variant="ghost" 
                size="sm" 
                icon="delete" 
                class="danger-btn" 
                onclick={() => handleDelete(key.id)} 
                title="Eliminar"
              />
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</Card>

<Modal show={showModal} title="Añadir Clave de Despliegue" onclose={() => (showModal = false)}>
  <form
    onsubmit={(e) => {
      e.preventDefault();
      handleAdd();
    }}
    class="key-form"
  >
    <Input 
      label="Nombre de la Clave" 
      bind:value={newKey.title} 
      placeholder="Ej: Servidor de CI"
    />
    <Input 
      label="Contenido de la Clave (SSH)" 
      bind:value={newKey.key} 
      placeholder="ssh-rsa AAAAB3Nza..."
      multiline
      rows={5}
    />

    <Checkbox
      bind:checked={newKey.read_only}
      label="Solo lectura"
      description="Si está habilitado, esta clave no podrá realizar pushes al repositorio."
    />

    <footer class="modal-actions-premium">
      <Button variant="ghost" onclick={() => (showModal = false)}>Cancelar</Button>
      <Button type="submit" variant="primary" loading={isSubmitting}>
        Guardar Clave
      </Button>
    </footer>
  </form>
</Modal>

<style>
  .settings-content { padding: 2.5rem; }
  @media (max-width: 600px) { .settings-content { padding: 1.5rem; } }

  .keys-stack { display: flex; flex-direction: column; gap: 1rem; }
  
  .key-tile { padding: 1.25rem 1.75rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); display: flex; justify-content: space-between; align-items: center; gap: 2rem; transition: all 0.2s; }
  .key-tile:hover { transform: translateX(4px); border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.02); }

  .key-main { display: flex; align-items: center; gap: 1.5rem; flex: 1; min-width: 0; }
  .key-icon-wrap { width: 44px; height: 44px; border-radius: 1rem; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); display: flex; align-items: center; justify-content: center; }
  .key-icon-wrap span { font-size: 1.5rem; }

  .key-info { display: flex; flex-direction: column; gap: 0.25rem; flex: 1; min-width: 0; }
  .key-header-row { display: flex; align-items: center; gap: 1rem; }
  .key-title { font-weight: 900; font-size: 1.1rem; color: var(--text-color); }
  
  .badge { font-size: 0.65rem; font-weight: 950; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 100px; background: var(--md-sys-color-surface-container-highest); color: var(--text-muted); opacity: 0.6; }
  .badge.active { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); opacity: 1; }

  .key-fingerprint { margin: 0.25rem 0; }
  .key-fingerprint code { font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; color: var(--text-muted); opacity: 0.7; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; display: block; }

  .key-meta { font-size: 0.8rem; font-weight: 600; color: var(--text-muted); opacity: 0.5; }

  .key-actions { display: flex; gap: 0.5rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  @media (max-width: 800px) {
    .key-tile { flex-direction: column; align-items: stretch; gap: 1.5rem; }
    .key-actions { justify-content: flex-end; border-top: 1px solid var(--md-sys-color-outline-variant); padding-top: 1rem; }
  }

  .key-form { display: flex; flex-direction: column; gap: 1.5rem; }
  .modal-actions-premium { display: flex; justify-content: flex-end; gap: 1rem; padding-top: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  .empty-state { padding: 4rem 2rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-state span { font-size: 4rem; color: var(--primary-color); opacity: 0.2; }
  .empty-state p { font-weight: 700; color: var(--text-muted); }

  .skeleton-stack { display: flex; flex-direction: column; gap: 1rem; }
  .skeleton-item { height: 90px; border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }

  .error-banner { display: flex; align-items: center; gap: 1rem; padding: 1.5rem; border-radius: 1.5rem; background: rgba(var(--error-rgb), 0.1); color: var(--error); font-weight: 700; }
</style>
