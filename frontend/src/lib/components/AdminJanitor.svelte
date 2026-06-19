<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { tool } from "$lib/tool";
  import { i18n } from "$lib/i18n.svelte";
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  
import Sparkles from "lucide-svelte/icons/sparkles";
import Book from "lucide-svelte/icons/book";
import Trash2 from "lucide-svelte/icons/trash-2";
import Settings from "lucide-svelte/icons/settings";
import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
import Trash from "lucide-svelte/icons/trash";


  let loading = $state(true);
  let stats = $state<any[]>([]);

  async function fetchStats() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/admin/janitor");
      if (res.ok) {
        stats = res.data || [];
        stats.sort((a, b) => b.waste - a.waste);
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  onMount(fetchStats);
</script>

<div class="admin-janitor-container">
  <PageHeader 
    title={i18n.t("admin.janitor_global") || "Ggit Janitor"}
    subtitle={i18n.t("admin.janitor_global_desc") || "Limpieza y mantenimiento de espacio en disco de repositorios"}
    icon={Sparkles}
  />

  <div class="admin-janitor-body">

  {#if loading}
    <div class="skeleton-grid">
      {#each Array(3) as _}
        <div class="skeleton-card"></div>
      {/each}
    </div>
  {:else if stats.length === 0}
    <Card class="empty-card">
      <CheckCircle2 size={48} />
      <p>{i18n.t("admin.empty_janitor")}</p>
    </Card>
  {:else}
    <div class="stats-grid">
      {#each stats as repo}
        <Card class="repo-stat-card" tonal={repo.waste > 0}>
          <div class="repo-info">
            <div class="repo-main">
              <Book size={24} class="repo-icon" />
              <div class="repo-names">
                <span class="owner">{repo.owner} /</span>
                <span class="name">{repo.name}</span>
              </div>
            </div>
            <div class="waste-badge" class:warning={repo.waste > 0}>
              <Trash size={18} />
              {tool.fileSize(repo.waste)}
            </div>
          </div>

          <div class="repo-details">
            <div class="detail-item">
              <span class="label">{i18n.t("admin.issues")}</span>
              <span class="value">{repo.issues_count}</span>
            </div>
            <div class="detail-item">
              <span class="label">{i18n.t("admin.last_scan")}</span>
              <span class="value">
                {repo.last_check > 0 
                  ? new Date(repo.last_check * 1000).toLocaleDateString() 
                  : (i18n.t("admin.never") || "Never")}
              </span>
            </div>
          </div>

          <div class="card-actions">
            <Button href="/{repo.owner}/{repo.name}/settings" variant="outline" size="sm" icon={Settings}>
              {i18n.t("admin.manage")}
            </Button>
          </div>
        </Card>
      {/each}
    </div>
  {/if}
  </div>
</div>

<style>
  .admin-janitor-body {
    padding: 2rem 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  :global(.admin-janitor) { display: flex; flex-direction: column; gap: 3rem; }
  
  .stats-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(400px, 1fr)); gap: 2rem; }
  
  :global(.repo-stat-card) { 
    display: flex; 
    flex-direction: column; 
    gap: 2rem !important; 
    padding: 2rem !important;
    border-radius: 1.25rem !important;
  }
  
  .repo-info { display: flex; justify-content: space-between; align-items: flex-start; gap: 1rem; }
  .repo-main { display: flex; gap: 1.25rem; align-items: center; }
  .repo-names { display: flex; flex-direction: column; gap: 0.15rem; }
  .owner { font-size: 0.85rem; color: var(--md-sys-color-on-surface-variant); font-weight: 700; opacity: 0.7; }
  .name { font-size: 1.35rem; font-weight: 900; color: var(--md-sys-color-on-surface); letter-spacing: -0.01em; }
  .waste-badge { 
    display: flex; 
    align-items: center; 
    gap: 0.5rem; 
    padding: 0.75rem 1.25rem; 
    border-radius: 0.5rem; 
    background: var(--md-sys-color-surface-container-high); 
    font-weight: 900; 
    font-size: 0.95rem;
    white-space: nowrap;
  }
  .waste-badge.warning { background: color-mix(in srgb, var(--md-sys-color-error), transparent 90%); color: var(--md-sys-color-error); }
  .repo-details { 
    display: grid; 
    grid-template-columns: 1fr 1fr; 
    gap: 1.5rem; 
    padding: 1.5rem; 
    background: var(--md-sys-color-surface-container-low); 
    border-radius: 0.75rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  .detail-item { display: flex; flex-direction: column; gap: 0.4rem; }
  .label { font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.08em; color: var(--md-sys-color-on-surface-variant); font-weight: 800; opacity: 0.8; }
  .value { font-size: 1.1rem; font-weight: 700; color: var(--md-sys-color-on-surface); }
  .card-actions { display: flex; justify-content: flex-end; margin-top: 0.5rem; }
  :global(.card-actions .btn) { border-radius: 0.5rem !important; font-weight: 700 !important; }
  :global(.empty-card) { padding: 4rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  :global(.empty-card span) { font-size: 4rem; opacity: 0.2; }
  .skeleton-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(400px, 1fr)); gap: 2rem; }
  .skeleton-card { height: 260px; background: var(--md-sys-color-surface-container-high); border-radius: 1.25rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 1; } 50% { opacity: 0.5; } 100% { opacity: 1; } }
</style>
