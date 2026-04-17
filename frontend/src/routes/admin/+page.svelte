<script lang="ts">
  import { onMount } from 'svelte';
  
  onMount(() => {
    window.location.href = '/admin/users';
  });
</script>

<div class="redirecting">
  <p>Loading Admin Panel...</p>
</div>
