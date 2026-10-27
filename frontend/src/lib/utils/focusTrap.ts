const FOCUSABLE =
  'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])';

/**
 * Atrapa el foco dentro de un diálogo mientras está montado y lo devuelve al
 * elemento anterior al cerrarse. Se usa como acción de Svelte.
 */
export function focusTrap(node: HTMLElement) {
  const previous = document.activeElement as HTMLElement | null;

  function focusable(): HTMLElement[] {
    return Array.from(node.querySelectorAll<HTMLElement>(FOCUSABLE)).filter(
      (el) => el.offsetParent !== null,
    );
  }

  function onKeydown(event: KeyboardEvent) {
    if (event.key !== 'Tab') return;
    const items = focusable();
    if (items.length === 0) {
      event.preventDefault();
      return;
    }
    const first = items[0];
    const last = items[items.length - 1];
    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault();
      last.focus();
    } else if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault();
      first.focus();
    }
  }

  if (!node.hasAttribute('tabindex')) node.setAttribute('tabindex', '-1');
  node.addEventListener('keydown', onKeydown);
  queueMicrotask(() => (focusable()[0] ?? node).focus());

  return {
    destroy() {
      node.removeEventListener('keydown', onKeydown);
      previous?.focus?.();
    },
  };
}
