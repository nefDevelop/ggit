<script lang="ts">
  import type { Component } from 'svelte';
  let { component: Comp, props = {}, childText = '' }: { component: Component; props: Record<string, any>; childText: string } = $props();
</script>

<Comp {...props}>
  {childText}
</Comp>
