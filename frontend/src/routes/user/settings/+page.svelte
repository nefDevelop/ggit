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
</style>
