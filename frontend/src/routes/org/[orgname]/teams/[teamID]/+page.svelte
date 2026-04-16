<script lang="ts">
  import { page } from '$app/stores';
  import TeamDetail from '$lib/components/TeamDetail.svelte';

  const orgName = $page.params.orgname;
  const teamID = $page.params.teamID;
</script>

<svelte:head>
  <title>Manage Team · {orgName}</title>
</svelte:head>

<div class="page-container">
  <TeamDetail {orgName} {teamID} />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
