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
    width: 100%;
    margin: 0;
    padding: 0 1.5rem 4rem;
  }
</style>
