<script lang="ts">
  import { page } from '$app/stores';
  import TeamDetail from '$lib/components/TeamDetail.svelte';

  const orgName = $page.params.orgname;
  const teamID = $page.params.teamID;
</script>

<svelte:head>
  <title>Manage Team · {orgName}</title>
</svelte:head>

<div class="page page--medium">
  <TeamDetail {orgName} {teamID} />
</div>


