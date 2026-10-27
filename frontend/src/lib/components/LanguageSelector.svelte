<script lang="ts">
    import { i18n } from '$lib/i18n.svelte';
    import { getAvailableLanguages } from '$lib/locales';

    const availableLanguages = getAvailableLanguages();

    function changeLanguage(event: Event) {
        const select = event.currentTarget as HTMLSelectElement;
        i18n.init(select.value);
    }
  import Languages from "lucide-svelte/icons/languages";
  import ChevronDown from "lucide-svelte/icons/chevron-down";
</script>

<div class="switches-grid">
    <div class="switch-row-item">
        <div class="item-info">
            <div class="title-with-icon">
                <Languages size={16} />
                <span class="label">{i18n.t('language') || 'System Language'}</span>
            </div>
            <p class="desc">{i18n.t('settings.language_desc') || 'Choose your preferred language for the interface.'}</p>
        </div>
        
        <div class="select-wrapper">
            <select id="lang-select" onchange={changeLanguage} value={i18n.locale}>
                {#each availableLanguages as lang}
                    <option value={lang.code}>{lang.name}</option>
                {/each}
            </select>
            <ChevronDown class="dropdown-icon" size={18} />
        </div>
    </div>
</div>

<style>
    .switches-grid { display: flex; flex-direction: column; background: var(--surface-variant); border: 1px solid var(--md-sys-color-outline-variant); border-radius: 1.5rem; overflow: hidden; box-shadow: var(--shadow-sm); }
    
    .switch-row-item { display: flex; align-items: center; justify-content: space-between; padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 2rem; }
    .switch-row-item:last-child { border-bottom: none; }
    .switch-row-item:hover { background: rgba(255,255,255,0.02); }

    .item-info { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; }
    .title-with-icon { display: flex; align-items: center; gap: 1rem; font-weight: 700; color: var(--text-color); }
    .desc { font-size: 0.85rem; color: var(--text-muted); margin: 0; opacity: 0.8; }
    .mini { font-size: 1.25rem; color: var(--primary-color); }

    .select-wrapper { position: relative; display: flex; align-items: center; }
    
    select {
        appearance: none;
        background: rgba(255,255,255,0.05);
        border: 1px solid var(--md-sys-color-outline-variant);
        color: var(--text-color);
        padding: 0.6rem 2.5rem 0.6rem 1.2rem;
        border-radius: 0.75rem;
        font-family: inherit;
        font-size: 0.9rem;
        font-weight: 700;
        cursor: pointer;
        outline: none;
        transition: all 0.2s;
        min-width: 160px;
    }

    select:focus {
        border-color: var(--primary-color);
        background: rgba(var(--primary-rgb), 0.05);
    }

    .dropdown-icon {
        position: absolute;
        right: 1rem;
        pointer-events: none;
        font-size: 1.25rem;
        color: var(--text-muted);
    }

    @media (max-width: 600px) {
        .switch-row-item { flex-direction: column; align-items: flex-start; gap: 1.5rem; }
        .select-wrapper { width: 100%; }
        select { width: 100%; }
    }
</style>
