<script lang="ts">
  import AdminNotices from '$lib/components/AdminNotices.svelte';
</script>

<svelte:head>
  <title>System Notices · Admin · ggit</title>
</svelte:head>

<div class="page-container">
  <AdminNotices />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
