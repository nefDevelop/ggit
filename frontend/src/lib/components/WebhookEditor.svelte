<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { getContext } from 'svelte';

  let { hook = null, onclose } = $props();
  const repo: any = getContext('repo');

  let url = $state(hook?.url || '');
  let contentType = $state(hook?.content_type || 'json');
  let secret = $state(hook?.secret || '');
  let active = $state(hook ? hook.active : true);

  let events = $state((hook?.events || ['push']) as string[]);

  let submitting = $state(false);  let error = $state('');

  const availableEvents = [
    { id: 'push', label: 'Push Events', desc: 'Any git push to the repository.' },
    { id: 'create', label: 'Branch/Tag Created', desc: 'When a new branch or tag is created.' },
    { id: 'delete', label: 'Branch/Tag Deleted', desc: 'When a branch or tag is deleted.' },
    { id: 'issues', label: 'Issues', desc: 'Opened, closed, reopened, edited, assigned, etc.' },
    { id: 'issue_comment', label: 'Issue Comments', desc: 'Created, edited, or deleted.' },
    { id: 'pull_request', label: 'Pull Requests', desc: 'Opened, closed, reopened, edited, merged, etc.' },
    { id: 'release', label: 'Releases', desc: 'Published, edited, or deleted.' },
  ];

  function toggleEvent(id: string) {
    if (events.includes(id)) {
      events = events.filter(e => e !== id);
    } else {
      events = [...events, id];
    }
  }

  async function handleSave() {
    submitting = true;
    error = '';

    const payload = {
      url,
      content_type: contentType,
      secret,
      active,
      events
    };

    const method = hook ? 'PATCH' : 'POST';
    const apiUrl = hook 
      ? `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/${hook.id}`
      : `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks`;

    try {
      const res = await apiFetch(apiUrl, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (res.ok) {
        onclose(true);
      } else {
        const data = await res.json();
        error = data.error || 'Failed to save webhook.';
      }
    } catch (e) {
      error = 'A network error occurred.';
    } finally {
      submitting = false;
    }
  }

  async function handleTest() {
    if (!hook) return;
    submitting = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/hooks/${hook.id}/test`, {
        method: 'POST'
      });
      if (res.ok) {
        alert('Test delivery triggered successfully!');
      }
    } finally {
      submitting = false;
    }
  }
</script>

<div class="webhook-editor">
  <header class="modal-header">
    <div class="title-zone">
      <span class="material-symbols-rounded">webhook</span>
      <h3>{hook ? 'Edit Webhook' : 'Add New Webhook'}</h3>
    </div>
    <button class="btn-close" onclick={() => onclose()}>
      <span class="material-symbols-rounded">close</span>
    </button>
  </header>

  <form onsubmit={(e) => { e.preventDefault(); handleSave(); }}>
    <div class="form-body">
      <div class="field">
        <label for="url">Payload URL</label>
        <input id="url" type="url" bind:value={url} placeholder="https://example.com/postreceive" required />
        <p class="helper">The URL that will receive the POST requests.</p>
      </div>

      <div class="form-row">
        <div class="field">
          <label for="content-type">Content Type</label>
          <select id="content-type" bind:value={contentType}>
            <option value="json">application/json</option>
            <option value="form">application/x-www-form-urlencoded</option>
          </select>
        </div>
      </div>

      <div class="field">
        <label for="secret">Secret</label>
        <input id="secret" type="password" bind:value={secret} placeholder="Optional secret for HMAC verification" />
        <p class="helper">Used to generate a signature in <code>X-Ggit-Signature-256</code> header.</p>
      </div>

      <div class="events-picker">
        <label>Which events would you like to trigger this webhook?</label>
        <div class="events-grid">
          {#each availableEvents as ev}
            <label class="event-card" class:active={events.includes(ev.id)}>
              <input type="checkbox" checked={events.includes(ev.id)} onchange={() => toggleEvent(ev.id)} />
              <div class="event-content">
                <strong>{ev.label}</strong>
                <span>{ev.desc}</span>
              </div>
            </label>
          {/each}
        </div>
      </div>

      <div class="toggle-field">
         <label class="switch">
           <input type="checkbox" bind:checked={active} />
           <span class="slider"></span>
         </label>
         <div class="toggle-text">
            <strong>Active</strong>
            <span>Deliveries will be attempted for this webhook.</span>
         </div>
      </div>

      {#if error}
        <div class="error-msg">{error}</div>
      {/if}
    </div>

    <footer class="modal-footer">
      {#if hook}
        <button type="button" class="btn btn-tonal" onclick={handleTest} disabled={submitting}>
          <span class="material-symbols-rounded">send</span>
          Test Webhook
        </button>
      {/if}
      <div class="spacer"></div>
      <button type="button" class="btn btn-outline" onclick={() => onclose()}>Cancel</button>
      <button type="submit" class="btn btn-primary lg" disabled={submitting}>
        {submitting ? 'Saving...' : hook ? 'Update Webhook' : 'Add Webhook'}
      </button>
    </footer>
  </form>
</div>

<style>
  .webhook-editor { display: flex; flex-direction: column; }
  
  .modal-header { padding: 1.5rem 2rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .title-zone { display: flex; align-items: center; gap: 0.75rem; }
  .title-zone span { color: var(--primary-color); }
  .title-zone h3 { margin: 0; font-size: 1.25rem; font-weight: 700; }
  
  .btn-close { background: transparent; border: none; color: var(--text-muted); cursor: pointer; display: flex; padding: 0.25rem; border-radius: 50%; }
  .btn-close:hover { background: rgba(255,255,255,0.05); color: var(--text-color); }

  .form-body { padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
  
  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.85rem; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }
  .field input, .field select { background: rgba(0,0,0,0.2); border: 1px solid var(--md-sys-color-outline-variant); color: #fff; padding: 0.75rem 1rem; border-radius: 0.75rem; font-size: 1rem; width: 100%; transition: border-color 0.2s; }
  .field input:focus { border-color: var(--primary-color); outline: none; }
  
  .helper { margin: 0; font-size: 0.8rem; color: var(--text-muted); opacity: 0.7; }

  .events-picker { display: flex; flex-direction: column; gap: 1rem; }
  .events-picker label { font-size: 0.9rem; font-weight: 600; }
  
  .events-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0.75rem; }
  @media (max-width: 600px) { .events-grid { grid-template-columns: 1fr; } }
  .event-card { display: flex; gap: 1rem; padding: 1rem; background: rgba(255,255,255,0.02); border: 1px solid var(--md-sys-color-outline-variant); border-radius: 1rem; cursor: pointer; transition: all 0.2s; }
  .event-card:hover { background: rgba(255,255,255,0.05); }
  .event-card.active { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  .event-card input { margin-top: 0.2rem; }
  
  .event-content { display: flex; flex-direction: column; gap: 0.1rem; }
  .event-content strong { font-size: 0.9rem; }
  .event-content span { font-size: 0.75rem; color: var(--text-muted); }

  .toggle-field { display: flex; align-items: flex-start; gap: 1rem; background: rgba(255,255,255,0.03); padding: 1.25rem; border-radius: 1rem; }
  .toggle-text { display: flex; flex-direction: column; gap: 0.1rem; }
  .toggle-text strong { font-size: 0.95rem; }
  .toggle-text span { font-size: 0.85rem; color: var(--text-muted); }

  .error-msg { color: #f87171; font-weight: 600; font-size: 0.9rem; }

  .modal-footer { padding: 1.5rem 2rem; display: flex; align-items: center; gap: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); background: rgba(0,0,0,0.1); }
  .spacer { flex: 1; }
  
  .switch { position: relative; display: inline-block; width: 44px; height: 24px; flex-shrink: 0; }
  .switch input { opacity: 0; width: 0; height: 0; }
  .slider { position: absolute; cursor: pointer; inset: 0; background-color: rgba(255,255,255,0.1); transition: .4s; border-radius: 24px; }
  .slider:before { position: absolute; content: ""; height: 18px; width: 18px; left: 3px; bottom: 3px; background-color: white; transition: .4s; border-radius: 50%; }
  input:checked + .slider { background-color: var(--primary-color); }
  input:checked + .slider:before { transform: translateX(20px); }

  @media (max-width: 600px) {
    .events-grid { grid-template-columns: 1fr; }
    .modal-footer { flex-direction: column; align-items: stretch; }
  }
</style>
