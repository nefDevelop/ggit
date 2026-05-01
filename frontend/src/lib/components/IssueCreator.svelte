<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { goto } from "$app/navigation";
  import PageHeader from "./ui/PageHeader.svelte";
  import Button from "./ui/Button.svelte";
  import Card from "./ui/Card.svelte";

  const repo: any = getContext("repo");

  let title = $state("");
  let body = $state("");
  let submitting = $state(false);
  let error = $state("");

  async function createIssue() {
    if (!title.trim()) {
      error = "Title is required.";
      return;
    }
    submitting = true;
    error = "";
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/issues`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, body }),
      });
      if (res.ok) {
        const data = await res.json();
        // Redirect to newly created issue
        goto(`/${repo.data.owner.username}/${repo.data.name}/issues/${data.number || data.index}`);
      } else {
        error = "Failed to create issue. Please check the inputs.";
      }
    } catch (e) {
      error = "An error occurred. Please try again.";
    } finally {
      submitting = false;
    }
  }
</script>

<div class="issue-creator-view">
  <PageHeader title="New Issue" icon="bug_report" compact>
    {#snippet actions()}
      <Button variant="ghost" onclick={() => goto(`/${repo.data.owner.username}/${repo.data.name}/issues`)}>
        Cancel
      </Button>
    {/snippet}
  </PageHeader>

  <Card class="issue-form-card" padding="lg">
    <div class="form-stack">
      {#if error}
        <div class="form-error-msg">{error}</div>
      {/if}

      <div class="field">
        <label for="title">Title</label>
        <input 
          id="title" 
          type="text" 
          bind:value={title} 
          placeholder="Give your issue a clear title..." 
          disabled={submitting} 
          required 
        />
      </div>

      <div class="field">
        <label for="body">Description</label>
        <textarea 
          id="body" 
          bind:value={body} 
          placeholder="Describe the problem, steps to reproduce, etc..." 
          disabled={submitting}
        ></textarea>
      </div>

      <div class="form-actions">
        <Button variant="primary" onclick={createIssue} disabled={submitting || !title.trim()}>
          {submitting ? "Creating..." : "Submit Issue"}
        </Button>
      </div>
    </div>
  </Card>
</div>

<style>
  .issue-creator-view {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    animation: fadeIn 0.3s ease-out;
  }

  .form-stack {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .form-error-msg {
    padding: 1rem 1.25rem;
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
    border-radius: 0.75rem;
    font-size: 0.9rem;
    font-weight: 700;
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .field label {
    font-size: 0.8rem;
    font-weight: 800;
    color: var(--text-color);
  }

  .field input {
    background: var(--surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 0.85rem 1.25rem;
    border-radius: 1rem;
    color: var(--text-color);
    font-size: 1.1rem;
    font-family: inherit;
    outline: none;
    transition: all 0.2s;
  }
  .field input:focus {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }

  .field textarea {
    background: var(--surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 1.25rem;
    border-radius: 1rem;
    color: var(--text-color);
    font-size: 1rem;
    font-family: inherit;
    min-height: 250px;
    resize: vertical;
    outline: none;
    transition: all 0.2s;
  }
  .field textarea:focus {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    border-top: 1px solid var(--md-sys-color-outline-variant);
    padding-top: 1.5rem;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>
