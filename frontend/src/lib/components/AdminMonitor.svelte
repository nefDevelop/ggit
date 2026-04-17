<script lang="ts">
  import { onMount } from 'svelte';
  
  let monitorData = $state<any>(null);
  let loading = $state(true);

  async function fetchMonitor() {
    loading = true;
    try {
      const res = await fetch('/api/v2/admin/monitor');
      if (res.ok) monitorData = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  function formatTime(dateStr: string) {
    const date = new Date(dateStr);
    return date.toLocaleString();
  }

  onMount(fetchMonitor);
</script>

<div class="admin-monitor">
  <header class="section-header">
    <h2>System Monitor</h2>
    <p class="subtitle">Monitor background tasks, cron jobs, and active processes.</p>
  </header>

  {#if loading}
    <div class="skeleton-monitor"></div>
  {:else if monitorData}
    <section class="island card">
      <h3 class="card-title">Cron Tasks</h3>
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>Description</th>
              <th>Schedule</th>
              <th>Next Run</th>
              <th>Prev Run</th>
              <th>Exec Times</th>
            </tr>
          </thead>
          <tbody>
            {#each monitorData.cron as task}
              <tr>
                <td><strong>{task.Description}</strong></td>
                <td><code>{task.Spec}</code></td>
                <td>{formatTime(task.Next)}</td>
                <td>{task.Prev ? formatTime(task.Prev) : 'N/A'}</td>
                <td>{task.ExecTimes}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </section>

    <section class="island card">
      <h3 class="card-title">Active Processes</h3>
      <div class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>PID</th>
              <th>Description</th>
              <th>Started At</th>
            </tr>
          </thead>
          <tbody>
            {#each monitorData.processes as proc}
              <tr>
                <td><code>{proc.PID}</code></td>
                <td>{proc.Description}</td>
                <td>{formatTime(proc.Start)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </section>
  {/if}
</div>

<style>
  .admin-monitor { display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header h2 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .subtitle { color: var(--text-muted); margin: 0.25rem 0 0 0; }

  .card { padding: 0; overflow: hidden; }
  .card-title { padding: 1rem 2rem; margin: 0; font-size: 1rem; background: rgba(255,255,255,0.02); border-bottom: 1px solid var(--border-color); }

  .table-wrapper { overflow-x: auto; }
  .data-table { width: 100%; border-collapse: collapse; text-align: left; }
  .data-table th { padding: 1rem 2rem; font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); border-bottom: 1px solid var(--border-color); }
  .data-table td { padding: 1rem 2rem; border-bottom: 1px solid var(--border-color); font-size: 0.95rem; }
  .data-table tr:last-child td { border-bottom: none; }

  code { font-family: monospace; background: rgba(255,255,255,0.05); padding: 0.2rem 0.4rem; border-radius: 0.4rem; font-size: 0.85rem; }

  .skeleton-monitor { height: 500px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
