<script lang="ts">
  import { apiFetch } from "$lib/api";
    import { onMount } from 'svelte';
    import { authState } from '$lib/auth.svelte';
    import { Globe, Lock, Shield, Info, ArrowUpRight, Save, Database, User, Users } from 'lucide-svelte';

    let organizations = $state<any[]>([]);
    let loading = $state(false);
    let error = $state('');
    let success = $state(false);

    let form = $state({
        clone_addr: '',
        auth_username: '',
        auth_password: '',
        uid: 0,
        repo_name: '',
        description: '',
        private: false,
        unlisted: false
    });

    onMount(async () => {
        form.uid = authState.user?.id || 0;
        await fetchOrganizations();
    });

    async function fetchOrganizations() {
        try {
            const res = await apiFetch('/api/v2/user/settings/orgs');
            if (res.ok) {
                organizations = await res.json();
            }
        } catch (e) {
            console.error('Error fetching orgs:', e);
        }
    }

    async function handleMigrate() {
        loading = true;
        error = '';
        success = false;

        try {
            const res = await apiFetch('/api/v2/user/repos/migrate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(form)
            });

            if (res.ok) {
                const repo = await res.json();
                success = true;
                // Redirect after a short delay
                setTimeout(() => {
                    window.location.href = `/${repo.owner.username}/${repo.name}`;
                }, 1500);
            } else {
                const data = await res.json();
                error = data.message || 'Error al migrar el repositorio';
            }
        } catch (e) {
            error = 'Error de conexión con el servidor';
        } finally {
            loading = false;
        }
    }

    $effect(() => {
        // Auto-fill repo name from URL if empty
        if (form.clone_addr && !form.repo_name) {
            const match = form.clone_addr.match(/\/([^\/.]+)(\.git)?$/);
            if (match) {
                form.repo_name = match[1];
            }
        }
    });

</script>

