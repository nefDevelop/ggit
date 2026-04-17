<script lang="ts">
  let { children } = $props();
</script>

<div class="admin-layout-container">
  <main class="admin-content-full">
    {@render children()}
  </main>
</div>

<style>
  .admin-layout-container {
    animation: fadeIn 0.3s ease-out;
    padding: 2rem 4rem;
    max-width: 1600px;
    margin: 0 auto;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .admin-content-full {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  /* Adjust for mobile if needed */
  @media (max-width: 768px) {
    .admin-layout-container {
      padding: 1rem;
    }
  }
</style>
