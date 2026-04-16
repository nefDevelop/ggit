<script lang="ts">
  import AdminRepos from '$lib/components/AdminRepos.svelte';
</script>

<svelte:head>
  <title>Repository Oversight · Admin · ggit</title>
</svelte:head>

<div class="page-container">
  <AdminRepos />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
