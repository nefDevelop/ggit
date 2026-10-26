<script lang="ts">
  import { page } from '$app/stores';
  import OrgSettings from '$lib/components/OrgSettings.svelte';

  const orgName = $page.params.orgname;
</script>

<svelte:head>
  <title>Settings · {orgName}</title>
</svelte:head>

<div class="page page--medium">
  <OrgSettings {orgName} />
</div>


