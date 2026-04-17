<script lang="ts">
  import { onMount } from 'svelte';
  
  let config = $state<any>(null);
  let loading = $state(true);

  async function fetchConfig() {
    loading = true;
    try {
      const res = await fetch('/api/v2/admin/config');
      if (res.ok) config = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  const sections = [
    { key: 'app', label: 'Application' },
    { key: 'server', label: 'Server' },
    { key: 'database', label: 'Database' },
    { key: 'security', label: 'Security' },
    { key: 'repository', label: 'Repository' },
    { key: 'git', label: 'Git' },
    { key: 'email', label: 'Email' }
  ];

  let testEmail = $state('');
  let sendingTest = $state(false);
  let testMsg = $state({ text: '', type: '' });

  async function sendTestEmail() {
    if (!testEmail) return;
    sendingTest = true;
    testMsg = { text: '', type: '' };
    try {
      const res = await fetch('/api/v2/admin/config/test_email', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: testEmail })
      });
      if (res.ok) {
        testMsg = { text: 'Test email successfully sent!', type: 'success' };
      } else {
        const data = await res.json();
        testMsg = { text: data.error || 'Failed to send test email.', type: 'error' };
      }
    } catch (e) {
      testMsg = { text: 'A network error occurred.', type: 'error' };
    } finally {
      sendingTest = false;
    }
  }

  onMount(fetchConfig);
</script>

<div class="admin-config">
  <header class="section-header">
    <h2>Server Configuration</h2>
    <p class="subtitle">Read-only view of the <code>app.ini</code> configuration file.</p>
  </header>

  {#if loading}
    <div class="skeleton-config"></div>
  {:else if config}
    <div class="config-sections">
      {#each sections as section}
        <div class="config-card island">
          <h3 class="card-title">{section.label}</h3>
          <div class="config-grid">
            {#each Object.entries(config[section.key] || {}) as [label, value]}
              <div class="config-item">
                <span class="key">{label}</span>
                <span class="val">{typeof value === 'object' ? JSON.stringify(value) : value}</span>
              </div>
            {/each}
          </div>
        </div>
      {/each}

      <div class="config-card island tonal test-email-section">
        <h3 class="card-title">Test Email Configuration</h3>
        <div class="test-email-content">
          <p>Send a test email to ensure your SMTP settings are correct.</p>
          <div class="input-group">
            <input type="email" placeholder="target@example.com" bind:value={testEmail} />
            <button class="btn btn-primary" onclick={sendTestEmail} disabled={sendingTest || !testEmail}>
              {sendingTest ? 'Sending...' : 'Send Test Mail'}
            </button>
          </div>
          {#if testMsg.text}
            <p class="status-msg" class:error={testMsg.type === 'error'}>{testMsg.text}</p>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .admin-config { display: flex; flex-direction: column; gap: 2rem; }
  
  .section-header h2 { margin: 0; font-size: 1.75rem; font-weight: 800; }
  .subtitle { color: var(--text-muted); margin: 0.25rem 0 0 0; }

  .config-sections { display: flex; flex-direction: column; gap: 2rem; }
  .config-card { padding: 0; overflow: hidden; }
  .card-title { padding: 1rem 2rem; margin: 0; font-size: 1rem; background: rgba(255,255,255,0.02); border-bottom: 1px solid var(--border-color); }

  .config-grid { padding: 1.5rem 2rem; display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  .config-item { display: flex; flex-direction: column; gap: 0.2rem; }
  .key { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.6; }
  .val { font-size: 0.95rem; font-family: monospace; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

  .skeleton-config { height: 600px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }

  .test-email-content { padding: 2rem; display: flex; flex-direction: column; gap: 1rem; }
  .test-email-content p { margin: 0; color: var(--text-muted); font-size: 0.9rem; }
  .input-group { display: flex; gap: 1rem; max-width: 500px; }
  .input-group input { flex: 1; background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); color: #fff; padding: 0.6rem 1rem; border-radius: 0.5rem; }
  .status-msg { font-size: 0.85rem; font-weight: 700; color: #10b981; }
  .status-msg.error { color: #ef4444; }
</style>
