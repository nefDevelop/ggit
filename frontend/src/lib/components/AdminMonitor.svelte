<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Card from './ui/Card.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  import { Activity, Clock, CalendarOff, Network, Hourglass, Play, Trash2 } from 'lucide-svelte';
  import { toast } from "$lib/toast.svelte";
  import Button from './ui/Button.svelte';
  
  let monitorData = $state<any>(null);
  let rawResponse = $state<string>("");
  let loading = $state(true);
  let debugInfo = $derived(monitorData?.debug_info);
  let cronList = $derived(monitorData?.cron_v2 || monitorData?.cron || []);

  async function fetchMonitor() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/admin/health-snapshot');
      rawResponse = (await res.text()).trim();
      
      console.log(`[v5] Raw Length: ${rawResponse.length}`);
      console.log(`[v5] Raw Start: ${rawResponse.substring(0, 50)}`);

      if (res.ok) {
        try {
          monitorData = JSON.parse(rawResponse);
        } catch (e) {
          console.error("Failed to parse v5 JSON", e);
          monitorData = { error: 'JSON parse failed', preview: rawResponse.substring(0, 100) };
        }
      } else {
        monitorData = { error: `Status ${res.status}`, url: res.url };
      }
    } catch (e) {
      console.error(e);
      monitorData = { error: 'Fetch failed' };
    } finally {
      loading = false;
    }
  }

  async function runTask(description: string) {
    try {
      const res = await apiFetch('/api/v2/admin/monitor/cron', {
        method: 'POST',
        body: JSON.stringify({ description })
      });
      if (res.ok) {
        toast.success(i18n.t('admin.monitor.task_triggered') || "Tarea iniciada correctamente");
        fetchMonitor();
      } else {
        toast.error(i18n.t('admin.monitor.task_failed') || "Error al iniciar la tarea");
      }
    } catch (e) {
      console.error(e);
      toast.error("Error de red");
    }
  }

  async function killProcess(pid: number) {
    try {
      const res = await apiFetch(`/api/v2/admin/monitor/processes/${pid}`, {
        method: 'POST'
      });
      if (res.ok) {
        toast.success(i18n.t('admin.monitor.process_killed') || "Proceso detenido");
        fetchMonitor();
      } else {
        toast.error(i18n.t('admin.monitor.kill_failed') || "Error al detener el proceso");
      }
    } catch (e) {
      console.error(e);
      toast.error("Error de red");
    }
  }

  function formatTime(dateStr: string) {
    if (!dateStr || dateStr.startsWith('0001')) return '—';
    const date = new Date(dateStr);
    return date.toLocaleString();
  }

  onMount(fetchMonitor);
</script>

