<script lang="ts">
  import { page } from '$app/stores';
  import OrgView from '$lib/components/OrgView.svelte';

  const orgName = $page.params.orgname;
</script>

<svelte:head>
  <title>{orgName} · Organization</title>
</svelte:head>

<div class="page-container">
  <OrgView {orgName} />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
