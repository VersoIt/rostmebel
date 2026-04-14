<script setup lang="ts">
import { computed, ref } from 'vue';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';
import { budgetOptions, contactOptions, projectTypeOptions } from '@/utils/orderOptions';
import { formatRussianPhone } from '@/utils/phone';
import {
  LucideArrowLeft,
  LucideArrowRight,
  LucideCheck,
  LucideCheckCircle2,
  LucideSend,
} from 'lucide-vue-next';

const props = withDefaults(defineProps<{
  projectId?: number;
  initialProjectType?: string;
}>(), {
  initialProjectType: 'Кухня с техникой',
});

const emit = defineEmits(['success']);

const steps = [
  { title: 'Что делаем?', hint: 'Один выбор, чтобы не тратить ваше время.' },
  { title: 'Какой бюджет держим в голове?', hint: 'Можно выбрать “Нужен расчет”.' },
  { title: 'Где и как удобнее связаться?', hint: 'Город нужен, чтобы понимать выезд и логистику.' },
  { title: 'Куда отправить расчет?', hint: 'Мы уточним детали и предложим следующий шаг.' },
] as const;

const currentStep = ref(0);
const isSubmitting = ref(false);
const isSuccess = ref(false);
const error = ref('');

const form = ref({
  client_name: '',
  client_phone: '',
  client_email: '',
  project_type: props.initialProjectType,
  budget_range: '',
  city: '',
  contact_method: 'phone',
  comment: '',
  website: '',
});

const selectedContactLabel = computed(() => {
  return contactOptions.find((option) => option.value === form.value.contact_method)?.label || 'Звонок';
});

const canContinue = computed(() => {
  if (currentStep.value === 0) return Boolean(form.value.project_type);
  if (currentStep.value === 1) return Boolean(form.value.budget_range);
  if (currentStep.value === 2) return Boolean(form.value.contact_method);
  const hasPhone = form.value.client_phone.replace(/\D/g, '').length >= 10;
  const hasEmail = /.+@.+\..+/.test(form.value.client_email);

  return form.value.client_name.trim().length >= 2
    && hasPhone
    && (form.value.contact_method !== 'email' || hasEmail);
});

const progress = computed(() => ((currentStep.value + 1) / steps.length) * 100);

const formatPhone = (event: Event) => {
  const target = event.target as HTMLInputElement;
  form.value.client_phone = formatRussianPhone(target.value);
};

const nextStep = () => {
  if (!canContinue.value || currentStep.value >= steps.length - 1) return;
  currentStep.value += 1;
};

const prevStep = () => {
  if (currentStep.value === 0) return;
  currentStep.value -= 1;
};

