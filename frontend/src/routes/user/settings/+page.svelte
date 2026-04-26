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
    max-width: 1600px;
    margin: 0 auto;
    padding: 2rem 4rem;
  }

  @media (max-width: 1024px) {
    .page-container {
      padding: 1rem;
    }
  }
</style>
