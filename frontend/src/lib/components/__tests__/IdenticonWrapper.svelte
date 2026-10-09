<script lang="ts">
  import Identicon from '../Identicon.svelte';
  let { seed = '', size = 100, type = 'user', langColor = '' }: { seed?: string; size?: number; type?: string; langColor?: string } = $props();
</script>
<Identicon {seed} {size} {type} {langColor} />
