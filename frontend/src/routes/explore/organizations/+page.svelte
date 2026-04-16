<script lang="ts">
  import ExploreOrganizations from '$lib/components/ExploreOrganizations.svelte';
</script>

<svelte:head>
  <title>Explore Organizations · Ggit</title>
</svelte:head>

<div class="page-container">
  <ExploreOrganizations />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
