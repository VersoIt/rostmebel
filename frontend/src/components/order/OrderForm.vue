<script setup lang="ts">
import { ref } from 'vue';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';
import { LucideSend, LucideCheck } from 'lucide-vue-next';

const props = defineProps<{
  projectId?: number;
}>();

const emit = defineEmits(['success']);

const projectTypeOptions = [
  'Кухня',
  'Шкаф или гардеробная',
  'Мебель для всей квартиры',
  'Коммерческий объект',
  'Пока не знаю',
];

const budgetOptions = [
  'До 200 000 ₽',
  '200 000-400 000 ₽',
  '400 000-700 000 ₽',
  'От 700 000 ₽',
  'Нужен расчет',
];

const contactOptions = [
  { value: 'phone', label: 'Звонок' },
  { value: 'whatsapp', label: 'WhatsApp' },
  { value: 'telegram', label: 'Telegram' },
  { value: 'email', label: 'Email' },
];

const form = ref({
  client_name: '',
  client_phone: '',
  client_email: '',
  project_type: '',
  budget_range: '',
  city: '',
  contact_method: 'phone',
  comment: '',
  website: '', // Honeypot
});

const formatPhone = (e: Event) => {
  const target = e.target as HTMLInputElement;
  let input = target.value.replace(/\D/g, '');
  if (input.startsWith('7')) input = input.substring(1);
  if (input.startsWith('8')) input = input.substring(1);
  
  let formatted = '+7 ';
  if (input.length > 0) formatted += '(' + input.substring(0, 3);
  if (input.length >= 4) formatted += ') ' + input.substring(3, 6);
  if (input.length >= 7) formatted += '-' + input.substring(6, 8);
  if (input.length >= 9) formatted += '-' + input.substring(8, 10);
  
  form.value.client_phone = formatted.substring(0, 18);
};

const isSubmitting = ref(false);
const isSuccess = ref(false);
const error = ref('');

const handleSubmit = async () => {
  isSubmitting.value = true;
  error.value = '';
  
  try {
    await api.post('/orders', {
      ...form.value,
      project_id: props.projectId,
      fingerprint: btoa(navigator.userAgent), // Basic fingerprint
    });
    isSuccess.value = true;
    setTimeout(() => {
      emit('success');
    }, 2000);
  } catch (err: any) {
    error.value = getApiErrorMessage(err);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div class="relative">
    <div v-if="isSuccess" class="flex flex-col items-center justify-center py-12 text-center">
      <div class="w-20 h-20 bg-green-100 text-green-600 rounded-lg flex items-center justify-center mb-6">
        <LucideCheck :size="40" />
      </div>
      <h3 class="text-2xl font-serif text-brand-brown mb-2">Спасибо!</h3>
      <p class="text-brand-brown/60">Ваша заявка успешно отправлена</p>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Honeypot -->
      <input v-model="form.website" type="text" name="website" class="hidden">
      
      <div>
        <label for="order-client-name" class="block text-sm font-medium text-brand-brown/60 mb-2">Имя *</label>
        <input 
          id="order-client-name"
          v-model="form.client_name"
          required
          type="text"
          autocomplete="name"
          class="w-full px-4 py-3 rounded-lg border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30"
          placeholder="Иван Иванов"
        >
      </div>

      <div>
        <label for="order-client-phone" class="block text-sm font-medium text-brand-brown/60 mb-2">Телефон *</label>
        <input 
          id="order-client-phone"
          v-model="form.client_phone"
          @input="formatPhone"
          required
          type="tel"
          autocomplete="tel"
          class="w-full px-4 py-3 rounded-lg border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30"
          placeholder="+7 (___) ___-__-__"
        >
      </div>

      <div>
        <label for="order-client-email" class="block text-sm font-medium text-brand-brown/60 mb-2">Email (необязательно)</label>
        <input 
          id="order-client-email"
          v-model="form.client_email"
          type="email"
          autocomplete="email"
          class="w-full px-4 py-3 rounded-lg border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30"
          placeholder="ivan@example.com"
        >
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label for="order-project-type" class="block text-sm font-medium text-brand-brown/60 mb-2">Что делаем?</label>
          <select
            id="order-project-type"
            v-model="form.project_type"
            class="w-full px-4 py-3 rounded-lg border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30"
          >
            <option value="">Выберите вариант</option>
            <option v-for="option in projectTypeOptions" :key="option" :value="option">{{ option }}</option>
          </select>
        </div>

        <div>
          <label for="order-budget-range" class="block text-sm font-medium text-brand-brown/60 mb-2">Ориентир по бюджету</label>
          <select
            id="order-budget-range"
            v-model="form.budget_range"
            class="w-full px-4 py-3 rounded-lg border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30"
          >
            <option value="">Пока не знаю</option>
            <option v-for="option in budgetOptions" :key="option" :value="option">{{ option }}</option>
          </select>
        </div>
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label for="order-city" class="block text-sm font-medium text-brand-brown/60 mb-2">Город</label>
          <input
            id="order-city"
            v-model="form.city"
            type="text"
            autocomplete="address-level2"
            class="w-full px-4 py-3 rounded-lg border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30"
            placeholder="Симферополь, Ялта..."
          >
        </div>

        <div>
          <label for="order-contact-method" class="block text-sm font-medium text-brand-brown/60 mb-2">Как удобнее связаться?</label>
          <select
            id="order-contact-method"
            v-model="form.contact_method"
            class="w-full px-4 py-3 rounded-lg border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30"
          >
            <option v-for="option in contactOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
        </div>
      </div>

      <div>
        <label for="order-comment" class="block text-sm font-medium text-brand-brown/60 mb-2">Комментарий</label>
        <textarea 
          id="order-comment"
          v-model="form.comment"
          rows="3"
          class="w-full px-4 py-3 rounded-lg border border-brand-brown/10 focus:border-brand-gold outline-none bg-brand-gray/30"
          placeholder="Размеры, сроки, адрес объекта, что важно учесть"
        ></textarea>
      </div>

      <div v-if="error" class="text-red-500 text-sm bg-red-50 p-3 rounded-lg">
        {{ error }}
      </div>

      <button 
        type="submit"
        :disabled="isSubmitting"
        class="w-full bg-brand-brown text-white py-4 rounded-lg font-medium hover:bg-brand-gold disabled:opacity-50 transition-all flex items-center justify-center gap-2"
      >
        <LucideSend v-if="!isSubmitting" :size="20" />
        {{ isSubmitting ? 'Отправка...' : 'Отправить заявку' }}
      </button>
    </form>
  </div>
</template>
