<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { getContext } from 'svelte';
  import Input from "$lib/components/ui/Input.svelte";
  import Select from "$lib/components/ui/Select.svelte";
  import { toast } from "$lib/toast.svelte";

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

  const contentTypeOptions = [
    { value: 'json', label: 'application/json' },
    { value: 'form', label: 'application/x-www-form-urlencoded' }
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
        toast.success('Test delivery triggered successfully!');
      }
    } finally {
      submitting = false;
    }
  }
</script>

<div class="webhook-editor">
  <form onsubmit={(e) => { e.preventDefault(); handleSave(); }}>
    <div class="form-body">
      <Input 
        label="Payload URL" 
        bind:value={url} 
        placeholder="https://example.com/postreceive"
      />

      />
      <p class="helper">The URL that will receive the POST requests.</p>

      <Select
        label="Content Type"
        options={contentTypeOptions}
        bind:value={contentType}
      />

      <Input
        label="Secret"
        type="password"
        bind:value={secret}
        placeholder="Optional secret for HMAC verification"
      />
      <p class="helper">Used to generate a signature in <code>X-Ggit-Signature-256</code> header.</p>

      <div class="events-picker">
        <label class="section-label">Which events would you like to trigger this webhook?</label>
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
  
  .form-body { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .helper { margin: -0.5rem 0 0 0.25rem; font-size: 0.8rem; color: var(--text-muted); opacity: 0.7; }

  .section-label { font-size: 0.85rem; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }

  .events-picker { display: flex; flex-direction: column; gap: 1rem; }
  
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

  .modal-footer { margin-top: 2rem; padding-top: 1.5rem; display: flex; align-items: center; gap: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); }
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
