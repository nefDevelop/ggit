<script lang="ts">
  import Button from '../../Button.svelte';
  let { props = {}, childText = '' }: { props: Record<string, any>; childText: string } = $props();
</script>

<Button {...props}>
  {childText}
</Button>
