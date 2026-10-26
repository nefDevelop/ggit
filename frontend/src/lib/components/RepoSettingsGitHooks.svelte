<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Modal from "$lib/components/Modal.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Card from "$lib/components/ui/Card.svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  import { toast } from "$lib/toast.svelte";
  import { i18n } from "$lib/i18n.svelte";
  

import Terminal from "lucide-svelte/icons/terminal";
import AlertCircle from "lucide-svelte/icons/alert-circle";
import Code from "lucide-svelte/icons/code";
import Edit from "lucide-svelte/icons/edit";
import Info from "lucide-svelte/icons/info";



  const repo = getContext<RepoContext>('repo');

  let hooks = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  let editingHook = $state<any>(null);
  let hookContent = $state("");
  let isSubmitting = $state(false);

  async function fetchHooks() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/git`,
      );
      if (res.ok) {
        hooks = await res.json();
      } else {
        error = "Failed to load git hooks.";
      }
    } catch (e) {
      error = "An error occurred while fetching git hooks.";
    } finally {
      loading = false;
    }
  }

  function handleEdit(hook: any) {
    editingHook = hook;
    hookContent = hook.content;
  }

  async function handleSave() {
    isSubmitting = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/git/${editingHook.name}`,
        {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ content: hookContent }),
        },
      );

      if (res.ok) {
        toast.success(`Hook ${editingHook.name} updated successfully.`);
        editingHook = null;
        fetchHooks();
      } else {
        const data = await res.json();
        toast.error(data.message || "Failed to update hook.");
      }
    } catch (e) {
      console.error(e);
      toast.error("An error occurred while saving the hook.");
    } finally {
      isSubmitting = false;
    }
  }

  onMount(fetchHooks);
</script>

<Card padding="none" class="settings-premium-card">
  <PageHeader 
    title="Git Hooks" 
    subtitle="Gestiona scripts personalizados que se ejecutan en el servidor ante eventos de Git." 
    icon={Terminal}
  />

  <div class="settings-content">
    {#if loading}
      <div class="skeleton-list">
        {#each Array(3) as _}
          <div class="skeleton-item tonal"></div>
        {/each}
      </div>
    {:else if error}
      <div class="error-banner">
        <AlertCircle size={24} />
        {error}
      </div>
    {:else}
      <div class="hooks-grid">
        {#each hooks as hook}
          <div class="hook-tile tonal">
            <div class="hook-main">
              <div class="hook-identity">
                <Code size={18} class="hook-icon" />
                <div class="hook-names">
                  <span class="hook-name">{hook.name}</span>
                  {#if hook.is_active}
                    <span class="badge active">Configurado</span>
                  {:else}
                    <span class="badge">No configurado</span>
                  {/if}
                </div>
              </div>
              <p class="hook-desc">
                {#if hook.name === 'pre-receive'}
                  Se ejecuta antes de aceptar cualquier commit. Ideal para validar políticas.
                {:else if hook.name === 'update'}
                  Se ejecuta por cada rama que se intenta actualizar.
                {:else if hook.name === 'post-receive'}
                  Se ejecuta después de que los cambios han sido aceptados. Útil para despliegues o avisos.
                {:else}
                  Script de servidor ejecutado durante el ciclo de vida de Git.
                {/if}
              </p>
            </div>
            <div class="hook-footer">
              <Button variant="secondary" size="sm" icon={Edit} onclick={() => handleEdit(hook)}>
                Editar Script
              </Button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</Card>

<Modal 
  show={!!editingHook} 
  title={`Editar Git Hook: ${editingHook?.name}`} 
  onclose={() => editingHook = null}
  maxWidth="900px"
>
  {#if editingHook}
    <div class="hook-editor">
      <div class="editor-info island">
        <Info size={18} />
        <p>Este script se ejecutará en el servidor. Asegúrate de incluir el shebang correcto (ej: <code>#!/bin/bash</code>).</p>
      </div>
      
      <div class="textarea-container">
        <textarea 
          bind:value={hookContent} 
          class="hook-textarea"
          placeholder="#!/bin/bash\n\n# Tu lógica personalizada aquí..."
          spellcheck="false"
        ></textarea>
      </div>

      <footer class="modal-actions">
        <Button variant="ghost" onclick={() => editingHook = null}>Cancelar</Button>
        <Button variant="primary" loading={isSubmitting} onclick={handleSave}>
          Guardar Cambios
        </Button>
      </footer>
    </div>
  {/if}
</Modal>

<style>
  .settings-content { padding: 2.5rem; }
  @media (max-width: 600px) { .settings-content { padding: 1.5rem; } }

  .hooks-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.5rem; }
  @media (max-width: 800px) { .hooks-grid { grid-template-columns: 1fr; } }

  .hook-tile { padding: 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); display: flex; flex-direction: column; gap: 1.5rem; transition: all 0.2s; }
  .hook-tile:hover { transform: translateY(-3px); border-color: var(--primary-color); }
  
  .hook-identity { display: flex; align-items: center; gap: 1rem; }
  :global(.hook-icon) { color: var(--primary-color); opacity: 0.7; }
  
  .hook-names { display: flex; flex-direction: column; gap: 0.1rem; }
  .hook-name { font-family: 'JetBrains Mono', monospace; font-weight: 900; font-size: 1rem; color: var(--text-color); }
  
  .badge { font-size: 0.65rem; font-weight: 950; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 100px; width: fit-content; background: var(--md-sys-color-surface-container-highest); color: var(--text-muted); opacity: 0.6; }
  .badge.active { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); opacity: 1; }

  .hook-desc { margin: 0; font-size: 0.9rem; color: var(--text-muted); line-height: 1.5; opacity: 0.8; }

  .hook-footer { display: flex; justify-content: flex-end; margin-top: auto; padding-top: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  /* EDITOR */
  .hook-editor { display: flex; flex-direction: column; gap: 1.5rem; }
  .editor-info { display: flex; align-items: center; gap: 1rem; color: var(--text-muted); padding: 1rem 1.5rem; border-radius: 1rem; font-size: 0.9rem; }
  .editor-info code { background: var(--sunken-3); padding: 0.1rem 0.3rem; border-radius: 4px; color: var(--primary-color); }

  .textarea-container { border-radius: 1.5rem; overflow: hidden; border: 2px solid var(--md-sys-color-outline-variant); background: #000; }
  .hook-textarea { width: 100%; height: 400px; background: transparent; color: #fff; border: none; padding: 1.5rem; font-family: 'JetBrains Mono', 'Fira Code', monospace; font-size: 0.9rem; line-height: 1.6; resize: vertical; outline: none; }
  .hook-textarea:focus { caret-color: var(--primary-color); }

  .modal-actions { display: flex; justify-content: flex-end; gap: 1rem; padding-top: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  .skeleton-list { display: flex; flex-direction: column; gap: 1.5rem; }
  .skeleton-item { height: 120px; border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }

  .error-banner { display: flex; align-items: center; gap: 1rem; padding: 1.5rem; border-radius: 1.5rem; background: rgba(var(--error-rgb), 0.1); color: var(--error); font-weight: 700; }
</style>
