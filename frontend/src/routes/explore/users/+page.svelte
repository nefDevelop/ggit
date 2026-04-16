<script lang="ts">
  import ExploreUsers from '$lib/components/ExploreUsers.svelte';
</script>

<svelte:head>
  <title>Explore Users · Ggit</title>
</svelte:head>

<div class="page-container">
  <ExploreUsers />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
