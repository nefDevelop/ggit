<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let settings = $state<any>(null);
  let loading = $state(true);
  let form = $state<any>({});
  let submitting = $state(false);
  let error = $state('');

  async function fetchSettings() {
    try {
      const res = await apiFetch('/api/v2/install/settings');
      if (res.ok) {
        const data = await res.json();
        settings = data;
        form = { ...data.default_settings };
      }
    } finally {
      loading = false;
    }
  }

  async function install() {
    submitting = true;
    error = '';
    try {
      const res = await apiFetch('/api/v2/install', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form)
      });
      if (res.ok) {
        window.location.href = '/user/login';
      } else {
        const data = await res.json();
        error = data.error || data.message || 'Installation failed.';
      }
    } catch (e) {
      error = 'A network error occurred.';
    } finally {
      submitting = false;
    }
  }

  onMount(fetchSettings);
</script>

<div class="install-wizard island tonal">
  <header>
    <h1>Ggit Installation</h1>
    <p>Configure your self-hosted Git platform in minutes.</p>
  </header>

  {#if loading}
    <div class="spinner"></div>
  {:else if settings}
    <form onsubmit={(e) => { e.preventDefault(); install(); }}>
      <section class="island card">
        <h3>Database Settings</h3>
        <div class="grid">
          <div class="field">
            <label for="db_type">Database Type</label>
            <select id="db_type" bind:value={form.db_type}>
              {#each settings.db_types as type}
                <option value={type}>{type}</option>
              {/each}
            </select>
          </div>
          <div class="field">
            <label for="db_host">Host</label>
            <input id="db_host" type="text" bind:value={form.db_host} />
          </div>
        </div>
      </section>

      <section class="island card">
        <h3>General Settings</h3>
        <div class="grid">
           <div class="field">
             <label for="app_name">Application Name</label>
             <input id="app_name" type="text" bind:value={form.app_name} />
           </div>
           <div class="field">
             <label for="run_user">Run User</label>
             <input id="run_user" type="text" bind:value={form.run_user} />
           </div>
        </div>
      </section>

      <div class="info-banner island tonal">
        <span class="material-symbols-rounded">admin_panel_settings</span>
        <div class="info-text">
          <strong>Automatic Administration</strong>
          <p>For your convenience, the first user account created after installation will automatically be granted full administrator privileges.</p>
        </div>
      </div>

      {#if error}
        <div class="error-banner">{error}</div>
      {/if}

      <div class="actions">
        <button class="btn btn-primary lg" type="submit" disabled={submitting}>
          {submitting ? 'Installing...' : 'Install Ggit'}
        </button>
      </div>
    </form>
  {/if}
</div>

<style>
  .install-wizard { max-width: 800px; margin: 4rem auto; padding: 3rem; box-shadow: var(--shadow-2xl); display: flex; flex-direction: column; gap: 2rem; }
  h1 { margin: 0; font-size: 2.25rem; font-weight: 900; color: #fff; }
  p { color: var(--text-muted); opacity: 0.8; }
  
  form { display: flex; flex-direction: column; gap: 2rem; }
  .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  @media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }
  
  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); }
  .field input, .field select { background: rgba(0,0,0,0.2); border: 1px solid var(--md-sys-color-outline-variant); color: #fff; padding: 0.75rem 1rem; border-radius: 0.75rem; }

  .card { border: 1px solid var(--md-sys-color-outline-variant); padding: 1.5rem; }
  
  .error-banner { background: rgba(239, 68, 68, 0.1); color: #ef4444; padding: 1rem; border-radius: 0.75rem; text-align: center; font-weight: 700; }

  .info-banner { display: flex; gap: 1rem; align-items: center; padding: 1.5rem; border: 1px solid var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  .info-banner span { font-size: 2.5rem; color: var(--primary-color); }
  .info-text strong { display: block; font-size: 1.1rem; color: #fff; margin-bottom: 0.25rem; }
  .info-text p { margin: 0; font-size: 0.9rem; line-height: 1.4; }

  .lg { width: 100%; padding: 1rem; font-weight: 800; font-size: 1.1rem; }
</style>