<div class="migrate-container">
    <div class="migrate-header">
        <div class="header-icon">
            <ArrowUpRight size={32} />
        </div>
        <div class="header-text">
            <h1>Migrar Repositorio</h1>
            <p>Importa un repositorio existente desde GitHub, GitLab, Gitea o cualquier URL de Git.</p>
        </div>
    </div>

    {#if success}
        <div class="success-alert">
            <div class="alert-icon">✨</div>
            <div class="alert-content">
                <h3>¡Migración Iniciada!</h3>
                <p>El repositorio se ha creado. El contenido se importará en segundo plano. Redirigiendo...</p>
            </div>
        </div>
    {/if}

    {#if error}
        <div class="error-alert">
            <Info size={20} />
            <p>{error}</p>
        </div>
    {/if}

    <div class="migrate-grid">
        <div class="form-section main-info">
            <div class="section-title">
                <Globe size={18} />
                <h2>Origen del Clon</h2>
            </div>

            <div class="input-group">
                <label for="clone_addr">URL de Clonado</label>
                <input 
                    type="text" 
                    id="clone_addr" 
                    placeholder="https://github.com/usuario/repo.git"
                    bind:value={form.clone_addr}
                    required
                />
                <span class="help-text">Protocolos soportados: HTTP(S), Git.</span>
            </div>

            <div class="auth-optional">
                <div class="auth-toggle">
                    <Shield size={16} />
                    <span>Autenticación (Opcional)</span>
                </div>
                <div class="auth-fields">
                    <div class="input-pair">
                        <div class="input-group">
                            <label for="auth_username">Usuario</label>
                            <input type="text" id="auth_username" bind:value={form.auth_username} />
                        </div>
                        <div class="input-group">
                            <label for="auth_password">Contraseña / Token</label>
                            <input type="password" id="auth_password" bind:value={form.auth_password} />
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="form-section destination-info">
            <div class="section-title">
                <Database size={18} />
                <h2>Destino en Ggit</h2>
            </div>

            <div class="owner-repo-pair">
                <div class="input-group owner-select">
                    <label for="owner">Propietario</label>
                    <select id="owner" bind:value={form.uid}>
                        <option value={authState.user?.id}>
                            {authState.user?.username} (Tú)
                        </option>
                        {#each organizations as org}
                            <option value={org.id}>{org.username}</option>
                        {/each}
                    </select>
                </div>
                <div class="input-group repo-name">
                    <label for="repo_name">Nombre del Repositorio</label>
                    <input type="text" id="repo_name" bind:value={form.repo_name} required />
                </div>
            </div>

            <div class="input-group">
                <label for="description">Descripción (Opcional)</label>
                <textarea id="description" rows="2" bind:value={form.description}></textarea>
            </div>
        </div>

        <div class="form-section settings-section">
            <div class="section-title">
                <Save size={18} />
                <h2>Ajustes de Migración</h2>
            </div>

            <div class="checkbox-options">
                <label class="checkbox-item">
                    <input type="checkbox" bind:checked={form.private} />
                    <div class="checkbox-content">
                        <span class="label">Repositorio Privado</span>
                        <span class="desc">Solo tú y los colaboradores podrán verlo.</span>
                    </div>
                </label>

                <label class="checkbox-item">
                    <input type="checkbox" bind:checked={form.unlisted} />
                    <div class="checkbox-content">
                        <span class="label">No listado (Unlisted)</span>
                        <span class="desc">Ocultar de la página de exploración.</span>
                    </div>
                </label>
            </div>
        </div>

        <div class="actions">
            <button class="btn-migrate" onclick={handleMigrate} disabled={loading || !form.clone_addr || !form.repo_name}>
                {#if loading}
                    <div class="spinner"></div>
                    Migrando...
                {:else}
                    <ArrowUpRight size={18} />
                    Migrar Repositorio
                {/if}
            </button>
            <button class="btn-cancel" onclick={() => window.location.href = '/'}>Cancelar</button>
        </div>
    </div>
</div>

<style>
    .migrate-container {
        max-width: 900px;
        margin: 2rem auto;
        padding: 0 1rem;
    }

    .migrate-header {
        display: flex;
        gap: 1.5rem;
        align-items: center;
        margin-bottom: 3rem;
    }

    .header-icon {
        width: 64px;
        height: 64px;
        background: var(--primary-container);
        color: var(--on-primary-container);
        border-radius: 18px;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: var(--shadow-md);
    }

    .header-text h1 {
        font-size: 2rem;
        font-weight: 600;
        margin: 0;
        color: var(--md-sys-color-on-surface);
    }

    .header-text p {
        color: var(--md-sys-color-on-surface-variant);
        margin: 0.25rem 0 0 0;
    }

    .migrate-grid {
        display: grid;
        grid-template-columns: 1fr;
        gap: 2rem;
    }

    .form-section {
        background: var(--island island-interactive-bg);
        border: 1px solid var(--md-sys-color-outline-variant);
        border-radius: 24px;
        padding: 2rem;
        box-shadow: var(--shadow-sm);
    }

    .section-title {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        margin-bottom: 2rem;
        color: var(--primary);
    }

    .section-title h2 {
        font-size: 1.25rem;
        margin: 0;
        color: var(--md-sys-color-on-surface);
    }

    .input-group {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        margin-bottom: 1.5rem;
    }

    .input-group label {
        font-size: 0.9rem;
        font-weight: 500;
        color: var(--md-sys-color-on-surface-variant);
        padding-left: 0.25rem;
    }

    input[type="text"], input[type="password"], textarea, select {
        background: var(--surface-variant);
        border: 1px solid var(--md-sys-color-outline-variant);
        border-radius: 12px;
        padding: 0.75rem 1rem;
        color: var(--md-sys-color-on-surface);
        font-size: 1rem;
        transition: all 0.2s;
    }

    input:focus, textarea:focus, select:focus {
        border-color: var(--primary);
        outline: none;
        box-shadow: 0 0 0 4px rgba(var(--primary-rgb), 0.1);
    }

    .help-text {
        font-size: 0.8rem;
        color: var(--text-tertiary);
    }

    .auth-optional {
        margin-top: 1rem;
        padding-top: 1.5rem;
        border-top: 1px solid var(--md-sys-color-outline-variant);
    }

    .auth-toggle {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.9rem;
        color: var(--md-sys-color-on-surface-variant);
        margin-bottom: 1rem;
    }

    .input-pair {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }

    .owner-repo-pair {
        display: grid;
        grid-template-columns: 200px 1fr;
        gap: 1rem;
    }

    .checkbox-options {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
    }

    .checkbox-item {
        display: flex;
        gap: 1rem;
        cursor: pointer;
        align-items: flex-start;
    }

    .checkbox-item input {
        width: 20px;
        height: 20px;
        margin-top: 0.25rem;
        cursor: pointer;
    }

    .checkbox-content .label {
        display: block;
        font-weight: 500;
        color: var(--md-sys-color-on-surface);
    }

    .checkbox-content .desc {
        display: block;
        font-size: 0.85rem;
        color: var(--md-sys-color-on-surface-variant);
    }

    .actions {
        display: flex;
        gap: 1rem;
        margin-bottom: 4rem;
    }

    .btn-migrate {
        flex: 1;
        background: var(--primary);
        color: white;
        border: none;
        border-radius: 16px;
        padding: 1rem;
        font-size: 1.1rem;
        font-weight: 600;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.75rem;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-migrate:hover:not(:disabled) {
        background: var(--primary-hover);
        transform: translateY(-2px);
        box-shadow: var(--shadow-lg);
    }

    .btn-migrate:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .btn-cancel {
        background: transparent;
        border: 1px solid var(--md-sys-color-outline-variant);
        color: var(--md-sys-color-on-surface-variant);
        border-radius: 16px;
        padding: 0 2rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-cancel:hover {
        background: var(--surface-variant);
    }

    .success-alert {
        background: #d4edda;
        border-left: 6px solid #28a745;
        color: #155724;
        padding: 1.5rem;
        border-radius: 16px;
        margin-bottom: 2rem;
        display: flex;
        gap: 1rem;
    }

    .alert-icon {
        font-size: 2rem;
    }

    .alert-content h3 {
        margin: 0;
        font-size: 1.1rem;
    }

    .alert-content p {
        margin: 0.25rem 0 0 0;
    }

    .error-alert {
        background: #f8d7da;
        border-left: 6px solid #dc3545;
        color: #721c24;
        padding: 1rem;
        border-radius: 12px;
        margin-bottom: 2rem;
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .spinner {
        width: 20px;
        height: 20px;
        border: 3px solid rgba(255,255,255,0.3);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }

    @media (max-width: 640px) {
        .owner-repo-pair, .input-pair {
            grid-template-columns: 1fr;
        }
        
        .actions {
            flex-direction: column;
        }
    }
</style>
