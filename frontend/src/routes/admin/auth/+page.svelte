<script lang="ts">
  import AdminAuth from '$lib/components/AdminAuth.svelte';
</script>

<svelte:head>
  <title>Login Sources · Admin · ggit</title>
</svelte:head>

<div class="page-container">
  <AdminAuth />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
