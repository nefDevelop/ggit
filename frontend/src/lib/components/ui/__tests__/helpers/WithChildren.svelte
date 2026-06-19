<script lang="ts">
  import type { Component } from 'svelte';
  let { comp: Comp, childText = '', ...rest }: Record<string, any> = $props();
</script>

<Comp {...rest}>
  {childText}
</Comp>
