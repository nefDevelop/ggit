<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Card from './ui/Card.svelte';
  import { i18n } from '$lib/i18n.svelte';
  
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

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="header-flex">
      <div class="title-with-icon">
        <span class="material-symbols-rounded hero-icon">monitoring</span>
        <div>
          <h2>Monitor del Sistema</h2>
          <p class="subtitle">Tareas programadas, procesos activos y recursos</p>
        </div>
      </div>
    </div>
  </header>

  {#if loading}
    <div class="skeleton-monitor"></div>
  {:else if monitorData}
    <div class="monitor-layout">
      <!-- CRON TASKS SECTION -->
      <section class="monitor-section">
        <div class="section-tag">
          <span class="material-symbols-rounded">schedule</span>
          Tareas Cron
        </div>
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
      </section>

      <!-- ACTIVE PROCESSES SECTION -->
      <section class="monitor-section">
        <div class="section-tag">
          <span class="material-symbols-rounded">account_tree</span>
          Procesos en Ejecución
        </div>
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
      </section>
    </div>
  {/if}
</Card>

<style>
  .section-header { margin-bottom: 4rem; }
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  @media (max-width: 600px) {
    .header-flex { flex-direction: column; align-items: stretch; text-align: center; }
    .title-with-icon { flex-direction: column; gap: 1rem; }
  }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; flex-shrink: 0; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  @media (max-width: 600px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { font-size: 0.9rem; }
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
  .section-tag span { font-size: 1.25rem; }

  .task-name { font-size: 1.05rem; color: var(--text-color); font-weight: 800; }
  .spec-code { font-family: 'JetBrains Mono', monospace; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); padding: 0.4rem 0.75rem; border-radius: 0.75rem; font-size: 0.85rem; font-weight: 700; }
  .time-stamp { color: var(--text-muted); font-size: 0.9rem; font-weight: 600; }

  .exec-pill { display: flex; flex-direction: column; align-items: center; background: rgba(255,255,255,0.03); padding: 0.4rem 1rem; border-radius: 1rem; border: 1px solid rgba(255,255,255,0.05); }
  .exec-pill .val { font-weight: 900; font-size: 1rem; color: var(--text-color); }
  .exec-pill .lab { font-size: 0.6rem; text-transform: uppercase; font-weight: 800; color: var(--text-muted); }

  .pid-cell code { font-family: 'JetBrains Mono', monospace; font-weight: 900; color: var(--primary-color); }
  .proc-desc { color: var(--text-color); font-weight: 700; max-width: 500px; }

  .skeleton-monitor { height: 600px; background: var(--surface-variant); border-radius: 2.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }
</style>
