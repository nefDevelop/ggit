<script lang="ts">
  import Button from "$lib/components/ui/Button.svelte";
  import GitMerge from "lucide-svelte/icons/git-merge";

  let { title = $bindable(), content = $bindable(), submitting = false, onCreate }: {
    title: string;
    content: string;
    submitting: boolean;
    onCreate: () => void;
  } = $props();
</script>

<section class="creation-details shadow-xl">
  <div class="form-header">
    <h3>New Pull Request</h3>
    <p class="text-muted">
      Fill in the details to submit your changes for review.
    </p>
  </div>

  <div class="form-fields">
    <div class="field">
      <label for="pr-title">Title</label>
      <input
        id="pr-title"
        type="text"
        bind:value={title}
        placeholder="A concise summary of changes"
        required
      />
    </div>
    <div class="field">
      <label for="pr-body">Description</label>
      <textarea
        id="pr-body"
        bind:value={content}
        placeholder="Explain what this PR does and why it's needed..."
      ></textarea>
    </div>
  </div>

  <div class="form-actions">
    <Button
      variant="primary"
      size="lg"
      onclick={onCreate}
      disabled={!title || submitting}
      loading={submitting}
      icon={GitMerge}
      class="w-full-mobile"
    >
      Create Pull Request
    </Button>
  </div>
</section>

<style>
  .creation-details {
    padding: 2.5rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
  .form-header h3 { margin: 0; font-size: 1.5rem; font-weight: 800; }
  .form-header p { margin: 0.25rem 0 0; font-size: 0.95rem; }
  .form-fields { display: flex; flex-direction: column; gap: 1.5rem; }
  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.8rem; font-weight: 700; color: var(--text-muted); }
  .field input {
    background: var(--sunken-2);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: #fff;
    padding: 1rem;
    border-radius: 1rem;
    font-size: 1.1rem;
    width: 100%;
    outline: none;
  }
  .field input:focus { border-color: var(--md-sys-color-primary); background: var(--md-sys-color-surface-container-highest); box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary), transparent 85%); }
  .field textarea {
    background: var(--sunken-2);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: #fff;
    padding: 1rem;
    border-radius: 1rem;
    font-size: 1rem;
    height: 350px;
    resize: none;
    width: 100%;
    outline: none;
    font-family: inherit;
  }
  .field textarea:focus { border-color: var(--md-sys-color-primary); background: var(--md-sys-color-surface-container-highest); box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary), transparent 85%); }
  .form-actions { margin-top: 1rem; }
  @media (max-width: 1200px) {
    .creation-details { padding: 1.5rem; }
  }
</style>
