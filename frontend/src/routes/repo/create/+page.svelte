<script lang="ts">
  import { authState } from '$lib/auth.svelte';
  import { goto } from '$app/navigation';

  let name = $state('');
  let description = $state('');
  let isPrivate = $state(false);
  let autoInit = $state(true);
  let errorMsg = $state('');
  let loading = $state(false);

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleSubmit() {
    if (!name) return;
    loading = true;
    errorMsg = '';
    
    try {
      const res = await fetch('/api/v2/user/repos', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'X-Csrf-Token': getCsrfToken()
        },
        body: JSON.stringify({
          name,
          description,
          private: isPrivate,
          auto_init: autoInit
        })
      });
      
      const data = await res.json();
      if (res.ok) {
        goto(`/${authState.user.username}/${name}`);
      } else {
        errorMsg = data.message || 'Failed to create repository';
      }
    } catch (e) {
      errorMsg = 'An error occurred. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<div class="create-container">
  <div class="island create-island">
    <header class="create-header">
      <div class="icon-circle">
        <span class="material-symbols-rounded">add_box</span>
      </div>
      <div class="title-stack">
        <h1>Create New Repository</h1>
        <p class="text-muted">A repository contains all your project's files and revision history.</p>
      </div>
    </header>

    <form class="create-form" onsubmit={handleSubmit}>
      {#if errorMsg}
        <div class="error-box">{errorMsg}</div>
      {/if}

      <div class="form-group">
        <label for="repo-name">Repository Name <span class="required">*</span></label>
        <input 
          type="text" 
          id="repo-name" 
          bind:value={name} 
          placeholder="e.g. awesome-project"
          maxlength="100"
          required
        />
        <p class="field-hint">Great repository names are short and memorable.</p>
      </div>

      <div class="form-group">
        <label for="repo-desc">Description <span class="optional">(optional)</span></label>
        <textarea 
          id="repo-desc" 
          bind:value={description} 
          placeholder="What is this project about?"
        ></textarea>
      </div>

      <div class="options-grid">
        <label class="option-card" class:selected={!isPrivate}>
          <input type="radio" name="visibility" bind:group={isPrivate} value={false} />
          <div class="option-icon">
            <span class="material-symbols-rounded">public</span>
          </div>
          <div class="option-text">
            <strong>Public</strong>
            <p>Anyone on the internet can see this repository.</p>
          </div>
          {#if !isPrivate}
            <span class="material-symbols-rounded check-icon">check_circle</span>
          {/if}
        </label>

        <label class="option-card" class:selected={isPrivate}>
          <input type="radio" name="visibility" bind:group={isPrivate} value={true} />
          <div class="option-icon">
            <span class="material-symbols-rounded">lock</span>
          </div>
          <div class="option-text">
            <strong>Private</strong>
            <p>Only you and people you choose can see this repo.</p>
          </div>
          {#if isPrivate}
            <span class="material-symbols-rounded check-icon">check_circle</span>
          {/if}
        </label>
      </div>

      <div class="checkbox-group island-mini">
        <label class="checkbox-label">
          <input type="checkbox" bind:checked={autoInit} />
          <div class="checkbox-text">
            <strong>Initialize this repository with a README</strong>
            <p>This allows you to immediately clone the repository to your computer.</p>
          </div>
        </label>
      </div>

      <div class="form-footer">
        <a href="/dashboard" class="btn btn-tonal">Cancel</a>
        <button type="submit" class="btn primary-btn" disabled={loading || !name}>
          {#if loading}
            <span class="material-symbols-rounded spin">sync</span>
            Creating...
          {:else}
            Create Repository
          {/if}
        </button>
      </div>
    </form>
  </div>
</div>

<style>
  .create-container {
    max-width: 800px;
    margin: 0 auto;
    padding-top: 2rem;
  }

  .create-island {
    padding: 3rem;
  }

  .create-header {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    margin-bottom: 3rem;
  }

  .icon-circle {
    width: 64px;
    height: 64px;
    background: linear-gradient(135deg, var(--primary-color), #7aa9f7);
    color: var(--on-primary);
    border-radius: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 8px 16px rgba(168, 199, 250, 0.2);
  }

  .icon-circle span {
    font-size: 32px;
  }

  .title-stack h1 {
    font-size: 2rem;
    font-weight: 800;
    margin: 0;
    letter-spacing: -0.04em;
  }

  .title-stack p {
    margin: 0.5rem 0 0 0;
    font-size: 1.1rem;
    opacity: 0.8;
  }

  .create-form {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  label {
    font-weight: 600;
    font-size: 1rem;
  }

  .required { color: #f87171; }
  .optional { font-weight: 400; opacity: 0.6; }

  input[type="text"], textarea {
    background-color: var(--md3-surface-2);
    border: 1px solid var(--border-color);
    border-radius: 1rem;
    padding: 1rem;
    color: var(--text-color);
    font-size: 1rem;
    transition: all 0.2s;
  }

  input:focus, textarea:focus {
    outline: none;
    border-color: var(--primary-color);
    background-color: var(--md3-surface-3);
    box-shadow: 0 0 0 4px rgba(168, 199, 250, 0.1);
  }

  textarea {
    min-height: 100px;
    resize: vertical;
  }

  .field-hint {
    font-size: 0.85rem;
    color: var(--text-muted);
    margin: 0;
  }

  .options-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
  }

  .option-card {
    position: relative;
    border: 1px solid var(--border-color);
    border-radius: 1.5rem;
    padding: 1.5rem;
    cursor: pointer;
    display: flex;
    gap: 1rem;
    transition: all 0.2s;
    background: var(--md3-surface-1);
  }

  .option-card:hover {
    background: var(--md3-surface-2);
  }

  .option-card.selected {
    border-color: var(--primary-color);
    background: rgba(168, 199, 250, 0.05);
  }

  .option-card input {
    display: none;
  }

  .option-icon {
    color: var(--primary-color);
  }

  .option-text strong {
    display: block;
    font-size: 1.1rem;
    margin-bottom: 0.25rem;
  }

  .option-text p {
    margin: 0;
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .check-icon {
    position: absolute;
    top: 1rem;
    right: 1rem;
    color: var(--primary-color);
    font-size: 20px;
  }

  .island-mini {
    background: var(--md3-surface-1);
    border: 1px solid var(--border-color);
    border-radius: 1rem;
    padding: 1.25rem;
  }

  .checkbox-label {
    display: flex;
    gap: 1rem;
    cursor: pointer;
    align-items: flex-start;
  }

  .checkbox-label input {
    margin-top: 0.25rem;
    transform: scale(1.2);
  }

  .checkbox-text strong {
    display: block;
    margin-bottom: 0.25rem;
  }

  .checkbox-text p {
    margin: 0;
    font-size: 0.85rem;
    color: var(--text-muted);
  }

  .form-footer {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 1rem;
    padding-top: 2rem;
    border-top: 1px solid var(--border-color);
  }

  .primary-btn {
    min-width: 200px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
  }

  .error-box {
    background: rgba(239, 68, 68, 0.1);
    color: #f87171;
    padding: 1rem;
    border-radius: 1rem;
    border: 1px solid rgba(239, 68, 68, 0.2);
    font-weight: 600;
  }

  .spin {
    animation: rotate 2s linear infinite;
  }

  @keyframes rotate {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }

  @media (max-width: 600px) {
    .options-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
