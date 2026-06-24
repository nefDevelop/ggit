<script lang="ts">
  import ViewHeader from '../ui/ViewHeader.svelte';
  let { tabs = [], activeTab = $bindable(''), showSearch = true, noPills = false }: Record<string, any> = $props();
</script>
<ViewHeader {tabs} bind:activeTab {showSearch} {noPills} />
