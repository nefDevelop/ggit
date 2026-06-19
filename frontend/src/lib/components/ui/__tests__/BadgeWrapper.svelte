<script lang="ts">
  import Badge from '../Badge.svelte';
  let { childText = '', ...rest }: Record<string, any> = $props();
</script>

<Badge {...rest}>
  {childText}
</Badge>
