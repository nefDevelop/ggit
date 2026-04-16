<script lang="ts">
  import { page } from '$app/stores';
  import OrgSettings from '$lib/components/OrgSettings.svelte';

  const orgName = $page.params.orgname;
</script>

<svelte:head>
  <title>Settings · {orgName}</title>
</svelte:head>

<div class="page-container">
  <OrgSettings {orgName} />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
