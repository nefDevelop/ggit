<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import { toast } from "$lib/toast.svelte";
  

import ShieldCheck from "lucide-svelte/icons/shield-check";
import ShieldAlert from "lucide-svelte/icons/shield-alert";
import GitMerge from "lucide-svelte/icons/git-merge";
import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
import XCircle from "lucide-svelte/icons/x-circle";
import AlertCircle from "lucide-svelte/icons/alert-circle";



  let { index, onMerged } = $props<{ index: number, onMerged: () => void }>();
  const repo = getContext<RepoContext>('repo');

  let protection = $state<any>(null);
  let loading = $state(true);
  let merging = $state(false);

  async function fetchProtection() {
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/protection`);
      if (res.ok) protection = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function doMerge() {
    merging = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/pulls/${index}/merge`, {
        method: 'POST'
      });
      if (res.ok) {
        toast.success("Pull Request fusionado correctamente.");
        onMerged();
      } else {
        const data = await res.json();
        toast.error(data.error || "Error al fusionar.");
      }
    } catch (e) {
      toast.error("Error de red.");
    } finally {
      merging = false;
    }
  }

  onMount(fetchProtection);
</script>

<Card class="merge-box-card" padding="none">
  <div class="merge-box">
    {#if loading}
      <div class="loading-merge">Comprobando estado de fusión...</div>
    {:else if protection}
      <div class="merge-status-list">
        <!-- Requisito de Aprobaciones -->
        {#if protection.is_protected && protection.required_approvals > 0}
          <div class="status-item" class:success={protection.approvals >= protection.required_approvals} class:error={protection.approvals < protection.required_approvals}>
            <div class="status-icon">
              {#if protection.approvals >= protection.required_approvals}
                <CheckCircle2 size={20} />
              {:else}
                <ShieldAlert size={20} />
              {/if}
            </div>
            <div class="status-text">
              <strong>Aprobaciones: {protection.approvals} / {protection.required_approvals}</strong>
              <span>{protection.approvals >= protection.required_approvals ? 'Requisito de revisión cumplido.' : 'Se necesitan más aprobaciones antes de fusionar.'}</span>
            </div>
          </div>
        {/if}

        <!-- Requisito de Conversaciones Resueltas -->
        {#if protection.unresolved_count > 0}
          <div class="status-item error">
            <div class="status-icon"><AlertCircle size={20} /></div>
            <div class="status-text">
              <strong>Conversaciones pendientes: {protection.unresolved_count}</strong>
              <span>Todas las conversaciones deben ser resueltas antes de fusionar.</span>
            </div>
          </div>
        {:else}
          <div class="status-item success">
            <div class="status-icon"><CheckCircle2 size={20} /></div>
            <div class="status-text">
              <strong>Conversaciones resueltas</strong>
              <span>Todas las conversaciones han sido resueltas.</span>
            </div>
          </div>
        {/if}

        <!-- Requisito de CI/CD -->
        {#if protection.combined_status}
          <div class="status-item" class:success={protection.combined_status === 'success'} class:error={protection.combined_status === 'failure' || protection.combined_status === 'error'} class:warning={protection.combined_status === 'pending'}>
            <div class="status-icon">
              {#if protection.combined_status === 'success'}
                <CheckCircle2 size={20} />
              {:else if protection.combined_status === 'pending'}
                <AlertCircle size={20} />
              {:else}
                <XCircle size={20} />
              {/if}
            </div>
            <div class="status-text">
              <strong>Estado de CI/CD: {protection.combined_status}</strong>
              <span>
                {#if protection.combined_status === 'success'}
                  Todas las pruebas han pasado satisfactoriamente.
                {:else if protection.combined_status === 'pending'}
                  Hay pruebas en ejecución...
                {:else}
                  Algunas pruebas han fallado. Revisa los logs.
                {/if}
              </span>
            </div>
          </div>
        {/if}

        <!-- Estado General -->
        {#if protection.mergeable}
          <div class="status-item success highlight">
            <div class="status-icon"><GitMerge size={20} /></div>
            <div class="status-text">
              <strong>Lista para fusionar</strong>
              <span>Esta rama no tiene conflictos y cumple con todos los criterios de calidad.</span>
            </div>
          </div>
        {:else}
          <div class="status-item warning">
            <div class="status-icon"><AlertCircle size={20} /></div>
            <div class="status-text">
              <strong>Fusión bloqueada</strong>
              <span>{protection.merge_error || 'No se puede fusionar hasta que se cumplan todos los requisitos.'}</span>
            </div>
          </div>
        {/if}
      </div>

      <footer class="merge-actions">
        <Button 
          variant="primary" 
          icon="merge" 
          disabled={!protection.mergeable || merging} 
          onclick={doMerge}
        >
          {merging ? 'Fusionando...' : 'Fusionar Pull Request'}
        </Button>
      </footer>
    {/if}
  </div>
</Card>

<style>
  :global(.merge-box-card) { border-radius: 1.5rem; border: 2px solid var(--md-sys-color-outline-variant); overflow: hidden; margin-top: 2rem; }
  
  .merge-box { display: flex; flex-direction: column; background: rgba(var(--primary-rgb), 0.02); }
  
  .loading-merge { padding: 2rem; text-align: center; color: var(--text-muted); font-style: italic; }

  .merge-status-list { padding: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
  
  .status-item { display: flex; gap: 1.25rem; padding: 1.25rem; border-radius: 1rem; background: var(--surface-variant); border: 1px solid transparent; }
  
  .status-icon { flex-shrink: 0; padding-top: 0.1rem; }
  
  .status-text { display: flex; flex-direction: column; gap: 0.2rem; }
  .status-text strong { font-size: 1rem; }
  .status-text span { font-size: 0.85rem; color: var(--text-muted); }

  .status-item.success { background: rgba(16, 185, 129, 0.05); border-color: rgba(16, 185, 129, 0.2); color: #10b981; }
  .status-item.success span { color: #10b981; opacity: 0.8; }
  
  .status-item.error { background: rgba(var(--error-rgb), 0.05); border-color: rgba(var(--error-rgb), 0.2); color: var(--error); }
  .status-item.error span { color: var(--error); opacity: 0.8; }

  .status-item.warning { background: rgba(245, 158, 11, 0.05); border-color: rgba(245, 158, 11, 0.2); color: #f59e0b; }
  .status-item.warning span { color: #f59e0b; opacity: 0.8; }

  .merge-actions { padding: 1.5rem; background: var(--surface-variant); border-top: 1px solid var(--md-sys-color-outline-variant); display: flex; justify-content: flex-end; }
</style>
