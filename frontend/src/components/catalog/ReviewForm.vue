<script setup lang="ts">
import { ref } from 'vue';
import { LucideStar, LucideSend, LucideCamera, LucideX, LucideCheckCircle } from 'lucide-vue-next';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';
import { PLACEHOLDER_IMAGE } from '@/utils/constants';

const props = defineProps<{
  projectId?: number;
}>();

const emit = defineEmits(['success']);

const rating = ref(5);
const hoverRating = ref(0);
const comment = ref('');
const phone = ref('');
const images = ref<{url: string}[]>([]);
const isUploading = ref(false);
const isSubmitting = ref(false);
const isSuccess = ref(false);
const error = ref('');

const setRating = (val: number) => rating.value = val;

const handleFileUpload = async (e: Event) => {
  const target = e.target as HTMLInputElement;
  if (!target.files?.length) return;

  const formData = new FormData();
  formData.append('image', target.files[0]);

  isUploading.value = true;
  try {
    const { data } = await api.post('/uploads/images', formData);
    images.value.push({ url: data.url });
  } catch (err) {
    error.value = getApiErrorMessage(err);
  } finally {
    isUploading.value = false;
  }
};

const formatPhone = (e: any) => {
  let input = e.target.value.replace(/\D/g, '');
  if (input.startsWith('7')) input = input.substring(1);
  if (input.startsWith('8')) input = input.substring(1);
  let formatted = '+7 ';
  if (input.length > 0) formatted += '(' + input.substring(0, 3);
  if (input.length >= 4) formatted += ') ' + input.substring(3, 6);
  if (input.length >= 7) formatted += '-' + input.substring(6, 8);
  if (input.length >= 9) formatted += '-' + input.substring(8, 10);
  phone.value = formatted.substring(0, 18);
};

const submit = async () => {
  if (!phone.value || !comment.value) {
    error.value = 'Пожалуйста, заполните все обязательные поля';
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
      images: images.value
    });
    isSuccess.value = true;
    setTimeout(() => emit('success'), 2000);
  } catch (err: any) {
    error.value = getApiErrorMessage(err);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <div class="p-8">
    <div v-if="isSuccess" class="py-12 text-center">
      <div class="w-20 h-20 bg-green-100 text-green-600 rounded-full flex items-center justify-center mx-auto mb-6">
        <LucideCheckCircle :size="40" />
      </div>
      <h3 class="text-3xl font-serif text-brand-brown mb-2">Отзыв отправлен</h3>
      <p class="text-brand-brown/60">Благодарим за обратную связь! Отзыв появится на сайте после модерации.</p>
    </div>

    <div v-else>
      <div class="mb-10 text-center">
        <h3 class="text-3xl font-serif text-brand-brown mb-2">Поделитесь впечатлениями</h3>
        <p class="text-brand-brown/40 text-sm font-medium">Отзыв могут оставить только клиенты с завершенным заказом</p>
      </div>

      <div class="space-y-8">
        <!-- Stars -->
        <div class="flex flex-col items-center gap-3">
          <span class="text-xs font-black uppercase tracking-widest text-brand-brown/30">Ваша оценка</span>
          <div class="flex gap-2">
            <button 
              v-for="i in 5" :key="i"
              @click="setRating(i)"
              @mouseenter="hoverRating = i"
              @mouseleave="hoverRating = 0"
              class="transition-all duration-300 transform hover:scale-125"
            >
              <LucideStar 
                :size="36" 
                :class="[ (hoverRating || rating) >= i ? 'text-brand-gold fill-brand-gold' : 'text-brand-gray fill-brand-gray' ]"
              />
            </button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-[10px] font-black uppercase tracking-widest text-brand-brown/40 mb-2">Номер телефона *</label>
            <input 
              v-model="phone"
              @input="formatPhone"
              type="tel"
              placeholder="+7 (___) ___-__-__"
              class="w-full px-6 py-4 rounded-2xl bg-brand-gray/30 border border-transparent focus:border-brand-gold outline-none font-bold"
            >
          </div>
          <div>
            <label class="block text-[10px] font-black uppercase tracking-widest text-brand-brown/40 mb-2">Фото готовой мебели</label>
            <label class="cursor-pointer group flex items-center gap-3 px-6 py-4 rounded-2xl bg-brand-gold/5 border-2 border-dashed border-brand-gold/20 hover:border-brand-gold transition-all">
              <LucideCamera class="text-brand-gold" :size="20" />
              <span class="text-sm font-bold text-brand-gold">Добавить фото</span>
              <input type="file" class="hidden" accept="image/*" @change="handleFileUpload">
            </label>
          </div>
        </div>

        <!-- Image Previews -->
        <div v-if="images.length > 0" class="flex gap-3 overflow-x-auto pb-2">
          <div v-for="(img, idx) in images" :key="idx" class="relative w-20 h-20 rounded-xl overflow-hidden shrink-0 border border-brand-brown/5">
            <img :src="img.url" class="w-full h-full object-cover">
            <button @click="images.splice(idx, 1)" class="absolute top-1 right-1 bg-black/50 text-white rounded-full p-1 hover:bg-red-500">
              <LucideX :size="12" />
            </button>
          </div>
          <div v-if="isUploading" class="w-20 h-20 rounded-xl bg-brand-gray/30 flex items-center justify-center animate-pulse">
            <div class="w-6 h-6 border-2 border-brand-gold border-t-transparent rounded-full animate-spin"></div>
          </div>
        </div>

        <div>
          <label class="block text-[10px] font-black uppercase tracking-widest text-brand-brown/40 mb-2">Ваш комментарий *</label>
          <textarea 
            v-model="comment"
            rows="4"
            placeholder="Расскажите о качестве сборки, материалах и сервисе..."
            class="w-full px-6 py-4 rounded-2xl bg-brand-gray/30 border border-transparent focus:border-brand-gold outline-none font-medium"
          ></textarea>
        </div>

        <div v-if="error" class="bg-red-50 text-red-600 p-4 rounded-xl text-xs font-bold text-center">
          {{ error }}
        </div>

        <button 
          @click="submit"
          :disabled="isSubmitting || isUploading"
          class="w-full bg-brand-brown text-white py-6 rounded-[2rem] text-lg font-black uppercase tracking-widest hover:bg-brand-gold shadow-2xl transition-all disabled:opacity-50"
        >
          {{ isSubmitting ? 'Отправка...' : 'Опубликовать отзыв' }}
        </button>
      </div>
    </div>
  </div>
</template>
