<script lang="ts">
  import { apiFetch } from "$lib/api";
  import Button from './ui/Button.svelte';
  
  let { org, onSave, onCancel } = $props();
  
  let editingOrg = $state({ username: '', full_name: '', website: '', description: '', location: '' });
  $effect(() => {
    if (org) editingOrg = { ...org };
  });
  let saving = $state(false);
  let error = $state('');

  async function handleSubmit() {
    saving = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/orgs/${org.username}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(editingOrg)
      });
      
      if (res.ok) {
        onSave();
      } else {
        const data = await res.json();
        error = data.message || 'Error al actualizar organización';
      }
    } catch (e) {
      error = 'Error de conexión';
    } finally {
      saving = false;
    }
  }
</script>

<div class="editor-overlay">
  <div class="editor-modal island">
    <header class="modal-header">
      <h3>Editar Organización: {org.username}</h3>
      <Button variant="ghost" icon="close" onclick={onCancel} />
    </header>

    <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      <div class="form-grid">
        <div class="form-group">
          <label for="full_name">Nombre Público</label>
          <input type="text" id="full_name" bind:value={editingOrg.full_name} placeholder="Nombre de la entidad" />
        </div>

        <div class="form-group">
          <label for="website">Sitio Web</label>
          <input type="url" id="website" bind:value={editingOrg.website} placeholder="https://..." />
        </div>

        <div class="form-group col-span-2">
          <label for="description">Descripción</label>
          <textarea id="description" bind:value={editingOrg.description} placeholder="Breve biografía de la organización"></textarea>
        </div>

        <div class="form-group">
          <label for="location">Ubicación</label>
          <input type="text" id="location" bind:value={editingOrg.location} placeholder="Ciudad, País" />
        </div>
      </div>

      {#if error}
        <p class="error-msg">{error}</p>
      {/if}

      <footer class="modal-footer">
        <Button variant="ghost" onclick={onCancel}>Cancelar</Button>
        <Button type="submit" variant="primary" disabled={saving}>
          {saving ? 'Guardando...' : 'Guardar Cambios'}
        </Button>
      </footer>
    </form>
  </div>
</div>

<style>
  .editor-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.8); backdrop-filter: blur(8px); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 2rem; }
  .editor-modal { width: 100%; max-width: 700px; background: var(--md-sys-color-surface); border: 1px solid var(--md-sys-color-outline-variant); overflow: hidden; }
  
  .modal-header { padding: 1.5rem 2rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .modal-header h3 { margin: 0; font-size: 1.5rem; font-weight: 800; color: var(--primary-color); }

  form { padding: 2rem; }
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  @media (max-width: 600px) { .form-grid { grid-template-columns: 1fr; } }
  .col-span-2 { grid-column: span 2; }

  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.8rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); }
  
  input, textarea { background: rgba(var(--primary-rgb), 0.05); border: 1px solid var(--md-sys-color-outline-variant); padding: 0.75rem 1rem; border-radius: 1rem; color: var(--text-color); font-family: inherit; font-size: 1rem; }
  textarea { min-height: 100px; resize: vertical; }
  input:focus, textarea:focus { outline: none; border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); }

  .error-msg { color: var(--error); font-size: 0.9rem; font-weight: 600; margin-top: 1rem; }

  .modal-footer { margin-top: 2rem; display: flex; justify-content: flex-end; gap: 1rem; padding-top: 1.5rem; border-top: 1px solid var(--md-sys-color-outline-variant); }
</style>