const handleSubmit = async () => {
  if (!canContinue.value) return;

  isSubmitting.value = true;
  error.value = '';

  try {
    await api.post('/orders', {
      ...form.value,
      project_id: props.projectId,
      fingerprint: window.btoa(navigator.userAgent),
    });
    isSuccess.value = true;
    window.setTimeout(() => emit('success'), 1200);
  } catch (err) {
    error.value = getApiErrorMessage(err);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div>
    <div v-if="isSuccess" class="py-8 text-center">
      <div class="mx-auto mb-5 flex h-16 w-16 items-center justify-center rounded-lg bg-green-50 text-green-600">
        <LucideCheck :size="34" />
      </div>
      <h3 class="ui-title-md mb-2">Заявка отправлена</h3>
      <p class="text-brand-brown/60">Свяжемся с вами, уточним детали и подскажем следующий шаг.</p>
    </div>

    <form v-else @submit.prevent="handleSubmit">
      <input v-model="form.website" type="text" name="website" class="hidden" tabindex="-1" autocomplete="off">

      <div class="mb-5 flex items-start justify-between gap-4">
        <div>
          <p class="ui-eyebrow mb-2">Расчет за минуту</p>
          <h3 class="font-serif text-2xl font-bold leading-tight text-brand-brown sm:text-3xl">
            {{ steps[currentStep].title }}
          </h3>
          <p class="mt-2 text-sm leading-6 text-brand-brown/55">{{ steps[currentStep].hint }}</p>
        </div>
        <div class="shrink-0 rounded-lg bg-brand-gray px-3 py-2 text-xs font-black text-brand-brown/55">
          {{ currentStep + 1 }}/{{ steps.length }}
        </div>
      </div>

      <div class="mb-6 h-2 overflow-hidden rounded-full bg-brand-gray">
        <div class="h-full rounded-full bg-brand-gold transition-all duration-300" :style="{ width: `${progress}%` }"></div>
      </div>

      <transition name="quiz-step" mode="out-in">
        <div :key="currentStep" class="min-h-[270px]">
          <div v-if="currentStep === 0" class="grid grid-cols-1 gap-2 sm:grid-cols-2">
            <button
              v-for="option in projectTypeOptions"
              :key="option"
              type="button"
              :class="[
                'flex min-h-14 items-center justify-between gap-3 rounded-lg border px-4 py-3 text-left text-sm font-bold transition-all',
                form.project_type === option
                  ? 'border-brand-brown bg-brand-brown text-white shadow-sm'
                  : 'border-brand-brown/10 bg-white text-brand-brown hover:border-brand-gold hover:text-brand-gold'
              ]"
              @click="form.project_type = option"
            >
              {{ option }}
              <LucideCheckCircle2 v-if="form.project_type === option" :size="18" class="shrink-0 text-brand-gold" />
            </button>
          </div>

          <div v-else-if="currentStep === 1" class="grid grid-cols-1 gap-2">
            <button
              v-for="option in budgetOptions"
              :key="option"
              type="button"
              :class="[
                'flex min-h-14 items-center justify-between gap-3 rounded-lg border px-4 py-3 text-left text-sm font-bold transition-all',
                form.budget_range === option
                  ? 'border-brand-brown bg-brand-brown text-white shadow-sm'
                  : 'border-brand-brown/10 bg-white text-brand-brown hover:border-brand-gold hover:text-brand-gold'
              ]"
              @click="form.budget_range = option"
            >
              {{ option }}
              <LucideCheckCircle2 v-if="form.budget_range === option" :size="18" class="shrink-0 text-brand-gold" />
            </button>
          </div>

          <div v-else-if="currentStep === 2" class="space-y-4">
            <div>
              <label for="quote-city" class="ui-label">Город или район</label>
              <input
                id="quote-city"
                v-model="form.city"
                type="text"
                autocomplete="address-level2"
                class="ui-input bg-white"
                placeholder="Симферополь, Ялта, Севастополь..."
              >
            </div>

            <div>
              <div class="ui-label">Как удобнее связаться?</div>
              <div class="grid grid-cols-2 gap-2 sm:grid-cols-3">
                <button
                  v-for="option in contactOptions"
                  :key="option.value"
                  type="button"
                  :class="[
                    'rounded-lg border px-3 py-3 text-sm font-bold transition-all',
                    form.contact_method === option.value
                      ? 'border-brand-brown bg-brand-brown text-white shadow-sm'
                      : 'border-brand-brown/10 bg-white text-brand-brown hover:border-brand-gold hover:text-brand-gold'
                  ]"
                  @click="form.contact_method = option.value"
                >
                  {{ option.label }}
                </button>
              </div>
            </div>
          </div>

          <div v-else class="space-y-4">
            <div class="rounded-lg bg-brand-gray/60 p-4 text-sm leading-6 text-brand-brown/65">
              <span class="font-bold text-brand-brown">{{ form.project_type }}</span>,
              {{ form.budget_range.toLowerCase() }}.
              Связь: {{ selectedContactLabel }}.
            </div>

            <div class="grid gap-4 sm:grid-cols-2">
              <div>
                <label for="quote-name" class="ui-label">Имя *</label>
                <input
                  id="quote-name"
                  v-model="form.client_name"
                  required
                  type="text"
                  autocomplete="name"
                  class="ui-input bg-white"
                  placeholder="Иван"
                >
              </div>

              <div>
                <label for="quote-phone" class="ui-label">Телефон *</label>
                <input
                  id="quote-phone"
                  v-model="form.client_phone"
                  required
                  type="tel"
                  autocomplete="tel"
                  class="ui-input bg-white"
                  placeholder="+7 (___) ___-__-__"
                  @input="formatPhone"
                >
              </div>

              <div class="sm:col-span-2">
                <label for="quote-email" class="ui-label">Email{{ form.contact_method === 'email' ? ' *' : '' }}</label>
                <input
                  id="quote-email"
                  v-model="form.client_email"
                  :required="form.contact_method === 'email'"
                  type="email"
                  autocomplete="email"
                  class="ui-input bg-white"
                  placeholder="ivan@example.com"
                >
              </div>
            </div>

            <div>
              <label for="quote-comment" class="ui-label">Комментарий</label>
              <textarea
                id="quote-comment"
                v-model="form.comment"
                rows="3"
                class="ui-input bg-white"
                placeholder="Размеры, техника, сроки, что важно учесть"
              ></textarea>
            </div>
          </div>
        </div>
      </transition>

      <div v-if="error" class="mt-4 rounded-lg bg-red-50 p-3 text-sm font-semibold text-red-700">
        {{ error }}
      </div>

      <div class="mt-5 flex flex-col gap-2 sm:flex-row">
        <button
          v-if="currentStep > 0"
          type="button"
          class="ui-button ui-button-secondary sm:w-auto"
          @click="prevStep"
        >
          <LucideArrowLeft :size="18" />
          Назад
        </button>

        <button
          v-if="currentStep < steps.length - 1"
          type="button"
          :disabled="!canContinue"
          class="ui-button ui-button-primary flex-1"
          @click="nextStep"
        >
          Дальше
          <LucideArrowRight :size="18" />
        </button>

        <button v-else type="submit" :disabled="isSubmitting || !canContinue" class="ui-button ui-button-primary flex-1">
          <LucideSend v-if="!isSubmitting" :size="18" />
          {{ isSubmitting ? 'Отправляем...' : 'Получить расчет' }}
        </button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.quiz-step-enter-active,
.quiz-step-leave-active {
  transition: opacity 180ms ease, transform 180ms ease;
}

.quiz-step-enter-from {
  opacity: 0;
  transform: translateX(10px);
}

.quiz-step-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}
</style>
