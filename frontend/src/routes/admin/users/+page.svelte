<script lang="ts">
  import AdminUsers from '$lib/components/AdminUsers.svelte';
</script>

<svelte:head>
  <title>User Management · Admin · ggit</title>
</svelte:head>

<div class="page-container">
  <AdminUsers />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
