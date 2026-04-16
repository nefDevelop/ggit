<script lang="ts">
  import AdminOrgs from '$lib/components/AdminOrgs.svelte';
</script>

<svelte:head>
  <title>Organization Oversight · Admin · ggit</title>
</svelte:head>

<div class="page-container">
  <AdminOrgs />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
