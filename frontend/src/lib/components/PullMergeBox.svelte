<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import { toast } from "$lib/toast.svelte";
  import { ShieldCheck, ShieldAlert, GitMerge, CheckCircle2, XCircle, AlertCircle } from "lucide-svelte";

  let { index, onMerged } = $props<{ index: number, onMerged: () => void }>();
  const repo: any = getContext("repo");

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

        <!-- Requisito de Cambios Solicitados -->
        {#if protection.rejections > 0}
          <div class="status-item error">
            <div class="status-icon"><XCircle size={20} /></div>
            <div class="status-text">
              <strong>Cambios solicitados pendientes</strong>
              <span>Debes resolver todos los comentarios de "Solicitud de cambios" antes de continuar.</span>
            </div>
          </div>
        {/if}

        <!-- Estado General -->
        {#if protection.mergeable}
          <div class="status-item success">
            <div class="status-icon"><GitMerge size={20} /></div>
            <div class="status-text">
              <strong>Lista para fusionar</strong>
              <span>Esta rama no tiene conflictos y cumple con las reglas de protección.</span>
            </div>
          </div>
        {:else}
          <div class="status-item warning">
            <div class="status-icon"><AlertCircle size={20} /></div>
            <div class="status-text">
              <strong>Fusión bloqueada</strong>
              <span>No se puede fusionar hasta que se cumplan todos los requisitos anteriores.</span>
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
  .merge-box-card { border-radius: 1.5rem; border: 2px solid var(--md-sys-color-outline-variant); overflow: hidden; margin-top: 2rem; }
  
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
  
  .status-item.error { background: rgba(239, 68, 68, 0.05); border-color: rgba(239, 68, 68, 0.2); color: #ef4444; }
  .status-item.error span { color: #ef4444; opacity: 0.8; }

  .status-item.warning { background: rgba(245, 158, 11, 0.05); border-color: rgba(245, 158, 11, 0.2); color: #f59e0b; }
  .status-item.warning span { color: #f59e0b; opacity: 0.8; }

  .merge-actions { padding: 1.5rem; background: var(--surface-variant); border-top: 1px solid var(--md-sys-color-outline-variant); display: flex; justify-content: flex-end; }
</style>
