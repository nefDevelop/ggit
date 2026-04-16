<script lang="ts">
  import AdminView from '$lib/components/AdminView.svelte';
</script>

<svelte:head>
  <title>Administration · Ggit</title>
</svelte:head>

<div class="page-container">
  <AdminView />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
