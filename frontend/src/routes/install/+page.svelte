<script lang="ts">
  import InstallWizard from '$lib/components/InstallWizard.svelte';
</script>

<svelte:head>
  <title>Installation · ggit</title>
</svelte:head>

<div class="page-container">
  <InstallWizard />
</div>

<style>
  .page-container {
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background: radial-gradient(circle at top right, rgba(var(--primary-rgb), 0.1), transparent);
  }
</style>
