<script lang="ts">
  import Button from "./ui/Button.svelte";
  

import ShieldCheck from "lucide-svelte/icons/shield-check";
import MessageSquare from "lucide-svelte/icons/message-square";
import AlertCircle from "lucide-svelte/icons/alert-circle";


  
  let { onSubmitReview } = $props<{
    onSubmitReview: (type: number, content: string) => void
  }>();

  let showModal = $state(false);
  let type = $state(2); // 1: Approve, 2: Comment, 3: Reject
  let content = $state("");
</script>

<div class="review-toolbar">
  <Button variant="primary" icon="rate_review" onclick={() => showModal = true}>
    Finalizar Revisión
  </Button>

  {#if showModal}
    <div class="modal-overlay" onclick={() => showModal = false} onkeydown={(e) => { if (e.key === 'Escape') showModal = false; }} role="button" tabindex="-1">
      <div class="modal-content island" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="dialog" tabindex="-1" aria-modal="true">
        <header>
          <h3>Finalizar Revisión</h3>
          <p>Envía tus comentarios y decide el estado del Pull Request.</p>
        </header>

        <div class="review-types">
          <label class="review-type-opt" class:active={type === 2}>
            <input type="radio" name="rev-type" value={2} bind:group={type} />
            <div class="icon-box"><MessageSquare size={20} /></div>
            <div class="info">
              <strong>Comentar</strong>
              <span>Enviar comentarios generales sin aprobación.</span>
            </div>
          </label>

          <label class="review-type-opt approve" class:active={type === 1}>
            <input type="radio" name="rev-type" value={1} bind:group={type} />
            <div class="icon-box"><ShieldCheck size={20} /></div>
            <div class="info">
              <strong>Aprobar</strong>
              <span>Dar el visto bueno para fusionar estos cambios.</span>
            </div>
          </label>

          <label class="review-type-opt reject" class:active={type === 3}>
            <input type="radio" name="rev-type" value={3} bind:group={type} />
            <div class="icon-box"><AlertCircle size={20} /></div>
            <div class="info">
              <strong>Solicitar Cambios</strong>
              <span>Indicar que hay problemas que deben corregirse.</span>
            </div>
          </label>
        </div>

        <div class="review-comment">
           <label for="rev-comment">Comentario general (opcional)</label>
           <textarea id="rev-comment" bind:value={content} placeholder="Escribe un resumen de tu revisión..."></textarea>
        </div>

        <footer class="modal-footer">
          <Button variant="ghost" onclick={() => showModal = false}>Cancelar</Button>
          <Button variant="primary" onclick={() => { onSubmitReview(type, content); showModal = false; }}>
            Enviar Revisión
          </Button>
        </footer>
      </div>
    </div>
  {/if}
</div>

<style>
  .review-toolbar {
    position: sticky;
    top: 1rem;
    z-index: 100;
    display: flex;
    justify-content: flex-end;
    margin-bottom: -3rem;
  }

  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.8);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal-content {
    width: 100%;
    max-width: 600px;
    background: var(--md-sys-color-surface);
    border: 1px solid var(--md-sys-color-outline-variant);
    display: flex;
    flex-direction: column;
    gap: 2rem;
    padding: 2.5rem;
    box-shadow: var(--shadow-xl);
  }

  .modal-content h3 { margin: 0; font-size: 1.5rem; font-weight: 900; }
  .modal-content p { margin: 0.25rem 0 0; color: var(--text-muted); }

  .review-types { display: flex; flex-direction: column; gap: 1rem; }
  
  .review-type-opt {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 1.25rem;
    border-radius: 1.25rem;
    border: 2px solid var(--md-sys-color-outline-variant);
    cursor: pointer;
    transition: all 0.2s;
  }
  .review-type-opt input { display: none; }
  .review-type-opt:hover { background: var(--surface-variant); }
  .review-type-opt.active { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .icon-box { width: 44px; height: 44px; border-radius: 12px; background: var(--surface-variant); display: flex; align-items: center; justify-content: center; color: var(--text-muted); }
  .active .icon-box { background: var(--primary-color); color: white; }

  .approve.active { border-color: var(--success); background: rgba(16, 185, 129, 0.05); }
  .approve.active .icon-box { background: var(--success); }

  .reject.active { border-color: var(--error); background: rgba(239, 68, 68, 0.05); }
  .reject.active .icon-box { background: var(--error); }

  .info { display: flex; flex-direction: column; gap: 0.1rem; }
  .info strong { font-size: 1rem; }
  .info span { font-size: 0.85rem; color: var(--text-muted); }

  .review-comment { display: flex; flex-direction: column; gap: 0.75rem; }
  .review-comment label { font-size: 0.85rem; font-weight: 800; color: var(--text-muted); text-transform: uppercase; }
  .review-comment textarea { background: var(--surface-variant); border: 2px solid transparent; border-radius: 1rem; padding: 1rem; color: var(--text-color); font-family: inherit; min-height: 100px; outline: none; }
  .review-comment textarea:focus { border-color: var(--primary-color); }

  .modal-footer { display: flex; justify-content: flex-end; gap: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); padding-top: 1.5rem; }
</style>
