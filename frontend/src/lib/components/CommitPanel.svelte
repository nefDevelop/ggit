<script lang="ts">
  let { onCommit, onCancel, isSaving = false } = $props();

  let summary = $state('');
  let message = $state('');
  let commitChoice = $state('direct');
  let newBranchName = $state('');
</script>

<div class="commit-panel island tonal">
  <div class="commit-header">
    <span class="material-symbols-rounded">commit</span>
    <h3>Commit changes</h3>
  </div>

  <div class="commit-body">
    <div class="form-group">
      <input 
        type="text" 
        bind:value={summary} 
        placeholder="Update file description..." 
        aria-label="Commit summary"
      />
    </div>
    
    <div class="form-group">
      <textarea 
        bind:value={message} 
        placeholder="Add an optional extended description..."
        aria-label="Extended description"
      ></textarea>
    </div>

    <div class="commit-choice">
       <label class="choice-item">
         <input type="radio" name="commit-choice" value="direct" bind:group={commitChoice} />
         <div class="choice-text">
            <strong>Commit directly to the base branch</strong>
            <span class="material-symbols-rounded">merge_type</span>
         </div>
       </label>
       
       <label class="choice-item">
         <input type="radio" name="commit-choice" value="branch" bind:group={commitChoice} />
         <div class="choice-text">
            <strong>Create a new branch for this commit</strong>
            <span class="material-symbols-rounded">branch</span>
         </div>
         {#if commitChoice === 'branch'}
           <input type="text" bind:value={newBranchName} placeholder="new-branch-name" class="inline-input" />
         {/if}
       </label>
    </div>
  </div>

  <footer class="commit-footer">
    <button class="btn btn-outline" onclick={onCancel}>Cancel</button>
    <button 
      class="btn btn-primary" 
      disabled={isSaving || (commitChoice === 'branch' && !newBranchName)} 
      onclick={() => onCommit({ summary, message, commitChoice, newBranchName })}
    >
      {isSaving ? 'Committing...' : 'Commit changes'}
    </button>
  </footer>
</div>

<style>
  .commit-panel { display: flex; flex-direction: column; gap: 1.5rem; max-width: 600px; padding: 2rem; border: 1px solid var(--border-color); box-shadow: 0 10px 40px rgba(0,0,0,0.4); }
  
  .commit-header { display: flex; align-items: center; gap: 0.75rem; color: var(--primary-color); }
  .commit-header h3 { margin: 0; font-size: 1.25rem; }
  .commit-header span { font-size: 2rem; }

  .commit-body { display: flex; flex-direction: column; gap: 1.25rem; }
  
  input[type="text"], textarea { background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); border-radius: 0.75rem; padding: 0.75rem 1rem; color: #fff; font-family: inherit; font-size: 0.95rem; }
  textarea { min-height: 100px; resize: vertical; }
  input:focus, textarea:focus { border-color: var(--primary-color); outline: none; background: rgba(0,0,0,0.4); }

  .commit-choice { display: flex; flex-direction: column; gap: 0.75rem; margin-top: 0.5rem; }
  .choice-item { display: flex; gap: 1rem; cursor: pointer; padding: 1rem; background: rgba(255,255,255,0.02); border-radius: 0.75rem; border: 1px solid transparent; transition: all 0.2s; }
  .choice-item:hover { background: rgba(255,255,255,0.04); }
  .choice-item:has(input:checked) { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .choice-text { flex: 1; display: flex; align-items: center; justify-content: space-between; }
  .choice-text strong { font-size: 0.95rem; }
  .choice-text span { opacity: 0.5; font-size: 1.25rem; }

  .inline-input { margin-top: 0.5rem; width: 100%; border-style: dashed; }

  .commit-footer { display: flex; justify-content: flex-end; gap: 1rem; padding-top: 1rem; border-top: 1px solid var(--border-color); }
</style>
