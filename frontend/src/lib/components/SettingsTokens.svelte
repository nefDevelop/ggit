<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Button from './ui/Button.svelte';
  import Input from './ui/Input.svelte';
  import Card from './ui/Card.svelte';
  import Badge from './ui/Badge.svelte';
  
  let tokens = $state<any[]>([]);
  let loading = $state(true);
  let newTokenName = $state('');
  let generatedToken = $state('');
  let submitting = $state(false);

  async function fetchTokens() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/tokens');
      if (res.ok) tokens = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function createToken() {
    if (!newTokenName) return;
    submitting = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/tokens', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newTokenName })
      });
      if (res.ok) {
        const token = await res.json();
        generatedToken = token.sha1; // Ggit returns the clear token once
        tokens = [token, ...tokens];
        newTokenName = '';
      }
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  async function deleteToken(id: number) {
    if (!confirm('Are you sure you want to delete this token?')) return;
    try {
      const res = await apiFetch(`/api/v2/user/settings/tokens/${id}`, {
        method: 'DELETE'
      });
      if (res.ok) {
        tokens = tokens.filter(t => t.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchTokens);
</script>

<Card padding="none" class="settings-card">
  <header class="card-header">
    <h3 class="card-title">Access Tokens</h3>
  </header>
  
  <div class="card-body">
    <p class="desc">Personal access tokens function like ordinary OAuth access tokens. They can be used instead of a password for Git over HTTP, or can be used to authenticate to the API over Basic Authentication.</p>

    {#if generatedToken}
      <div class="generated-token">
        <div class="token-warning">
           <span class="material-symbols-rounded">warning</span>
           <strong>Make sure to copy your personal access token now. You won’t be able to see it again!</strong>
        </div>
        <div class="token-display">
           <code>{generatedToken}</code>
           <Button variant="ghost" size="sm" icon="content_copy" onclick={() => navigator.clipboard.writeText(generatedToken)} />
        </div>
        <Button onclick={() => generatedToken = ''}>I've copied it</Button>
      </div>
    {/if}

    <div class="tokens-list">
      {#each tokens as token}
        <Card padding="sm" class="token-row">
          <div class="token-info">
            <strong>{token.name}</strong>
            <span class="token-meta">Created {new Date(token.created_at).toLocaleDateString()}</span>
          </div>
          <Button variant="ghost" size="sm" icon="delete" class="danger-icon" onclick={() => deleteToken(token.id)} />
        </Card>
      {/each}
      {#if !loading && tokens.length === 0}
        <div class="empty-state">No access tokens generated yet.</div>
      {/if}
    </div>

    <hr class="divider" />

    <form class="add-token-form" onsubmit={(e) => { e.preventDefault(); createToken(); }}>
      <h4>Generate New Token</h4>
      <div class="form-group">
        <div class="inline-form">
          <Input bind:value={newTokenName} placeholder="Project Alpha Release" label="Token Name" />
          <div class="btn-align">
            <Button type="submit" disabled={submitting || !newTokenName} icon={submitting ? 'sync' : 'key'}>
              {submitting ? 'Generating...' : 'Generate Token'}
            </Button>
          </div>
        </div>
      </div>
    </form>
  </div>
</Card>

<style>
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; color: var(--md-sys-color-on-surface); }
  
  .card-body { padding: 2rem; }
  .desc { color: var(--text-muted); margin-bottom: 2rem; }

  .generated-token { padding: 1.5rem; margin-bottom: 2rem; border-left: 4px solid var(--md-sys-color-primary); background: var(--md-sys-color-primary-container); border-radius: 0.5rem; display: flex; flex-direction: column; gap: 1rem; }
  .token-warning { display: flex; align-items: center; gap: 0.75rem; color: var(--md-sys-color-primary); font-size: 0.9rem; }
  .token-display { display: flex; gap: 1rem; align-items: center; background: var(--md-sys-color-background); border: 1px solid var(--md-sys-color-outline-variant); padding: 0.75rem 1rem; border-radius: 0.75rem; }
  .token-display code { flex: 1; font-family: monospace; font-size: 1.1rem; color: var(--md-sys-color-on-surface); overflow-x: auto; }

  .tokens-list { display: flex; flex-direction: column; gap: 1rem; margin-bottom: 2rem; }
  :global(.token-row) { display: flex; align-items: center; gap: 1.5rem; background: var(--md-sys-color-surface-container) !important; }
  .token-info { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; }
  .token-info strong { font-size: 1rem; color: var(--md-sys-color-on-surface); }
  .token-meta { font-size: 0.75rem; color: var(--text-muted); opacity: 0.7; }

  .divider { border: 0; border-top: 1px solid var(--md-sys-color-outline-variant); margin: 2rem 0; }

  .add-token-form { display: flex; flex-direction: column; gap: 1rem; }
  .add-token-form h4 { margin: 0; font-size: 1.1rem; margin-bottom: 0.5rem; color: var(--md-sys-color-on-surface); }
  
  .inline-form { display: flex; gap: 1rem; align-items: flex-end; }
  .btn-align { margin-bottom: 0.25rem; }
  
  .empty-state { text-align: center; padding: 2rem; color: var(--text-muted); border: 2px dashed var(--md-sys-color-outline-variant); border-radius: 1rem; }
  :global(.danger-icon:hover) { color: var(--color-danger) !important; }
</style>

