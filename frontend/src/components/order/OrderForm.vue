<script setup lang="ts">
import { ref } from 'vue';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';
import { LucideCheck, LucideSend } from 'lucide-vue-next';

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
  website: '',
});

const isSubmitting = ref(false);
const isSuccess = ref(false);
const error = ref('');

const formatPhone = (event: Event) => {
  const target = event.target as HTMLInputElement;
  let input = target.value.replace(/\D/g, '');
  if (input.startsWith('7')) input = input.substring(1);
  if (input.startsWith('8')) input = input.substring(1);

  let formatted = '+7 ';
  if (input.length > 0) formatted += `(${input.substring(0, 3)}`;
  if (input.length >= 4) formatted += `) ${input.substring(3, 6)}`;
  if (input.length >= 7) formatted += `-${input.substring(6, 8)}`;
  if (input.length >= 9) formatted += `-${input.substring(8, 10)}`;

  form.value.client_phone = formatted.substring(0, 18);
};

const handleSubmit = async () => {
  isSubmitting.value = true;
  error.value = '';

  try {
    await api.post('/orders', {
      ...form.value,
      project_id: props.projectId,
      fingerprint: btoa(navigator.userAgent),
    });
    isSuccess.value = true;
    window.setTimeout(() => emit('success'), 1800);
  } catch (err) {
    error.value = getApiErrorMessage(err);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div>
    <div v-if="isSuccess" class="py-10 text-center">
      <div class="mx-auto mb-5 flex h-16 w-16 items-center justify-center rounded-lg bg-green-50 text-green-600">
        <LucideCheck :size="34" />
      </div>
      <h3 class="ui-title-md mb-2">Заявка отправлена</h3>
      <p class="text-brand-brown/60">Свяжемся с вами и уточним детали проекта.</p>
    </div>

    <form v-else class="space-y-5" @submit.prevent="handleSubmit">
      <input v-model="form.website" type="text" name="website" class="hidden" tabindex="-1" autocomplete="off">

      <div>
        <label for="order-client-name" class="ui-label">Имя *</label>
        <input
          id="order-client-name"
          v-model="form.client_name"
          required
          type="text"
          autocomplete="name"
          class="ui-input"
          placeholder="Иван"
        >
      </div>

      <div>
        <label for="order-client-phone" class="ui-label">Телефон *</label>
        <input
          id="order-client-phone"
          v-model="form.client_phone"
          required
          type="tel"
          autocomplete="tel"
          class="ui-input"
          placeholder="+7 (___) ___-__-__"
          @input="formatPhone"
        >
      </div>

      <div>
        <label for="order-client-email" class="ui-label">Email</label>
        <input
          id="order-client-email"
          v-model="form.client_email"
          type="email"
          autocomplete="email"
          class="ui-input"
          placeholder="ivan@example.com"
        >
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label for="order-project-type" class="ui-label">Что делаем?</label>
          <select id="order-project-type" v-model="form.project_type" class="ui-input">
            <option value="">Выберите вариант</option>
            <option v-for="option in projectTypeOptions" :key="option" :value="option">{{ option }}</option>
          </select>
        </div>

        <div>
          <label for="order-budget-range" class="ui-label">Ориентир по бюджету</label>
          <select id="order-budget-range" v-model="form.budget_range" class="ui-input">
            <option value="">Пока не знаю</option>
            <option v-for="option in budgetOptions" :key="option" :value="option">{{ option }}</option>
          </select>
        </div>
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label for="order-city" class="ui-label">Город</label>
          <input
            id="order-city"
            v-model="form.city"
            type="text"
            autocomplete="address-level2"
            class="ui-input"
            placeholder="Симферополь, Ялта..."
          >
        </div>

        <div>
          <label for="order-contact-method" class="ui-label">Как удобнее связаться?</label>
          <select id="order-contact-method" v-model="form.contact_method" class="ui-input">
            <option v-for="option in contactOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
        </div>
      </div>

      <div>
        <label for="order-comment" class="ui-label">Комментарий</label>
        <textarea
          id="order-comment"
          v-model="form.comment"
          rows="3"
          class="ui-input"
          placeholder="Размеры, сроки, адрес объекта, что важно учесть"
        ></textarea>
      </div>

      <div v-if="error" class="rounded-lg bg-red-50 p-3 text-sm font-semibold text-red-700">
        {{ error }}
      </div>

      <button type="submit" :disabled="isSubmitting" class="ui-button ui-button-primary w-full">
        <LucideSend v-if="!isSubmitting" :size="19" />
        {{ isSubmitting ? 'Отправляем...' : 'Отправить заявку' }}
      </button>
    </form>
  </div>
</template>