<div class="admin-monitor-container">
  <PageHeader 
    title={(i18n.t('admin.monitor') || "Monitor del Sistema") + " (v2)"}
    subtitle="Tareas programadas, procesos activos y recursos"
    icon={Activity}
  />

  <div class="admin-monitor-body">

  {#if loading}
    <div class="skeleton-monitor"></div>
  {:else if monitorData}
    <div class="monitor-layout">
      <!-- CRON TASKS SECTION -->
      <section class="monitor-section">
        <div class="section-tag">
          <Clock size={18} />
          Tareas Cron
        </div>
        {#if cronList && cronList.length > 0}
          <div class="modern-table-card">
            <table class="modern-table">
              <thead>
                <tr>
                  <th>Descripción</th>
                  <th>Intervalo</th>
                  <th>Próxima Ejecución</th>
                  <th>Estado</th>
                  <th class="actions-head">Acción</th>
                </tr>
              </thead>
              <tbody>
                {#each cronList as task}
                  <tr class="hover-row">
                    <td data-label="Descripción"><strong class="task-name">{task.Description}</strong></td>
                    <td data-label="Intervalo"><code class="spec-code">{task.Spec}</code></td>
                    <td data-label="Próxima Ejecución"><span class="time-stamp">{formatTime(task.Next)}</span></td>
                    <td data-label="Estado">
                      {#if task.IsRunning}
                        <div class="status-pill running">
                          <span class="pulse-dot"></span>
                          Ejecutando
                        </div>
                      {:else}
                        <div class="exec-pill">
                          <span class="val">{task.ExecTimes}</span>
                          <span class="lab">Ejecs.</span>
                        </div>
                      {/if}
                    </td>
                    <td class="actions-cell">
                      <Button variant="ghost" size="sm" icon={Play} onclick={() => runTask(task.Description)}>
                        Ejecutar
                      </Button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else}
          <div class="empty-monitor-state">
            <CalendarOff size={48} />
            <p>No hay tareas programadas activas en este momento.</p>
          </div>
        {/if}
      </section>

      <!-- ACTIVE PROCESSES SECTION -->
      <section class="monitor-section">
        <div class="section-tag">
          <Network size={18} />
          Procesos en Ejecución
        </div>
        {#if monitorData.processes && monitorData.processes.length > 0}
          <div class="modern-table-card">
            <table class="modern-table">
              <thead>
                <tr>
                  <th>PID</th>
                  <th>Descripción del Proceso</th>
                  <th>Iniciado el</th>
                  <th class="actions-head">Acción</th>
                </tr>
              </thead>
              <tbody>
                {#each monitorData.processes as proc}
                  <tr class="hover-row">
                    <td class="pid-cell" data-label="PID"><code>{proc.PID}</code></td>
                    <td class="proc-desc" data-label="Descripción del Proceso">{proc.Description}</td>
                    <td data-label="Iniciado el"><span class="time-stamp">{formatTime(proc.Start)}</span></td>
                    <td class="actions-cell">
                      <Button variant="danger" size="sm" icon={Trash2} onclick={() => killProcess(proc.PID)}>
                        Detener
                      </Button>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else}
          <div class="empty-monitor-state">
            <Hourglass size={48} />
            <p>No hay procesos en segundo plano activos actualmente.</p>
          </div>
        {/if}
      </section>
  {/if}
  </div>
</div>

<style>
  .admin-monitor-body {
    padding: 2rem 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .debug-panel {
    margin-top: 4rem;
    padding: 1.5rem;
    background: rgba(0,0,0,0.2);
    border-radius: 1rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    font-family: monospace;
  }
  .debug-panel h3 { margin-top: 0; font-size: 0.9rem; color: var(--text-muted); text-transform: uppercase; }
  .debug-panel pre { margin: 0; font-size: 0.8rem; color: var(--primary-color); }
  .raw-box { margin-top: 1rem !important; color: #fff !important; background: #000; padding: 1rem; border-radius: 0.5rem; overflow-x: auto; }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  @media (max-width: 600px) {
    .modern-table-card { border-radius: 1.5rem; }
  }

  @media (max-width: 1024px) {
    .modern-table-card { border-radius: 1.5rem; }
    .modern-table td { text-align: right; }
    .task-name, .spec-code, .time-stamp, .exec-pill, .pid-cell code, .proc-desc { text-align: right; }
    .exec-pill { align-items: flex-end; margin-left: auto; }
    .pid-cell code, .spec-code { display: inline-block; }
    .actions-head { display: none; }
    .actions-cell { justify-content: flex-end; }
  }

  .monitor-layout { display: flex; flex-direction: column; gap: 4rem; }
  .monitor-section { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .section-tag { font-size: 0.75rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }

  .task-name { font-size: 1.05rem; color: var(--text-color); font-weight: 800; }
  .spec-code { font-family: 'JetBrains Mono', monospace; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.4rem 0.75rem; border-radius: 0.75rem; font-size: 0.85rem; font-weight: 700; }
  .time-stamp { color: var(--text-muted); font-size: 0.9rem; font-weight: 600; }

  .exec-pill { display: flex; flex-direction: column; align-items: center; background: rgba(255,255,255,0.03); padding: 0.4rem 1rem; border-radius: 0.5rem; border: 1px solid rgba(255,255,255,0.05); }
  .exec-pill .val { font-weight: 900; font-size: 1rem; color: var(--text-color); }
  .exec-pill .lab { font-size: 0.6rem; text-transform: uppercase; font-weight: 800; color: var(--text-muted); }

  .pid-cell code { font-family: 'JetBrains Mono', monospace; font-weight: 900; color: var(--primary-color); }
  .proc-desc { color: var(--text-color); font-weight: 700; max-width: 500px; }

  .actions-head { text-align: center; width: 120px; }
  .actions-cell { text-align: center; }

  .skeleton-monitor { height: 600px; background: var(--surface-variant); border-radius: 1.25rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }

  .empty-monitor-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 4rem; background: rgba(0,0,0,0.1); border-radius: 1.25rem; border: 1px dashed var(--md-sys-color-outline-variant); color: var(--text-muted); gap: 1rem; text-align: center; }
  .empty-monitor-state p { margin: 0; font-size: 1.1rem; font-weight: 600; }

  /* Running status pill */
  .status-pill {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.35rem 0.85rem;
    border-radius: 999px;
    font-size: 0.8rem;
    font-weight: 800;
  }
  .status-pill.running {
    background: rgba(34, 197, 94, 0.15);
    color: #22c55e;
    border: 1px solid rgba(34, 197, 94, 0.3);
  }
  .pulse-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #22c55e;
    animation: pulse-green 1.5s infinite;
  }
  @keyframes pulse-green {
    0%, 100% { opacity: 1; box-shadow: 0 0 0 0 rgba(34,197,94,0.4); }
    50% { opacity: 0.7; box-shadow: 0 0 0 5px rgba(34,197,94,0); }
  }
</style>
