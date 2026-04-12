import { defineStore } from 'pinia';
import { ref } from 'vue';

interface ConfirmOptions {
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  tone?: 'danger' | 'default';
}

export const useConfirmStore = defineStore('confirm', () => {
  const visible = ref(false);
  const title = ref('');
  const message = ref('');
  const confirmLabel = ref('Подтвердить');
  const cancelLabel = ref('Отмена');
  const tone = ref<'danger' | 'default'>('default');
  let resolveFn: ((value: boolean) => void) | null = null;

  const request = (options: ConfirmOptions) => {
    title.value = options.title;
    message.value = options.message;
    confirmLabel.value = options.confirmLabel || 'Подтвердить';
    cancelLabel.value = options.cancelLabel || 'Отмена';
    tone.value = options.tone || 'default';
    visible.value = true;

    return new Promise<boolean>((resolve) => {
      resolveFn = resolve;
    });
  };

  const resolve = (value: boolean) => {
    visible.value = false;
    resolveFn?.(value);
    resolveFn = null;
  };

  return {
    visible,
    title,
    message,
    confirmLabel,
    cancelLabel,
    tone,
    request,
    resolve,
  };
});
