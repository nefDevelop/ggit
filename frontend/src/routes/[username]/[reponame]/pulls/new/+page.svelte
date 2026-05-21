<script lang="ts">
  import PRCreator from '$lib/components/PRCreator.svelte';
</script>

<svelte:head>
  <title>New Pull Request · ggit</title>
</svelte:head>

<div class="page-container">
  <PRCreator />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
  }
</style>
