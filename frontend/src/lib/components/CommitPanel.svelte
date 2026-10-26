<script lang="ts">
  import { fade } from 'svelte/transition';
  let { onCommit, onCancel, isSaving = false } = $props();

  let summary = $state("");
  let message = $state("");
  let commitChoice = $state("direct");
  let newBranchName = $state("");
</script>

<div class="commit-panel glass-panel">
  <div class="commit-header">
    <div class="header-icon-ring">
      <span class="material-symbols-rounded">commit</span>
    </div>
    <div class="header-text">
        <h3>Commit changes</h3>
        <p>Save your progress to the repository history.</p>
    </div>
  </div>

  <div class="commit-body">
    <div class="form-group">
      <label for="commit-summary">Summary</label>
      <input
        id="commit-summary"
        type="text"
        bind:value={summary}
        placeholder="Update file description..."
        aria-label="Commit summary"
      />
    </div>

    <div class="form-group">
      <label for="commit-message">Extended Description</label>
      <textarea
        id="commit-message"
        bind:value={message}
        placeholder="Add an optional extended description..."
        aria-label="Extended description"
      ></textarea>
    </div>

    <div class="commit-choice">
      <label class="choice-item">
        <input
          type="radio"
          name="commit-choice"
          value="direct"
          bind:group={commitChoice}
        />
        <div class="choice-content">
          <div class="choice-meta">
            <strong>Commit directly to the base branch</strong>
            <p>Your changes will be immediately visible.</p>
          </div>
          <span class="material-symbols-rounded choice-icon">merge_type</span>
        </div>
      </label>

      <label class="choice-item">
        <input
          type="radio"
          name="commit-choice"
          value="branch"
          bind:group={commitChoice}
        />
        <div class="choice-content">
          <div class="choice-meta">
            <strong>Create a new branch for this commit</strong>
            <p>Safely iterate without affecting the main line.</p>
          </div>
          <span class="material-symbols-rounded choice-icon">branch</span>
        </div>
      </label>
      
      {#if commitChoice === "branch"}
        <div class="branch-input-container" transition:fade>
          <span class="material-symbols-rounded branch-prefix">subdirectory_arrow_right</span>
          <input
            type="text"
            bind:value={newBranchName}
            placeholder="new-branch-name"
            class="branch-input"
          />
        </div>
      {/if}
    </div>
  </div>

  <footer class="commit-footer">
    <button class="btn btn-ghost" onclick={onCancel}>Cancel</button>
    <button
      class="btn btn-action-primary"
      disabled={isSaving || (commitChoice === "branch" && !newBranchName)}
      onclick={() =>
        onCommit({ summary, message, commitChoice, newBranchName })}
    >
      <span class="material-symbols-rounded">{isSaving ? 'sync' : 'check'}</span>
      <span>{isSaving ? "Committing..." : "Commit changes"}</span>
    </button>
  </footer>
</div>

<style>
  .commit-panel {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    width: 100%;
    max-width: 540px;
    padding: 2.5rem;
    background: var(--surface);
    border-radius: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    box-shadow: var(--shadow-xl);
    position: relative;
  }

  .commit-header {
    display: flex;
    align-items: flex-start;
    gap: 1.25rem;
  }
  
  .header-icon-ring {
    width: 48px;
    height: 48px;
    border-radius: 1rem;
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }
  
  .header-text h3 { margin: 0; font-size: 1.5rem; font-weight: 800; color: var(--text-color); }
  .header-text p { margin: 0.25rem 0 0; font-size: 0.85rem; color: var(--text-muted); }

  .commit-body {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); padding-left: 0.5rem; }

  input[type="text"],
  textarea {
    background: var(--md3-surface-2);
    border: 1px solid var(--md-sys-color-outline-variant);
    border-radius: 1rem;
    padding: 0.85rem 1.15rem;
    color: var(--text-color);
    font-family: inherit;
    font-size: 0.95rem;
    transition: all 0.2s;
  }
  
  textarea { min-height: 120px; resize: vertical; }
  
  input:focus, textarea:focus {
    border-color: var(--primary-color);
    background: var(--surface);
    box-shadow: 0 0 0 4px rgba(var(--primary-rgb), 0.1);
    outline: none;
  }

  .commit-choice {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-top: 0.5rem;
  }
  .choice-item {
    display: flex;
    gap: 1rem;
    cursor: pointer;
    padding: 1rem 1.25rem;
    background: var(--md3-surface-2);
    border-radius: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.2s;
  }
  
  .choice-item:hover { background: rgba(var(--primary-rgb), 0.03); border-color: rgba(var(--primary-rgb), 0.2); }
  
  .choice-item:has(input:checked) {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.08);
    box-shadow: 0 4px 12px rgba(var(--primary-rgb), 0.1);
  }
  
  input[type="radio"] { width: 1.25rem; height: 1.25rem; accent-color: var(--primary-color); margin-top: 0.25rem; }

  .choice-content {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
  }
  .choice-meta strong { display: block; font-size: 0.95rem; color: var(--text-color); }
  .choice-meta p { margin: 0; font-size: 0.8rem; color: var(--text-muted); }
  .choice-icon { color: var(--primary-color); opacity: 0.5; font-size: 1.5rem; }

  .branch-input-container {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding-left: 1.5rem;
    animation: slideIn 0.3s ease-out;
  }
  
  @keyframes slideIn { from { opacity: 0; transform: translateX(-10px); } to { opacity: 1; transform: translateX(0); } }
  
  .branch-prefix { font-size: 1.25rem; color: var(--primary-color); opacity: 0.5; }
  .branch-input { flex: 1; font-family: monospace; border-style: dashed; }

  .commit-footer {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--md-sys-color-outline-variant);
  }
  
  .btn-ghost { background: transparent; border: none; color: var(--text-muted); font-weight: 600; padding: 0.75rem 1.5rem; border-radius: 1rem; cursor: pointer; }
  .btn-ghost:hover { background: var(--sunken-1); color: var(--text-color); }

  @media (max-width: 600px) {
    .commit-panel { padding: 1.5rem; border-radius: 1.5rem; }
    .header-text h3 { font-size: 1.25rem; }
    .commit-footer { flex-direction: column-reverse; }
    .commit-footer button { width: 100%; }
  }
</style>
