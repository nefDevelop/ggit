<script lang="ts">
  import AdminDashboard from '$lib/components/AdminDashboard.svelte';
</script>

<svelte:head>
  <title>Admin Dashboard · Ggit</title>
</svelte:head>

<div class="page-container">
  <AdminDashboard />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
