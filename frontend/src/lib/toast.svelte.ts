export type ToastType = 'success' | 'error' | 'info' | 'warning';

interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration: number;
}

class ToastManager {
  list = $state<Toast[]>([]);

  show(message: string, type: ToastType = 'info', duration = 3000) {
    const id = crypto.randomUUID();
    const toast: Toast = { id, message, type, duration };
    
    this.list = [...this.list, toast];

    setTimeout(() => {
      this.remove(id);
    }, duration);
  }

  success(message: string, duration?: number) {
    this.show(message, 'success', duration);
  }

  error(message: string, duration?: number) {
    this.show(message, 'error', duration);
  }

  info(message: string, duration?: number) {
    this.show(message, 'info', duration);
  }

  warning(message: string, duration?: number) {
    this.show(message, 'warning', duration);
  }

  remove(id: string) {
    this.list = this.list.filter(t => t.id !== id);
  }
}

export const toast = new ToastManager();
