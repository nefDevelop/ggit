<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";
  import Input from "./ui/Input.svelte";
  import Checkbox from "./ui/Checkbox.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";
  import { GitBranch, Shield, ShieldAlert, ShieldCheck, Settings2, ArrowLeft } from "lucide-svelte";

  const repo: any = getContext("repo");
  
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
           // Default values for new protection
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

<div class="branches-settings">
  {#if !selectedBranch}
    <Card padding="lg">
      <header class="section-header">
        <div class="title-info">
          <h3>Ramas del Repositorio</h3>
          <p>Gestiona y protege las ramas de tu proyecto.</p>
        </div>
      </header>

      <div class="branches-list">
        {#if loading}
          <div class="skeleton-list"></div>
        {:else}
          <!-- Rama por defecto -->
          {#if branches.default}
            <div class="branch-group">
              <span class="group-label">Rama por Defecto</span>
              <div class="branch-item default tonal">
                <div class="branch-info">
                  <div class="branch-icon main">
                    <GitBranch size={20} />
                  </div>
                  <div class="branch-meta">
                    <span class="branch-name">{branches.default.name}</span>
                    <span class="branch-status">
                      {#if branches.default.is_protected}
                        <span class="protected-badge">
                          <ShieldCheck size={12} /> Protegida
                        </span>
                      {:else}
                        <span class="unprotected-hint">Sin protección</span>
                      {/if}
                    </span>
                  </div>
                </div>
                <Button variant="ghost" icon="settings" onclick={() => editProtection(branches.default.name)}>
                  Configurar
                </Button>
              </div>
            </div>
          {/if}

          <!-- Ramas activas -->
          {#if branches.active.length > 0}
            <div class="branch-group">
              <span class="group-label">Ramas Activas</span>
              <div class="list-container">
                {#each branches.active as branch}
                  <div class="branch-item tonal">
                    <div class="branch-info">
                      <div class="branch-icon">
                        <GitBranch size={18} />
                      </div>
                      <div class="branch-meta">
                        <span class="branch-name">{branch.name}</span>
                        {#if branch.is_protected}
                          <span class="protected-badge mini">
                            <ShieldCheck size={10} /> Protegida
                          </span>
                        {/if}
                      </div>
                    </div>
                    <Button variant="ghost" icon="settings" onclick={() => editProtection(branch.name)} />
                  </div>
                {/each}
              </div>
            </div>
          {/if}
        {/if}
      </div>
    </Card>
  {:else}
    <div class="protection-editor-view">
      <button class="back-btn" onclick={() => selectedBranch = null}>
        <ArrowLeft size={20} /> Volver a ramas
      </button>

      <Card padding="lg">
        <header class="section-header">
          <div class="title-with-icon">
            <div class="icon-wrap" class:active={protection?.protected}>
               <Shield size={32} />
            </div>
            <div>
              <h3>Protección de Rama: {selectedBranch}</h3>
              <p>Define reglas para evitar cambios accidentales o no autorizados.</p>
            </div>
          </div>
        </header>

        {#if loadingProtection}
          <div class="skeleton-form"></div>
        {:else if protection}
          <form class="protection-form" onsubmit={(e) => { e.preventDefault(); saveProtection(); }}>
            <div class="main-toggle">
              <Checkbox 
                label="Habilitar Protección de Rama" 
                description="Activa las restricciones para esta rama." 
                bind:checked={protection.protected} 
              />
            </div>

            {#if protection.protected}
              <div class="protection-options">
                <div class="option-card tonal">
                  <Checkbox 
                    label="Requerir Pull Request" 
                    description="Prohíbe los pushes directos. Todos los cambios deben hacerse vía PR." 
                    bind:checked={protection.require_pull_request} 
                  />
                </div>

                <div class="option-card tonal">
                  <header class="sub-option-header">
                    <Checkbox 
                      label="Requerir Aprobaciones" 
                      description="Número mínimo de aprobaciones necesarias para fusionar." 
                      bind:checked={protection.required_approvals_enabled} 
                    />
                  </header>
                  <div class="sub-content">
                    <Input 
                      type="number" 
                      label="Aprobaciones requeridas" 
                      bind:value={protection.required_approvals} 
                    />
                  </div>
                </div>

                <div class="option-card tonal">
                  <Checkbox 
                    label="Bloquear en ramas desactualizadas" 
                    description="Impide el merge si la rama no está sincronizada con la base." 
                    bind:checked={protection.block_on_outdated_branch} 
                  />
                </div>

                <div class="option-card tonal">
                  <Checkbox 
                    label="Restringir quién puede fusionar (Merge Whitelist)" 
                    description="Solo usuarios o equipos específicos podrán completar el PR." 
                    bind:checked={protection.enable_merge_whitelist} 
                  />
                  {#if protection.enable_merge_whitelist}
                    <div class="sub-content">
                       <Input label="IDs de Usuarios (separados por coma)" bind:value={protection.merge_whitelist_user_ids} />
                       <Input label="IDs de Equipos (separados por coma)" bind:value={protection.merge_whitelist_team_ids} />
                    </div>
                  {/if}
                </div>

                <div class="option-card tonal">
                  <Checkbox 
                    label="Restringir Pushes (Push Whitelist)" 
                    description="Solo usuarios o equipos específicos podrán hacer push (si no se requiere PR)." 
                    bind:checked={protection.enable_whitelist} 
                  />
                   {#if protection.enable_whitelist}
                    <div class="sub-content">
                       <Input label="IDs de Usuarios (separados por coma)" bind:value={protection.whitelist_user_ids} />
                       <Input label="IDs de Equipos (separados por coma)" bind:value={protection.whitelist_team_ids} />
                    </div>
                  {/if}
                </div>
              </div>
            {/if}

            <div class="form-actions">
              <Button variant="primary" type="submit" disabled={saving}>
                {saving ? "Guardando..." : "Guardar Cambios"}
              </Button>
            </div>
          </form>
        {/if}
      </Card>
    </div>
  {/if}
</div>

<style>
  .branches-settings { display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header { margin-bottom: 2.5rem; }
  .section-header h3 { margin: 0; font-size: 1.5rem; font-weight: 900; letter-spacing: -0.02em; }
  .section-header p { margin: 0.5rem 0 0; color: var(--text-muted); font-size: 1rem; }

  .branches-list { display: flex; flex-direction: column; gap: 2.5rem; }
  .branch-group { display: flex; flex-direction: column; gap: 1rem; }
  .group-label { font-size: 0.75rem; font-weight: 900; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.1em; opacity: 0.8; padding-left: 0.5rem; }

  .branch-item { display: flex; align-items: center; justify-content: space-between; padding: 1.25rem 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.2s; }
  .branch-item:hover { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  
  .branch-item.default { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.03); }

  .branch-info { display: flex; align-items: center; gap: 1.5rem; }
  .branch-icon { width: 44px; height: 44px; border-radius: 1rem; background: var(--surface-variant); color: var(--text-muted); display: flex; align-items: center; justify-content: center; }
  .branch-icon.main { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); }
  
  .branch-meta { display: flex; flex-direction: column; gap: 0.2rem; }
  .branch-name { font-weight: 800; font-size: 1.1rem; }
  
  .protected-badge { display: inline-flex; align-items: center; gap: 0.4rem; background: rgba(16, 185, 129, 0.1); color: #10b981; padding: 0.2rem 0.6rem; border-radius: 2rem; font-size: 0.75rem; font-weight: 800; }
  .protected-badge.mini { padding: 0.1rem 0.4rem; font-size: 0.65rem; }
  .unprotected-hint { font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; }

  .list-container { display: flex; flex-direction: column; gap: 0.75rem; }

  /* Editor View */
  .protection-editor-view { display: flex; flex-direction: column; gap: 1.5rem; }
  .back-btn { background: none; border: none; color: var(--primary-color); font-weight: 800; display: flex; align-items: center; gap: 0.5rem; cursor: pointer; padding: 0.5rem; border-radius: 1rem; align-self: flex-start; }
  .back-btn:hover { background: rgba(var(--primary-rgb), 0.1); }

  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .icon-wrap { width: 64px; height: 64px; border-radius: 1.5rem; background: var(--surface-variant); color: var(--text-muted); display: flex; align-items: center; justify-content: center; }
  .icon-wrap.active { background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); }

  .protection-form { display: flex; flex-direction: column; gap: 2rem; margin-top: 2rem; }
  .main-toggle { padding: 1.5rem; background: rgba(var(--primary-rgb), 0.05); border-radius: 1.5rem; border: 1px solid var(--primary-color); }

  .protection-options { display: flex; flex-direction: column; gap: 1rem; animation: slideDown 0.3s ease-out; }
  @keyframes slideDown { from { opacity: 0; transform: translateY(-10px); } to { opacity: 1; transform: translateY(0); } }

  .option-card { padding: 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .sub-content { margin-top: 1.5rem; padding-left: 2rem; border-left: 2px solid var(--md-sys-color-outline-variant); display: flex; flex-direction: column; gap: 1rem; }

  .form-actions { margin-top: 1rem; display: flex; justify-content: flex-end; }

  .skeleton-list { height: 300px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }
</style>
