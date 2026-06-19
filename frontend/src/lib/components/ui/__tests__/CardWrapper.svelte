<script lang="ts">
  import Card from '../Card.svelte';
  let { childText = '', ...rest }: Record<string, any> = $props();
</script>

<Card {...rest}>
  {childText}
</Card>
