<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  import Input from "./ui/Input.svelte";
  import Checkbox from "./ui/Checkbox.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  import { i18n } from "$lib/i18n.svelte";
  

import GitBranch from "lucide-svelte/icons/git-branch";
import Shield from "lucide-svelte/icons/shield";
import ShieldAlert from "lucide-svelte/icons/shield-alert";
import ShieldCheck from "lucide-svelte/icons/shield-check";
import Settings2 from "lucide-svelte/icons/settings-2";
import ArrowLeft from "lucide-svelte/icons/arrow-left";



  const repo = getContext<RepoContext>('repo');
  
  let branches = $state<any>({ default: null, active: [], stale: [] });
  let loading = $state(true);
  let selectedBranch = $state<string | null>(null);
  let protection = $state<any>(null);
  let loadingProtection = $state(false);
  let saving = $state(false);

  async function fetchBranches() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/branches`
      );
      if (res.ok) branches = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function editProtection(branchName: string) {
    selectedBranch = branchName;
    loadingProtection = true;
    protection = null;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/branches/${branchName}/protection`
      );
      if (res.ok) {
        protection = await res.json();
        if (!protection.id) {
           protection = {
             protected: false,
             require_pull_request: false,
             enable_whitelist: false,
             required_approvals: 0,
             enable_merge_whitelist: false,
             block_on_outdated_branch: false
           };
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      loadingProtection = false;
    }
  }

  async function saveProtection() {
    if (!selectedBranch) return;
    saving = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/branches/${selectedBranch}/protection`,
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(protection),
        }
      );
      if (res.ok) {
        toast.success(`Protección de rama "${selectedBranch}" actualizada.`);
        selectedBranch = null;
        fetchBranches();
      }
    } catch (e) {
      console.error(e);
      toast.error("Error al guardar la protección.");
    } finally {
      saving = false;
    }
  }

  onMount(fetchBranches);
</script>

{#if !selectedBranch}
  <Card padding="none" class="settings-premium-card">
    <PageHeader 
      title="Ramas del Repositorio"
      subtitle="Gestiona la rama principal y configura reglas de protección para evitar cambios no deseados."
      icon={GitBranch}
    />

    <div class="settings-content">
      {#if loading}
        <div class="skeleton-stack">
          {#each Array(4) as _}
            <div class="skeleton-item tonal"></div>
          {/each}
        </div>
      {:else}
        <div class="branches-stack">
          <!-- Rama por defecto -->
          {#if branches.default}
            <div class="branch-group">
              <span class="group-tag">Rama por Defecto</span>
              <div class="branch-tile default-branch tonal">
                <div class="branch-core">
                  <div class="branch-icon-wrap primary">
                    <GitBranch size={22} />
                  </div>
                  <div class="branch-info">
                    <strong class="name">{branches.default.name}</strong>
                    <div class="status-row">
                      {#if branches.default.is_protected}
                        <span class="protected-pill">
                          <ShieldCheck size={12} /> Protegida
                        </span>
{:else}
                        <span class="unprotected-pill">Sin protección activa</span>
                      {/if}
                    </div>
                  </div>
                </div>
                <Button variant="outline" size="sm" icon="settings" onclick={() => editProtection(branches.default.name)}>
                  Configurar Reglas
                </Button>
              </div>
            </div>
          {/if}

          <!-- Ramas activas -->
          {#if branches.active.length > 0}
            <div class="branch-group">
              <span class="group-tag">Otras Ramas Activas</span>
              <div class="tiles-grid">
                {#each branches.active as branch}
                  <div class="branch-tile tonal hover-lift">
                    <div class="branch-core">
                      <div class="branch-icon-wrap">
                        <GitBranch size={18} />
                      </div>
                      <div class="branch-info">
                        <strong class="name">{branch.name}</strong>
                        {#if branch.is_protected}
                          <span class="protected-pill mini">
                            <ShieldCheck size={10} /> Protegida
                          </span>
                        {/if}
                      </div>
                    </div>
                    <Button variant="ghost" size="sm" icon="settings" onclick={() => editProtection(branch.name)} />
                  </div>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      {/if}
    </div>
  </Card>
{:else}
  <div class="protection-view">
    <Button variant="ghost" icon="arrow_back" onclick={() => selectedBranch = null} class="back-link">
      Volver a la lista de ramas
    </Button>

    <Card padding="none" class="settings-premium-card protection-card">
      <PageHeader 
        title={`Reglas de Protección: ${selectedBranch}`}
        subtitle="Asegura la integridad de tu código restringiendo quién puede realizar cambios directos."
        icon={Shield}
      />

      <div class="settings-content">
        {#if loadingProtection}
          <div class="skeleton-form"></div>
        {:else if protection}
          <form class="protection-form" onsubmit={(e) => { e.preventDefault(); saveProtection(); }}>
            <div class="main-switch island">
              <Checkbox 
                label="Habilitar Protección de Rama" 
                description="Activa todas las restricciones configuradas abajo para esta rama específica." 
                bind:checked={protection.protected} 
              />
            </div>

            {#if protection.protected}
              <div class="protection-options-grid">
                <div class="option-tile tonal">
                  <Checkbox 
                    label="Requerir Pull Request" 
                    description="Nadie podrá hacer push directamente. Los cambios deben venir de una rama diferente." 
                    bind:checked={protection.require_pull_request} 
                  />
                </div>

                <div class="option-tile tonal">
                  <div class="option-header">
                    <Checkbox 
                      label="Requerir Aprobaciones" 
                      description="Define cuántas aprobaciones externas son necesarias antes de fusionar." 
                      bind:checked={protection.required_approvals_enabled} 
                    />
                  </div>
                  {#if protection.required_approvals_enabled}
                    <div class="option-nested">
                      <Input 
                        type="number" 
                        label="Número de aprobaciones" 
                        bind:value={protection.required_approvals} 
                        min="1"
                      />
                    </div>
                  {/if}
                </div>

                <div class="option-tile tonal">
                  <Checkbox 
                    label="Estado de Rama Actualizado" 
                    description="Impide el merge si la rama no tiene los últimos cambios de su base." 
                    bind:checked={protection.block_on_outdated_branch} 
                  />
                </div>

                <div class="option-tile tonal">
                  <Checkbox 
                    label="Restringir Fusión (Merge Whitelist)" 
                    description="Solo usuarios o equipos específicos podrán completar PRs en esta rama." 
                    bind:checked={protection.enable_merge_whitelist} 
                  />
                  {#if protection.enable_merge_whitelist}
                    <div class="option-nested">
                       <Input label="IDs de Usuarios" bind:value={protection.merge_whitelist_user_ids} placeholder="12, 45, 89..." />
                       <Input label="IDs de Equipos" bind:value={protection.merge_whitelist_team_ids} placeholder="1, 5..." />
                    </div>
                  {/if}
                </div>

                <div class="option-tile tonal">
                  <Checkbox 
                    label="Restringir Escritura (Push Whitelist)" 
                    description="Solo usuarios seleccionados podrán hacer push si el PR no es obligatorio." 
                    bind:checked={protection.enable_whitelist} 
                  />
                   {#if protection.enable_whitelist}
                    <div class="option-nested">
                       <Input label="IDs de Usuarios" bind:value={protection.whitelist_user_ids} />
                       <Input label="IDs de Equipos" bind:value={protection.whitelist_team_ids} />
                    </div>
                  {/if}
                </div>
              </div>
            {/if}

            <div class="form-actions-premium">
              <Button variant="primary" type="submit" loading={saving} disabled={saving}>
                Guardar Reglas de Protección
              </Button>
            </div>
          </form>
        {/if}
      </div>
    </Card>
  </div>
{/if}

<style>
  .settings-content { padding: 2.5rem; }
  @media (max-width: 600px) { .settings-content { padding: 1.5rem; } }

  .branches-stack { display: flex; flex-direction: column; gap: 3rem; }
  .branch-group { display: flex; flex-direction: column; gap: 1.25rem; }
  .group-tag { font-size: 0.75rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.1em; opacity: 0.8; padding-left: 0.5rem; }

  .branch-tile { display: flex; align-items: center; justify-content: space-between; padding: 1.25rem 1.75rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.2s; }
  .branch-tile:hover { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.02); }
  
  .default-branch { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.03); padding: 1.5rem 2rem; }

  .branch-core { display: flex; align-items: center; gap: 1.5rem; }
  .branch-icon-wrap { width: 44px; height: 44px; border-radius: 1rem; background: var(--md-sys-color-surface-container-highest); color: var(--text-muted); display: flex; align-items: center; justify-content: center; opacity: 0.6; }
  .branch-icon-wrap.primary { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); opacity: 1; }
  
  .branch-info { display: flex; flex-direction: column; gap: 0.25rem; }
  .name { font-weight: 900; font-size: 1.15rem; color: var(--text-color); font-family: 'JetBrains Mono', monospace; }
  
  .status-row { display: flex; align-items: center; gap: 0.75rem; }
  .protected-pill { display: inline-flex; align-items: center; gap: 0.4rem; background: color-mix(in srgb, var(--status-success) 10%, transparent); color: var(--status-success); padding: 0.2rem 0.7rem; border-radius: 100px; font-size: 0.7rem; font-weight: 900; text-transform: uppercase; }
  .protected-pill.mini { padding: 0.1rem 0.5rem; font-size: 0.6rem; }
  .unprotected-pill { font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; font-weight: 600; }

  .tiles-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(400px, 1fr)); gap: 1rem; }
  @media (max-width: 800px) { .tiles-grid { grid-template-columns: 1fr; } }

  /* PROTECTION EDITOR */
  .protection-view { display: flex; flex-direction: column; gap: 2rem; animation: slideIn 0.4s cubic-bezier(0.16, 1, 0.3, 1); }
  @keyframes slideIn { from { opacity: 0; transform: translateX(20px); } to { opacity: 1; transform: translateX(0); } }

  :global(.back-link) { align-self: flex-start; color: var(--primary-color); font-weight: 800; }

  .protection-form { display: flex; flex-direction: column; gap: 2.5rem; }
  .main-switch { padding: 2rem; border-radius: 1.5rem; border: 2px solid var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .protection-options-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; animation: fadeIn 0.3s ease-out; }
  @media (max-width: 900px) { .protection-options-grid { grid-template-columns: 1fr; } }

  .option-tile { padding: 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); display: flex; flex-direction: column; gap: 1rem; }
  .option-nested { margin-top: 1rem; padding-left: 1.5rem; border-left: 2px solid var(--md-sys-color-outline-variant); display: flex; flex-direction: column; gap: 1rem; }

  .form-actions-premium { display: flex; justify-content: flex-end; padding-top: 2rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  .skeleton-stack { display: flex; flex-direction: column; gap: 1.5rem; }
  .skeleton-item { height: 100px; border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }
  @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }

  @media (max-width: 600px) {
    .branch-tile { flex-direction: column; align-items: stretch; gap: 1.5rem; }
    .branch-core { gap: 1rem; }
    .branch-icon-wrap { width: 36px; height: 36px; }
    .name { font-size: 1rem; }
  }
</style>
