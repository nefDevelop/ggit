<script lang="ts">
  import ExploreView from '$lib/components/ExploreView.svelte';
</script>

<svelte:head>
  <title>Explore · Ggit</title>
</svelte:head>

<div class="page-container">
  <ExploreView />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
