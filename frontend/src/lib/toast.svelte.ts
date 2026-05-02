export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration: number;
  timestamp: Date;
}

class ToastManager {
  list = $state<Toast[]>([]);
  history = $state<Toast[]>([]);
  hasUnread = $state(false);

  show(message: string, type: ToastType = 'info', duration = 3000) {
    const id = crypto.randomUUID();
    const t: Toast = { id, message, type, duration, timestamp: new Date() };
    
    // Add to history list instead of popups
    this.history = [t, ...this.history].slice(0, 50);
    this.hasUnread = true;
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

  clearUnread() {
    this.hasUnread = false;
  }

  clearHistory() {
    this.history = [];
    this.hasUnread = false;
  }

  remove(id: string) {
    this.history = this.history.filter(t => t.id !== id);
    if (this.history.length === 0) {
      this.hasUnread = false;
    }
  }
}

export const toast = new ToastManager();
