<script lang="ts">
  import ReleaseManager from '$lib/components/ReleaseManager.svelte';
</script>

<svelte:head>
  <title>Releases · ggit</title>
</svelte:head>

<div class="page-container">
  <ReleaseManager />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
