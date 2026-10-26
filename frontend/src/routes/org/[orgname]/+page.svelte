<script lang="ts">
  import { page } from '$app/stores';
  import OrgView from '$lib/components/OrgView.svelte';

  const orgName = $page.params.orgname;
</script>

<svelte:head>
  <title>{orgName} · Organization</title>
</svelte:head>

<div class="page page--medium">
  <OrgView {orgName} />
</div>


