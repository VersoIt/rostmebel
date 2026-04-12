<script setup lang="ts">
import { ref } from 'vue';
import { LucideCamera, LucideCheckCircle, LucideSend, LucideStar, LucideX } from 'lucide-vue-next';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';

const props = defineProps<{
  projectId?: number;
}>();

const emit = defineEmits(['success']);

const rating = ref(5);
const hoverRating = ref(0);
const comment = ref('');
const phone = ref('');
const images = ref<{ url: string }[]>([]);
const isUploading = ref(false);
const isSubmitting = ref(false);
const isSuccess = ref(false);
const error = ref('');

const setRating = (value: number) => {
  rating.value = value;
};

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (!target.files?.length) return;

  const formData = new FormData();
  formData.append('image', target.files[0]);

  isUploading.value = true;
  error.value = '';
  try {
    const { data } = await api.post('/uploads/images', formData);
    images.value.push({ url: data.url });
  } catch (err) {
    error.value = getApiErrorMessage(err);
  } finally {
    isUploading.value = false;
    target.value = '';
  }
};

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
  phone.value = formatted.substring(0, 18);
};

const submit = async () => {
  if (!phone.value || !comment.value) {
    error.value = 'Заполните телефон и комментарий.';
    return;
  }

  isSubmitting.value = true;
  error.value = '';
  try {
    await api.post('/reviews', {
      project_id: props.projectId,
      client_phone: phone.value,
      rating: rating.value,
      comment: comment.value,
      images: images.value,
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
  <div class="p-5 sm:p-8">
    <div v-if="isSuccess" class="py-10 text-center">
      <div class="mx-auto mb-5 flex h-16 w-16 items-center justify-center rounded-lg bg-green-50 text-green-600">
        <LucideCheckCircle :size="34" />
      </div>
      <h3 class="ui-title-md mb-2">Отзыв отправлен</h3>
      <p class="text-brand-brown/60">Спасибо. Отзыв появится на сайте после модерации.</p>
    </div>

    <div v-else>
      <div class="mb-8 text-center">
        <h3 class="ui-title-md mb-2">Поделитесь впечатлением</h3>
        <p class="text-sm font-medium text-brand-brown/45">Отзывы публикуются после проверки заказа.</p>
      </div>

      <div class="space-y-6">
        <div class="flex flex-col items-center gap-3">
          <span class="ui-label-compact">Оценка</span>
          <div class="flex gap-1">
            <button
              v-for="i in 5"
              :key="i"
              type="button"
              class="rounded-lg p-1 transition-transform hover:scale-105"
              @click="setRating(i)"
              @mouseenter="hoverRating = i"
              @mouseleave="hoverRating = 0"
            >
              <LucideStar
                :size="34"
                :class="[(hoverRating || rating) >= i ? 'fill-brand-gold text-brand-gold' : 'fill-brand-gray text-brand-gray']"
              />
            </button>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="ui-label">Телефон *</label>
            <input v-model="phone" type="tel" placeholder="+7 (___) ___-__-__" class="ui-input" @input="formatPhone">
          </div>
          <div>
            <label class="ui-label">Фото готовой мебели</label>
            <label class="flex min-h-12 cursor-pointer items-center gap-3 rounded-lg border border-dashed border-brand-gold/30 bg-brand-gold/5 px-4 py-3 text-brand-gold transition-colors hover:border-brand-gold">
              <LucideCamera :size="20" />
              <span class="text-sm font-bold">{{ isUploading ? 'Загружаем...' : 'Добавить фото' }}</span>
              <input type="file" class="hidden" accept="image/*" @change="handleFileUpload">
            </label>
          </div>
        </div>

        <div v-if="images.length > 0 || isUploading" class="flex gap-3 overflow-x-auto pb-2 no-scrollbar">
          <div v-for="(image, idx) in images" :key="image.url" class="relative h-20 w-20 shrink-0 overflow-hidden rounded-lg border border-brand-brown/10">
            <img :src="image.url" class="h-full w-full object-cover" alt="">
            <button type="button" class="absolute right-1 top-1 rounded-lg bg-black/50 p-1 text-white transition-colors hover:bg-red-600" @click="images.splice(idx, 1)">
              <LucideX :size="12" />
            </button>
          </div>
          <div v-if="isUploading" class="flex h-20 w-20 shrink-0 items-center justify-center rounded-lg bg-brand-gray">
            <div class="h-6 w-6 animate-spin rounded-full border-2 border-brand-gold border-t-transparent"></div>
          </div>
        </div>

        <div>
          <label class="ui-label">Комментарий *</label>
          <textarea
            v-model="comment"
            rows="4"
            class="ui-input"
            placeholder="Что понравилось в проекте, сборке, материалах или сервисе"
          ></textarea>
        </div>

        <div v-if="error" class="rounded-lg bg-red-50 p-4 text-sm font-semibold text-red-700">
          {{ error }}
        </div>

        <button type="button" :disabled="isSubmitting || isUploading" class="ui-button ui-button-primary w-full" @click="submit">
          <LucideSend v-if="!isSubmitting" :size="19" />
          {{ isSubmitting ? 'Отправляем...' : 'Отправить отзыв' }}
        </button>
      </div>
    </div>
  </div>
</template>
