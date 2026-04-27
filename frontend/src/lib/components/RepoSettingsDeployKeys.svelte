<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import Modal from "$lib/components/Modal.svelte";
  import Input from "$lib/components/ui/Input.svelte";
  import Checkbox from "$lib/components/ui/Checkbox.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { toast } from "$lib/toast.svelte";

  const repo: any = getContext("repo");

  let keys = $state<any[]>([]);
  let loading = $state(true);
  let error = $state("");

  let showModal = $state(false);
  let newKey = $state({ title: "", key: "", read_only: true });
  let isSubmitting = $state(false);

  async function fetchKeys() {
    loading = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/keys`,
      );
      if (res.ok) {
        keys = await res.json();
      } else {
        error = "Failed to load deploy keys.";
      }
    } catch (e) {
      error = "An error occurred while fetching deploy keys.";
    } finally {
      loading = false;
    }
  }

  async function handleAdd() {
    isSubmitting = true;
    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/keys`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(newKey),
        },
      );

      if (res.ok) {
        const key = await res.json();
        keys = [...keys, key];
        showModal = false;
        newKey = { title: "", key: "", read_only: true };
      } else {
        const data = await res.json();
        toast.error(data.message || "Failed to add deploy key.");
      }
    } catch (e) {
      console.error(e);
      toast.error("An error occurred while adding the key.");
    } finally {
      isSubmitting = false;
    }
  }

  async function handleDelete(id: number) {
    if (!(await dialog.confirm("Are you sure you want to delete this deploy key?"))) return;

    try {
      const res = await apiFetch(
        `/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/keys/${id}`,
        {
          method: "DELETE",
        },
      );
      if (res.ok) {
        keys = keys.filter((k) => k.id !== id);
      }
    } catch (e) {
      console.error(e);
    }
  }

  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleDateString(undefined, {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  }

  onMount(fetchKeys);
</script>

<div class="keys-manager">
  <header class="section-header">
    <div class="header-content">
      <h2>Deploy Keys</h2>
      <p>
        Deploy keys allow external servers to clone this repository via SSH
        without needing a user account.
      </p>
    </div>
    <button class="btn btn-primary" onclick={() => (showModal = true)}>
      <span class="material-symbols-rounded">add</span>
      Add Key
    </button>
  </header>

  {#if loading}
    <div class="skeleton-list">
      <div class="skeleton-item tonal"></div>
      <div class="skeleton-item tonal"></div>
    </div>
  {:else if error}
    <div class="error-banner">
      <span class="material-symbols-rounded">error</span>
      {error}
    </div>
  {:else if keys.length === 0}
    <div class="empty-state tonal">
      <span class="material-symbols-rounded size-xl">vpn_key</span>
      <h3>No Deploy Keys</h3>
      <p>You haven't configured any deploy keys for this repository yet.</p>
      <button class="btn btn-tonal" onclick={() => (showModal = true)}
        >Add your first deploy key</button
      >
    </div>
  {:else}
    <div class="keys-list">
      {#each keys as key}
        <div class="key-item tonal">
          <div class="key-icon-box">
            <span class="material-symbols-rounded">key</span>
          </div>
          <div class="key-info">
            <div class="key-main">
              <span class="key-title">{key.title}</span>
              {#if key.read_only}
                <span class="badge badge-read">Read-Only</span>
              {:else}
                <span class="badge badge-write">Read-Write</span>
              {/if}
            </div>
            <div class="key-fingerprint">
              <code>{key.key.substring(0, 32)}...</code>
            </div>
            <div class="key-meta">
              Added on {formatDate(key.created_at)}
            </div>
          </div>
          <div class="key-actions">
            <button
              class="btn-icon delete"
              onclick={() => handleDelete(key.id)}
              title="Delete"
            >
              <span class="material-symbols-rounded">delete</span>
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <Modal show={showModal} title="Add Deploy Key" onclose={() => (showModal = false)}>
    <form
      onsubmit={(e) => {
        e.preventDefault();
        handleAdd();
      }}
      style="display: flex; flex-direction: column; gap: 1.5rem;"
    >
      <Input 
        label="Key Name" 
        bind:value={newKey.title} 
        placeholder="e.g. CI Server"
      />
      <Input 
        label="Key Content" 
        bind:value={newKey.key} 
        placeholder="ssh-rsa AAAAB3Nza..."
        multiline
      />

        rows={5}
      />

      <Checkbox
        bind:checked={newKey.read_only}
        label="Read-only"
        description="If enabled, this key will not be able to push to the repository."
      />

      <footer class="modal-footer">
        <button
          type="button"
          class="btn btn-tonal"
          onclick={() => (showModal = false)}>Cancel</button
        >
        <button
          type="submit"
          class="btn btn-primary"
          disabled={isSubmitting}
        >
          {#if isSubmitting}
            <span class="spinner"></span>
          {:else}
            Add Key
          {/if}
        </button>
      </footer>
    </form>
  </Modal>
</div>

<style>
  .keys-manager {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    animation: fadeIn 0.3s ease-out;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }
  .header-content h2 {
    margin: 0 0 0.5rem 0;
    font-size: 1.5rem;
    font-weight: 800;
  }
  .header-content p {
    margin: 0;
    color: var(--text-muted);
    font-size: 0.95rem;
  }

  .keys-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .key-item {
    display: flex;
    gap: 1.25rem;
    align-items: center;
    padding: 1.25rem 1.5rem;
    transition:
      transform 0.2s,
      background 0.2s;
  }
  .key-item:hover {
    transform: translateX(5px);
    background: rgba(255, 255, 255, 0.03);
  }

  .key-icon-box {
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
    width: 48px;
    height: 48px;
    border-radius: 1rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .key-icon-box span {
    font-size: 1.75rem;
  }

  .key-info {
    flex: 1;
  }
  .key-main {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.25rem;
  }
  .key-title {
    font-weight: 700;
    font-size: 1.1rem;
  }

  .badge {
    font-size: 0.65rem;
    font-weight: 800;
    text-transform: uppercase;
    padding: 0.1rem 0.5rem;
    border-radius: 4px;
  }
  .badge-read {
    background: rgba(255, 255, 255, 0.1);
    color: var(--text-color);
  }
  .badge-write {
    background: rgba(var(--primary-rgb), 0.2);
    color: var(--primary-color);
  }

  .key-fingerprint {
    margin-bottom: 0.25rem;
  }
  .key-fingerprint code {
    font-size: 0.75rem;
    opacity: 0.6;
    font-family: "JetBrains Mono", monospace;
  }

  .key-meta {
    font-size: 0.8rem;
    color: var(--text-muted);
  }

  .btn-icon {
    background: transparent;
    border: none;
    color: var(--text-muted);
    padding: 0.5rem;
    border-radius: 0.5rem;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
  }
  .btn-icon:hover {
    background: rgba(255, 255, 255, 0.05);
    color: var(--text-color);
  }
  .btn-icon.delete:hover {
    background: rgba(248, 113, 113, 0.1);
    color: #f87171;
  }

  .empty-state {
    padding: 4rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    text-align: center;
  }
  .empty-state .size-xl {
    font-size: 4rem;
    opacity: 0.2;
  }
  .empty-state h3 {
    margin: 1rem 0 0 0;
  }
  .empty-state p {
    color: var(--text-muted);
    max-width: 400px;
    margin-bottom: 1rem;
  }

  .modal-footer {
    padding-top: 1rem;
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
  }

  .skeleton-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .skeleton-item {
    height: 90px;
    width: 100%;
    animation: pulse 1.5s infinite;
  }
  @keyframes pulse {
    0%,
    100% {
      opacity: 0.6;
    }
    50% {
      opacity: 0.3;
    }
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
