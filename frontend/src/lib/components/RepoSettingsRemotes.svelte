<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { Trash2, Plus, Globe, Shield } from "lucide-svelte";

  const repo: any = getContext("repo");
  let remotes = $state<string[]>([]);
  let loading = $state(true);
  let adding = $state(false);
  let deleting = $state<string | null>(null);

  let newName = $state("");
  let newURL = $state("");
  let message = $state("");
  let error = $state("");

  async function fetchRemotes() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/remotes`
      );
      if (res.ok) remotes = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addRemote() {
    if (!newName || !newURL) return;
    adding = true;
    error = "";
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/remotes`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name: newName, url: newURL }),
        }
      );
      if (res.ok) {
        newName = "";
        newURL = "";
        message = "Remoto añadido correctamente";
        setTimeout(() => (message = ""), 3000);
        fetchRemotes();
      } else {
        const data = await res.json();
        error = data.message || "Error al añadir remoto";
      }
    } catch (e) {
      error = "Error de conexión";
    } finally {
      adding = false;
    }
  }

  async function deleteRemote(name: string) {
    if (!await dialog.confirm(`¿Estás seguro de eliminar el remoto "${name}"?`)) return;
    deleting = name;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/remotes/${name}`,
        { method: "DELETE" }
      );
      if (res.ok) {
        fetchRemotes();
      }
    } catch (e) {
      console.error(e);
    } finally {
      deleting = null;
    }
  }

  onMount(fetchRemotes);
</script>

<div class="remotes-settings">
  <Card padding="lg">
    <header class="section-header">
      <div class="title-info">
        <h3>Gestionar Remotos</h3>
        <p>Configura destinos adicionales para sincronizar tu código.</p>
      </div>
    </header>

    <div class="add-remote-form tonal">
      <div class="form-row">
        <div class="input-group">
          <label for="remote-name">Nombre</label>
          <input 
            type="text" 
            id="remote-name" 
            bind:value={newName} 
            placeholder="ej: github, backup" 
          />
        </div>
        <div class="input-group grow">
          <label for="remote-url">URL del Repositorio</label>
          <input 
            type="text" 
            id="remote-url" 
            bind:value={newURL} 
            placeholder="https://github.com/usuario/repo.git" 
          />
        </div>
        <div class="btn-container">
          <Button 
            variant="primary" 
            onclick={addRemote} 
            disabled={adding || !newName || !newURL}
            icon="add"
          >
            {adding ? "Añadiendo..." : "Añadir"}
          </Button>
        </div>
      </div>
      {#if error}<p class="error-text">{error}</p>{/if}
      {#if message}<p class="success-text">{message}</p>{/if}
    </div>

    <div class="remotes-list">
      {#if loading}
        <div class="loading-placeholder">Cargando remotos...</div>
      {:else if remotes.length === 0}
        <div class="empty-state">
          <Globe size={40} />
          <p>No hay remotos adicionales configurados.</p>
        </div>
      {:else}
        <div class="list-grid">
          {#each remotes as remote}
            <div class="remote-item tonal">
              <div class="remote-info">
                <div class="remote-icon">
                  <Shield size={20} />
                </div>
                <div class="remote-meta">
                  <span class="remote-name">{remote}</span>
                  <span class="remote-hint">Configurado en el servidor</span>
                </div>
              </div>
              <button 
                class="btn-delete" 
                onclick={() => deleteRemote(remote)}
                disabled={deleting === remote}
                title="Eliminar remoto"
              >
                <Trash2 size={18} />
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </Card>

  <Card class="info-card variant" padding="md">
    <div class="info-content">
      <Globe size={24} />
      <div>
        <strong>Nota sobre Autenticación</strong>
        <p>Actualmente, el servidor de Ggit debe tener acceso pre-configurado (vía SSH o tokens en el entorno) para interactuar con estos remotos externos.</p>
      </div>
    </div>
  </Card>
</div>

<style>
  .remotes-settings {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .section-header {
    margin-bottom: 2rem;
  }

  .section-header h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 800;
  }

  .section-header p {
    margin: 0.5rem 0 0;
    color: var(--text-muted);
    font-size: 0.9rem;
  }

  .add-remote-form {
    padding: 1.5rem;
    border-radius: 1.25rem;
    margin-bottom: 2rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .form-row {
    display: flex;
    gap: 1rem;
    align-items: flex-end;
  }

  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .input-group.grow { flex: 1; }

  .input-group label {
    font-size: 0.75rem;
    font-weight: 800;
    text-transform: uppercase;
    color: var(--text-muted);
    padding-left: 0.5rem;
  }

  input {
    background: rgba(0, 0, 0, 0.1);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 0.75rem 1rem;
    border-radius: 10px;
    color: var(--text-color);
    font-size: 0.9rem;
    outline: none;
    transition: all 0.2s;
  }

  input:focus {
    border-color: var(--primary-color);
    background: rgba(0, 0, 0, 0.2);
  }

  .remotes-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .list-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
  }

  .remote-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1rem 1.25rem;
    border-radius: 1.25rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .remote-info {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .remote-icon {
    width: 40px;
    height: 40px;
    border-radius: 10px;
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .remote-meta {
    display: flex;
    flex-direction: column;
  }

  .remote-name {
    font-weight: 800;
    font-size: 1rem;
  }

  .remote-hint {
    font-size: 0.75rem;
    color: var(--text-muted);
  }

  .btn-delete {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 8px;
    transition: all 0.2s;
  }

  .btn-delete:hover {
    color: #ef4444;
    background: rgba(239, 68, 68, 0.1);
  }

  .empty-state {
    padding: 4rem 2rem;
    text-align: center;
    color: var(--text-muted);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .error-text { color: #ef4444; font-size: 0.85rem; margin-top: 1rem; }
  .success-text { color: #10b981; font-size: 0.85rem; margin-top: 1rem; }

  .info-card.variant {
    background: rgba(var(--primary-rgb), 0.05);
    border: 1px dashed var(--primary-color);
  }

  .info-content {
    display: flex;
    gap: 1.25rem;
    align-items: flex-start;
    color: var(--text-color);
  }

  .info-content strong { display: block; margin-bottom: 0.25rem; }
  .info-content p { margin: 0; font-size: 0.85rem; color: var(--text-muted); }
</style>
