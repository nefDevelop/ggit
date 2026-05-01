<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { toast } from "$lib/toast.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import Button from "./ui/Button.svelte";
  import Input from "./ui/Input.svelte";
  import Card from "./ui/Card.svelte";
  import PageHeader from "./ui/PageHeader.svelte";

  let loading = $state(true);
  let patterns = $state<any[]>([]);
  let newPattern = $state({
    name: "",
    regex: "",
    description: ""
  });

  async function fetchPatterns() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/admin/security/guardian/patterns");
      if (res.ok) patterns = (await res.json()) || [];
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addPattern() {
    if (!newPattern.name || !newPattern.regex) return;
    try {
      const res = await apiFetch("/api/v2/admin/security/guardian/patterns", {
        method: "POST",
        body: JSON.stringify(newPattern)
      });
      if (res.ok) {
        toast.success("Global pattern added.");
        newPattern = { name: "", regex: "", description: "" };
        fetchPatterns();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function deletePattern(id: number) {
    if (!await dialog.confirm("¿Estás seguro de eliminar este patrón global? Se eliminará de todos los repositorios.")) return;
    try {
      const res = await apiFetch(`/api/v2/admin/security/guardian/patterns/${id}`, {
        method: "DELETE"
      });
      if (res.ok) {
        toast.success("Global pattern deleted.");
        fetchPatterns();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchPatterns);
</script>

<div class="admin-security">
<div class="admin-security">
  <div class="admin-security-container">
    <PageHeader 
      title={i18n.t("admin.security_global") || "Seguridad Global"}
      subtitle={i18n.t("admin.security_global_desc") || "Gestión de patrones globales para repositorios"}
      icon="security"
    />

    <div class="admin-security-body">

    <div class="patterns-grid">
      {#if loading}
        <div class="skeleton-list"></div>
      {:else}
        {#each patterns as pattern}
          <div class="pattern-card tonal">
            <div class="pattern-header">
              <span class="pattern-name">{pattern.Name}</span>
              <Button variant="ghost" icon="delete" class="danger-text" onclick={() => deletePattern(pattern.ID)} />
            </div>
            <code class="regex-box">{pattern.Regex}</code>
            <p class="pattern-desc">{pattern.Description || i18n.t("admin.no_description")}</p>
          </div>
        {/each}
      {/if}
    </div>

    <div class="add-global-pattern tonal">
      <h4>{i18n.t("admin.add_global_pattern")}</h4>
      <div class="form-row">
        <Input label="Nombre" bind:value={newPattern.name} />
        <Input label="Regex" bind:value={newPattern.regex} />
      </div>
      <Input label="Descripción" bind:value={newPattern.description} />
      <div class="form-actions">
        <Button variant="primary" icon="add" onclick={addPattern} disabled={!newPattern.name || !newPattern.regex}>
          {i18n.t("admin.save_global_pattern")}
        </Button>
      </div>
    </div>
    </div>
  </div>
</div>

<style>
  .admin-security-body {
    padding: 0 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  .admin-header { display: flex; gap: 1.5rem; align-items: center; margin-bottom: 2.5rem; }
  .admin-header .icon { font-size: 3rem; color: var(--primary-color); }
  .admin-header h3 { margin: 0; font-size: 1.5rem; font-weight: 900; }
  .admin-header p { margin: 0.25rem 0 0; color: var(--text-muted); }

  .patterns-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(350px, 1fr)); gap: 1.5rem; margin-bottom: 3rem; }
  .pattern-card { padding: 1.5rem; border-radius: 1.5rem; display: flex; flex-direction: column; gap: 1rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .pattern-header { display: flex; justify-content: space-between; align-items: center; }
  .pattern-name { font-weight: 850; font-size: 1.1rem; }
  
  .regex-box { font-family: var(--font-mono); font-size: 0.85rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.08); padding: 0.75rem; border-radius: 0.75rem; word-break: break-all; }
  .pattern-desc { margin: 0; font-size: 0.9rem; color: var(--text-muted); }

  .add-global-pattern { padding: 2rem; border-radius: 2rem; border: 2px dashed var(--md-sys-color-outline-variant); }
  .add-global-pattern h4 { margin: 0 0 1.5rem; font-weight: 900; }
  .form-row { display: grid; grid-template-columns: 1fr 2fr; gap: 1rem; margin-bottom: 1rem; }
  .form-actions { display: flex; justify-content: flex-end; margin-top: 1.5rem; }

  @media (max-width: 800px) {
    .patterns-grid { grid-template-columns: 1fr; }
    .form-row { grid-template-columns: 1fr; }
  }
</style>
