<script lang="ts">
  import Identicon from '../Identicon.svelte';
  let { seed = '', size = 100, type = 'user' }: { seed?: string; size?: number; type?: string } = $props();
</script>
<Identicon {seed} {size} {type} />
