<script lang="ts">
  import UserSettings from '$lib/components/UserSettings.svelte';
</script>

<svelte:head>
  <title>User Settings · Ggit</title>
</svelte:head>

<div class="page-container">
  <UserSettings />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }

  @media (max-width: 600px) {
    .page-container {
      padding: 1rem 0.5rem;
    }
  }
</style>
