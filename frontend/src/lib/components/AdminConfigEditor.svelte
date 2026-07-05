<script lang="ts">
  import { apiFetch } from "$lib/api";
  import Button from './ui/Button.svelte';
  

import Settings from "lucide-svelte/icons/settings";
import X from "lucide-svelte/icons/x";
import TriangleAlert from "lucide-svelte/icons/triangle-alert";


  
  let { section, data, onSave, onCancel } = $props();
  
  let editingData = $state<Record<string, any>>({});
  $effect(() => {
    if (data) editingData = $state.snapshot(data);
  });
  let saving = $state(false);
  let error = $state('');

  async function handleSubmit() {
    saving = true;
    error = '';
    try {
      const res = await apiFetch(`/api/v2/admin/config`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          section: section,
          data: editingData
        })
      });
      
      if (res.ok) {
        onSave();
      } else {
        error = 'Error al actualizar configuración';
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
      <div class="title-wrap">
        <Settings size={20} />
        <h3>Editar Sección: [{section}]</h3>
      </div>
      <Button variant="ghost" icon={X} onclick={onCancel} />
    </header>

    <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      <div class="form-scroll-area">
        <div class="form-grid">
          {#each Object.entries(editingData) as [key, value]}
            <div class="form-group">
              <label for={"key-"+key}>{key}</label>
              <input 
                type="text" 
                id={"key-"+key} 
                bind:value={editingData[key]} 
                placeholder="Valor..." 
              />
            </div>
          {/each}
        </div>
      </div>

      {#if error}
        <p class="error-msg">{error}</p>
      {/if}

      <footer class="modal-footer">
        <div class="warning-box">
          <TriangleAlert size={20} />
          <p>Algunos cambios pueden requerir un reinicio del servidor para aplicarse.</p>
        </div>
        <div class="actions">
          <Button variant="ghost" onclick={onCancel}>Cancelar</Button>
          <Button type="submit" variant="primary" disabled={saving}>
            {saving ? 'Guardando...' : 'Aplicar Cambios'}
          </Button>
        </div>
      </footer>
    </form>
  </div>
</div>

<style>
  .editor-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.8); backdrop-filter: blur(12px); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 2rem; }
  .editor-modal { width: 100%; max-width: 900px; max-height: 90vh; display: flex; flex-direction: column; background: var(--md-sys-color-surface); border: 1px solid var(--md-sys-color-outline-variant); overflow: hidden; box-shadow: 0 24px 48px rgba(0,0,0,0.4); }
  
  .modal-header { padding: 1.5rem 2rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--md-sys-color-outline-variant); background: rgba(var(--primary-rgb), 0.05); }
  .title-wrap { display: flex; align-items: center; gap: 1rem; }
  .modal-header h3 { margin: 0; font-size: 1.25rem; font-weight: 900; letter-spacing: -0.02em; }

  form { display: flex; flex-direction: column; overflow: hidden; }
  .form-scroll-area { padding: 2rem; overflow-y: auto; flex: 1; mask-image: linear-gradient(to bottom, transparent, black 20px, black calc(100% - 20px), transparent); }
  
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; padding: 1rem 0; }
  @media (max-width: 600px) { .form-grid { grid-template-columns: 1fr; } }

  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.8; }
  
  input { background: rgba(var(--primary-rgb), 0.03); border: 1px solid var(--md-sys-color-outline-variant); padding: 0.85rem 1.25rem; border-radius: 0.75rem; color: var(--text-color); font-family: 'JetBrains Mono', monospace; font-size: 0.95rem; }
  input:focus { outline: none; border-color: var(--md-sys-color-primary); background: var(--md-sys-color-surface-container-highest); box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary), transparent 85%); }

  .error-msg { color: var(--error); font-size: 0.9rem; font-weight: 700; margin: 0 2rem 1rem; padding: 1rem; background: rgba(var(--error-rgb), 0.1); border-radius: 1rem; }

  .modal-footer { padding: 1.5rem 2rem; border-top: 1px solid var(--md-sys-color-outline-variant); display: flex; align-items: center; justify-content: space-between; gap: 2rem; background: rgba(0,0,0,0.1); }
  
  .warning-box { display: flex; align-items: center; gap: 0.75rem; color: #fbbf24; }
  .warning-box p { margin: 0; font-size: 0.85rem; font-weight: 700; max-width: 300px; line-height: 1.3; }

  .actions { display: flex; gap: 1rem; }
</style>
