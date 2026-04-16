<script lang="ts">
  import { onMount } from 'svelte';
  
  let tokens = $state<any[]>([]);
  let loading = $state(true);
  let newTokenName = $state('');
  let generatedToken = $state('');
  let submitting = $state(false);

  async function fetchTokens() {
    loading = true;
    try {
      const res = await fetch('/api/v2/user/settings/tokens');
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
      const res = await fetch('/api/v2/user/settings/tokens', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newTokenName })
      });
      if (res.ok) {
        const token = await res.json();
        generatedToken = token.sha1; // Gogs returns the clear token once
        tokens = [...tokens, token];
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
      const res = await fetch(`/api/v2/user/settings/tokens/${id}`, {
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

<div class="settings-card island">
  <header class="card-header">
    <h3 class="card-title">Access Tokens</h3>
  </header>
  
  <div class="card-body">
    <p class="desc">Personal access tokens function like ordinary OAuth access tokens. They can be used instead of a password for Git over HTTP, or can be used to authenticate to the API over Basic Authentication.</p>

    {#if generatedToken}
      <div class="generated-token island tonal warning">
        <div class="token-warning">
           <span class="material-symbols-rounded">warning</span>
           <strong>Make sure to copy your personal access token now. You won’t be able to see it again!</strong>
        </div>
        <div class="token-display">
           <code>{generatedToken}</code>
           <button class="btn btn-icon" onclick={() => navigator.clipboard.writeText(generatedToken)}>
             <span class="material-symbols-rounded">content_copy</span>
           </button>
        </div>
        <button class="btn btn-primary" onclick={() => generatedToken = ''}>I've copied it</button>
      </div>
    {/if}

    <div class="tokens-list">
      {#each tokens as token}
        <div class="token-row island tonal">
          <div class="token-info">
            <strong>{token.name}</strong>
            <span class="token-meta">Created {new Date(token.created_at).toLocaleDateString()}</span>
          </div>
          <button class="btn btn-icon btn-danger-text" onclick={() => deleteToken(token.id)}>
            <span class="material-symbols-rounded">delete</span>
          </button>
        </div>
      {/each}
      {#if !loading && tokens.length === 0}
        <div class="empty-state">No access tokens generated yet.</div>
      {/if}
    </div>

    <hr class="divider" />

    <form class="add-token-form" onsubmit={(e) => { e.preventDefault(); createToken(); }}>
      <h4>Generate New Token</h4>
      <div class="form-group">
        <label for="token_name">Token Name</label>
        <div class="inline-form">
          <input type="text" id="token_name" bind:value={newTokenName} placeholder="Project Alpha Release" required />
          <button type="submit" class="btn btn-primary" disabled={submitting || !newTokenName}>
            {submitting ? 'Generating...' : 'Generate Token'}
          </button>
        </div>
      </div>
    </form>
  </div>
</div>

<style>
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--border-color); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  
  .card-body { padding: 2rem; }
  .desc { color: var(--text-muted); margin-bottom: 2rem; }

  .generated-token { padding: 1.5rem; margin-bottom: 2rem; border-left: 4px solid var(--primary-color); background: rgba(var(--primary-rgb), 0.05); display: flex; flex-direction: column; gap: 1rem; }
  .token-warning { display: flex; align-items: center; gap: 0.75rem; color: var(--primary-color); font-size: 0.9rem; }
  .token-display { display: flex; gap: 1rem; align-items: center; background: rgba(0,0,0,0.3); padding: 0.75rem 1rem; border-radius: 0.75rem; }
  .token-display code { flex: 1; font-family: monospace; font-size: 1.1rem; color: #fff; overflow-x: auto; }

  .tokens-list { display: flex; flex-direction: column; gap: 1rem; margin-bottom: 2rem; }
  .token-row { display: flex; align-items: center; gap: 1.5rem; padding: 1rem 1.5rem; }
  .token-info { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; }
  .token-info strong { font-size: 1rem; }
  .token-meta { font-size: 0.75rem; color: var(--text-muted); opacity: 0.5; }

  .divider { border: 0; border-top: 1px solid var(--border-color); margin: 2rem 0; }

  .add-token-form { display: flex; flex-direction: column; gap: 1rem; }
  .add-token-form h4 { margin: 0; font-size: 1.1rem; margin-bottom: 0.5rem; }
  
  .inline-form { display: flex; gap: 1rem; }
  .inline-form input { flex: 1; }
  
  input { background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); padding: 0.75rem 1rem; border-radius: 0.75rem; color: var(--text-color); font-family: inherit; }
  input:focus { outline: none; border-color: var(--primary-color); background: rgba(0,0,0,0.4); }

  .empty-state { text-align: center; padding: 2rem; color: var(--text-muted); border: 2px dashed var(--border-color); border-radius: 1rem; }
  .btn-danger-text { color: #ef4444; }
  .btn-danger-text:hover { background: rgba(239, 68, 68, 0.1); }
</style>
