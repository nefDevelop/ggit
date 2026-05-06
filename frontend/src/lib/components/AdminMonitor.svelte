<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Card from './ui/Card.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  import { Activity, Clock, CalendarOff, Network, Hourglass } from 'lucide-svelte';
  
  let monitorData = $state<any>(null);
  let loading = $state(true);

  async function fetchMonitor() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/admin/monitor');
      if (res.ok) monitorData = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
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
    title={i18n.t('admin.monitor') || "Monitor del Sistema"}
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
        {#if monitorData.cron && monitorData.cron.length > 0}
          <div class="modern-table-card">
            <table class="modern-table">
              <thead>
                <tr>
                  <th>Descripción</th>
                  <th>Intervalo</th>
                  <th>Próxima Ejecución</th>
                  <th>Estado</th>
                </tr>
              </thead>
              <tbody>
                {#each monitorData.cron as task}
                  <tr class="hover-row">
                    <td data-label="Descripción"><strong class="task-name">{task.Description}</strong></td>
                    <td data-label="Intervalo"><code class="spec-code">{task.Spec}</code></td>
                    <td data-label="Próxima Ejecución"><span class="time-stamp">{formatTime(task.Next)}</span></td>
                    <td data-label="Estado">
                      <div class="exec-pill">
                        <span class="val">{task.ExecTimes}</span>
                        <span class="lab">Ejecs.</span>
                      </div>
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
                </tr>
              </thead>
              <tbody>
                {#each monitorData.processes as proc}
                  <tr class="hover-row">
                    <td class="pid-cell" data-label="PID"><code>{proc.PID}</code></td>
                    <td class="proc-desc" data-label="Descripción del Proceso">{proc.Description}</td>
                    <td data-label="Iniciado el"><span class="time-stamp">{formatTime(proc.Start)}</span></td>
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
    </div>
  {/if}
  </div>
</div>

<style>
  .admin-monitor-body {
    padding: 0 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

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

  .skeleton-monitor { height: 600px; background: var(--surface-variant); border-radius: 1.25rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }

  .empty-monitor-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 4rem; background: rgba(0,0,0,0.1); border-radius: 1.25rem; border: 1px dashed var(--md-sys-color-outline-variant); color: var(--text-muted); gap: 1rem; text-align: center; }
  .empty-monitor-state p { margin: 0; font-size: 1.1rem; font-weight: 600; }
</style>
