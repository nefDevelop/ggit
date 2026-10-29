<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { toast } from "$lib/toast.svelte";
  import Button from "./ui/Button.svelte";
  import Badge from "./ui/Badge.svelte";
  import Card from "./ui/Card.svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  

import Sparkles from "lucide-svelte/icons/sparkles";
import Search from "lucide-svelte/icons/search";
import Trash2 from "lucide-svelte/icons/trash-2";


  import { tool } from "$lib/tool";
  import { i18n } from "$lib/i18n.svelte";

  const repo = getContext<RepoContext>('repo');
  let loading = $state(true);
  let scanning = $state(false);
  let cleaning = $state(false);
  let report = $state<any>(null);
  let selectedPaths = $state<string[]>([]);

  async function fetchReport() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/janitor/report`);
      if (res.ok) {
        report = await res.json();
        selectedPaths = report.issues?.map((i: any) => i.path) || [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function runScan() {
    scanning = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/janitor/scan`, {
        method: "POST"
      });
      if (res.ok) {
        report = await res.json();
        selectedPaths = report.issues?.map((i: any) => i.path) || [];
        toast.success("Scan complete!");
      }
    } catch (e) {
      console.error(e);
    } finally {
      scanning = false;
    }
  }

  async function runCleanup() {
    if (selectedPaths.length === 0) return;
    cleaning = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/janitor/cleanup`, {
        method: "POST",
        body: JSON.stringify({ paths: selectedPaths })
      });
      if (res.ok) {
        toast.success("Cleanup successful! A new commit has been created.");
        report = null;
        selectedPaths = [];
        fetchReport();
      }
    } catch (e) {
      console.error(e);
      toast.error("Failed to perform cleanup.");
    } finally {
      cleaning = false;
    }
  }

  function toggleSelect(path: string) {
    if (selectedPaths.includes(path)) {
      selectedPaths = selectedPaths.filter(p => p !== path);
    } else {
      selectedPaths = [...selectedPaths, path];
    }
  }

  onMount(fetchReport);
  import BadgeCheck from "lucide-svelte/icons/badge-check";
  import Package from "lucide-svelte/icons/package";
  import Hammer from "lucide-svelte/icons/hammer";
  import Monitor from "lucide-svelte/icons/monitor";
  import FileCheck from "lucide-svelte/icons/file-check";
  import WandSparkles from "lucide-svelte/icons/wand-sparkles";
  import CircleCheck from "lucide-svelte/icons/circle-check";
  import Circle from "lucide-svelte/icons/circle";
</script>

<Card padding="none" class="settings-premium-card">
  <PageHeader 
    title="Ggit Janitor" 
    subtitle={i18n.t("settings.janitor_desc")} 
    icon={Sparkles}
  >
    {#snippet actions()}
      <Button variant="outline" icon={Search} onclick={runScan} loading={scanning} disabled={cleaning}>
        {report ? i18n.t("settings.re_scan") : i18n.t("settings.scan_now")}
      </Button>
      {#if report && report.issues?.length > 0}
        <Button variant="primary" icon={Trash2} onclick={runCleanup} loading={cleaning} disabled={scanning || selectedPaths.length === 0}>
          {i18n.t("settings.clean_selected")} ({selectedPaths.length})
        </Button>
      {/if}
    {/snippet}
  </PageHeader>

  <div class="settings-content">
    {#if loading}
      <div class="skeleton-form"></div>
    {:else if !report || !report.issues || report.issues.length === 0}
      <div class="empty-state-janitor island">
        <BadgeCheck size={18} />
        <h3>{i18n.t("settings.everything_clean")}</h3>
        <p>{i18n.t("settings.no_waste_detected")}</p>
        <Button variant="outline" size="sm" onclick={runScan}>{i18n.t("settings.repeat_scan")}</Button>
      </div>
    {:else}
      <div class="janitor-stack">
        <div class="summary-island island">
          <div class="stat-item">
            <span class="stat-value">{report.issues.length}</span>
            <span class="stat-label">{i18n.t("settings.issues_detected")}</span>
          </div>
          {#if report.total_waste > 0}
            <div class="stat-item">
              <span class="stat-value">{tool.fileSize(report.total_waste)}</span>
              <span class="stat-label">{i18n.t("settings.recoverable_space")}</span>
            </div>
          {/if}
          <div class="stat-item">
            <span class="stat-value">{new Date(report.generated_at * 1000).toLocaleDateString()}</span>
            <span class="stat-label">{i18n.t("settings.analysis_date")}</span>
          </div>
        </div>

        <div class="issues-grid">
          {#each report.issues as issue}
            <div 
              class="issue-tile tonal" 
              class:critical={issue.severity === 'critical'} 
              class:selected={selectedPaths.includes(issue.path)} 
              onclick={() => toggleSelect(issue.path)}
              onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); toggleSelect(issue.path); } }}
              role="button"
              tabindex="0"
            >
              <div class="issue-selector">
                {#if selectedPaths.includes(issue.path)}
                  <CircleCheck size={18} />
                {:else}
                  <Circle size={18} />
                {/if}
              </div>
              <div class="issue-icon">
                {#if issue.type === 'dependency'}
                  <Package size={18} />
                {:else if issue.type === 'build'}
                  <Hammer size={18} />
                {:else if issue.type === 'system'}
                  <Monitor size={18} />
                {:else}
                  <FileCheck size={18} />
                {/if}
              </div>
              <div class="issue-body">
                <div class="issue-header-row">
                  <code class="issue-path">{issue.path}</code>
                  <Badge variant={issue.severity === 'critical' ? 'danger' : 'warning'}>{issue.type}</Badge>
                </div>
                <p class="issue-description">{issue.description}</p>
              </div>
              {#if issue.size > 0}
                <div class="issue-meta-size">{tool.fileSize(issue.size)}</div>
              {/if}
            </div>
          {/each}
        </div>

        <div class="cleanup-hint island">
          <WandSparkles size={18} />
          <p>{@html i18n.t("settings.cleanup_hint")}</p>
        </div>
      </div>
    {/if}
  </div>
</Card>

<style>
  .settings-content { padding: 2.5rem; }
  @media (max-width: 600px) { .settings-content { padding: 1.5rem; } }

  .janitor-stack { display: flex; flex-direction: column; gap: 2.5rem; }

  .summary-island { display: flex; justify-content: space-around; padding: 2rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .stat-item { display: flex; flex-direction: column; align-items: center; text-align: center; gap: 0.25rem; }
  .stat-value { font-size: 1.75rem; font-weight: 950; color: var(--primary-color); }
  .stat-label { font-size: 0.75rem; font-weight: 900; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.1em; }

  @media (max-width: 600px) { .summary-island { flex-direction: column; gap: 2rem; } }

  .issues-grid { display: flex; flex-direction: column; gap: 1rem; }
  .issue-tile { display: flex; align-items: center; gap: 1.5rem; padding: 1.25rem 1.75rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); cursor: pointer; transition: all 0.2s; position: relative; }
  .issue-tile:hover { transform: translateX(4px); border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.02); }
  .issue-tile.selected { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  .issue-tile.critical { border-left: 6px solid var(--md-sys-color-error); }
  
  .issue-selector { color: var(--primary-color); transition: all 0.2s; }
  .issue-tile.selected .issue-selector { transform: scale(1.2); }

  .issue-icon { width: 44px; height: 44px; border-radius: 1rem; background: var(--md-sys-color-surface-container-highest); color: var(--primary-color); display: flex; align-items: center; justify-content: center; flex-shrink: 0; opacity: 0.7; }
  .issue-body { flex: 1; min-width: 0; }
  .issue-header-row { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.25rem; }
  .issue-path { font-family: 'JetBrains Mono', monospace; font-weight: 800; font-size: 1rem; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .issue-description { margin: 0; font-size: 0.9rem; color: var(--text-muted); opacity: 0.8; }
  .issue-meta-size { font-weight: 900; color: var(--text-muted); font-size: 0.85rem; opacity: 0.6; }

  .empty-state-janitor { padding: 5rem 2rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1.25rem; }
  .empty-state-janitor h3 { margin: 0; font-size: 1.75rem; font-weight: 950; }
  .empty-state-janitor p { color: var(--text-muted); font-size: 1.1rem; max-width: 450px; }

  .cleanup-hint { display: flex; align-items: center; gap: 1.25rem; padding: 1.5rem 2rem; border-radius: 1.5rem; background: rgba(var(--primary-rgb), 0.05); color: var(--text-muted); font-size: 0.95rem; border: 1px solid rgba(var(--primary-rgb), 0.1); }

  .skeleton-form { height: 400px; background: var(--md-sys-color-surface-variant); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }

  @media (max-width: 800px) {
    .issue-tile { flex-direction: column; align-items: flex-start; gap: 1.25rem; padding: 1.5rem; }
    .issue-meta-size { align-self: flex-end; }
  }
</style>
