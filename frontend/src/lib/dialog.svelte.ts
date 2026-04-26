type DialogType = 'alert' | 'confirm' | 'prompt';

interface DialogOptions {
  title?: string;
  message: string;
  defaultValue?: string;
  confirmText?: string;
  cancelText?: string;
  type: DialogType;
  resolve: (value: any) => void;
}

class DialogManager {
  current = $state<DialogOptions | null>(null);

  alert(message: string, title = 'Alerta'): Promise<void> {
    return new Promise((resolve) => {
      this.current = { type: 'alert', message, title, resolve };
    });
  }

  confirm(message: string, title = 'Confirmar'): Promise<boolean> {
    return new Promise((resolve) => {
      this.current = { type: 'confirm', message, title, resolve };
    });
  }

  prompt(message: string, defaultValue = '', title = 'Entrada'): Promise<string | null> {
    return new Promise((resolve) => {
      this.current = { type: 'prompt', message, title, defaultValue, resolve };
    });
  }

  close(value: any = null) {
    if (this.current) {
      this.current.resolve(value);
      this.current = null;
    }
  }
}

export const dialog = new DialogManager();
