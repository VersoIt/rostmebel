import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useNotificationStore = defineStore('notifications', () => {
  const message = ref('');
  const type = ref<'success' | 'error' | 'info'>('info');
  const visible = ref(false);
  let timeout: any = null;

  const show = (msg: string, t: 'success' | 'error' | 'info' = 'info') => {
    message.value = msg;
    type.value = t;
    visible.value = true;
    
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(() => {
      visible.value = false;
    }, 5000);
  };

  return { message, type, visible, show };
});
