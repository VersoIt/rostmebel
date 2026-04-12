<script setup lang="ts">
import { ref, watch } from 'vue';
import {
  LucideArrowLeft,
  LucideArrowRight,
  LucideCheck,
  LucidePlus,
  LucideTrash2,
  LucideUpload,
  LucideX,
} from 'lucide-vue-next';
import type { Category, Image, Product, ProjectStatus } from '@/types';
import api from '@/api/client';
import { getApiErrorMessage } from '@/api/errors';
import { useNotificationStore } from '@/stores/notifications';

const props = defineProps<{
  product?: Product | null;
  categories: Category[];
}>();

const emit = defineEmits<{
  close: [];
  saved: [];
}>();

type ProjectForm = {
  project_category_id?: number;
  name: string;
  slug: string;
  description: string;
  price: number | null;
  price_old?: number | null;
  images: Image[];
  specs: Record<string, string>;
  ai_tags: string;
  status: ProjectStatus;
};

type ProjectPayload = {
  project_category_id?: number;
  name: string;
  slug: string;
  description: string;
  price: number;
  price_old?: number;
  images: Image[];
  specs: Record<string, string>;
  ai_tags: string;
  status: ProjectStatus;
};

const notificationStore = useNotificationStore();

const emptyForm = (): ProjectForm => ({
  name: '',
  slug: '',
  description: '',
  price: 0,
  price_old: undefined,
  project_category_id: undefined,
  status: 'draft',
  ai_tags: '',
  images: [],
  specs: {},
});

const normalizeImages = (images: Image[] = []): Image[] => {
  const normalized = images
    .map((image) => ({
      url: String(image.url || '').trim(),
      is_main: Boolean(image.is_main),
    }))
    .filter((image) => image.url.length > 0);

  if (!normalized.length) return [];

  const mainIndex = normalized.findIndex((image) => image.is_main);
  normalized.forEach((image, index) => {
    image.is_main = mainIndex >= 0 ? index === mainIndex : index === 0;
  });

  return normalized;
};

const normalizeSpecs = (specs: Record<string, string> = {}) => {
  return Object.entries(specs).reduce<Record<string, string>>((acc, [key, value]) => {
    const cleanKey = key.trim();
    const cleanValue = String(value ?? '').trim();
    if (cleanKey && cleanValue) acc[cleanKey] = cleanValue;
    return acc;
  }, {});
};

const toOptionalNumber = (value: unknown) => {
  if (value === undefined || value === null || value === '') return undefined;
  const numberValue = Number(value);
  return Number.isFinite(numberValue) ? numberValue : undefined;
};

const toForm = (product?: Product | null): ProjectForm => {
  if (!product) return emptyForm();

  return {
    project_category_id: product.project_category_id,
    name: product.name || '',
    slug: product.slug || '',
    description: product.description || '',
    price: product.price ?? 0,
    price_old: product.price_old,
    images: normalizeImages(product.images),
    specs: normalizeSpecs(product.specs),
    ai_tags: product.ai_tags || '',
    status: product.status || 'draft',
  };
};

const buildProjectPayload = (): ProjectPayload => ({
  project_category_id: toOptionalNumber(form.value.project_category_id),
  name: form.value.name.trim(),
  slug: form.value.slug.trim(),
  description: form.value.description.trim(),
  price: toOptionalNumber(form.value.price) ?? 0,
  price_old: toOptionalNumber(form.value.price_old),
  images: normalizeImages(form.value.images),
  specs: normalizeSpecs(form.value.specs),
  ai_tags: form.value.ai_tags.trim(),
  status: form.value.status || 'draft',
});

const form = ref<ProjectForm>(emptyForm());
const newSpecKey = ref('');
const newSpecValue = ref('');
const fileInput = ref<HTMLInputElement | null>(null);
const isUploading = ref(false);
const isSaving = ref(false);
const errorMessage = ref('');

watch(
  () => props.product,
  (product) => {
    form.value = toForm(product);
    newSpecKey.value = '';
    newSpecValue.value = '';
    errorMessage.value = '';
  },
  { immediate: true },
);

const addSpec = () => {
  const key = newSpecKey.value.trim();
  const value = newSpecValue.value.trim();
  if (!key || !value) return;

  form.value.specs = {
    ...form.value.specs,
    [key]: value,
  };
  newSpecKey.value = '';
  newSpecValue.value = '';
};

const removeSpec = (key: string) => {
  const nextSpecs = { ...form.value.specs };
  delete nextSpecs[key];
  form.value.specs = nextSpecs;
};

const moveImage = (idx: number, direction: 'left' | 'right') => {
  const newIdx = direction === 'left' ? idx - 1 : idx + 1;
  if (newIdx < 0 || newIdx >= form.value.images.length) return;

  const images = [...form.value.images];
  [images[idx], images[newIdx]] = [images[newIdx], images[idx]];
  form.value.images = normalizeImages(images);
};

const removeImage = (idx: number) => {
  form.value.images = normalizeImages(form.value.images.filter((_, index) => index !== idx));
};

const setMainImage = (idx: number) => {
  form.value.images = form.value.images.map((image, index) => ({
    ...image,
    is_main: index === idx,
  }));
};

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  const formData = new FormData();
  formData.append('image', file);

  isUploading.value = true;
  errorMessage.value = '';
  try {
    const { data } = await api.post('/admin/upload', formData);
    const url = String(data?.url || '').trim();

    if (!url) {
      throw new Error('Upload response does not contain image URL');
    }

    form.value.images = normalizeImages([
      ...form.value.images,
      {
        url,
        is_main: form.value.images.length === 0,
      },
    ]);
    notificationStore.show('Изображение загружено', 'success');
  } catch (err: any) {
    errorMessage.value = getApiErrorMessage(err);
    notificationStore.show(errorMessage.value, 'error');
  } finally {
    isUploading.value = false;
    target.value = '';
  }
};

const triggerUpload = () => {
  fileInput.value?.click();
};

const save = async () => {
  if (isSaving.value) return;

  errorMessage.value = '';
  isSaving.value = true;
  try {
    const payload = buildProjectPayload();

    if (props.product?.id) {
      await api.put(`/admin/projects/${props.product.id}`, payload);
    } else {
      await api.post('/admin/projects', payload);
    }

    notificationStore.show('Проект сохранен', 'success');
    emit('saved');
  } catch (err: any) {
    errorMessage.value = getApiErrorMessage(err);
    notificationStore.show(errorMessage.value, 'error');
  } finally {
    isSaving.value = false;
  }
};

const statusLabels: Record<ProjectStatus, string> = {
  published: 'Опубликован',
  draft: 'Черновик',
  archived: 'Архив',
};

const projectStatuses: ProjectStatus[] = ['published', 'draft', 'archived'];
</script>

<template>
  <div class="fixed inset-0 z-[80] flex items-center justify-center p-2 sm:p-4">
    <div class="absolute inset-0 bg-brand-brown/80 backdrop-blur-sm" @click="emit('close')"></div>
    <div class="relative flex max-h-[calc(100vh-1rem)] w-full max-w-4xl flex-col overflow-hidden rounded-lg bg-white shadow-2xl sm:max-h-[90vh]">
      <header class="flex shrink-0 items-center justify-between gap-4 border-b border-brand-brown/10 p-4 sm:p-6">
        <div>
          <p class="text-xs font-black uppercase text-brand-gold">Каталог</p>
          <h2 class="mt-1 font-serif text-2xl font-bold leading-tight text-brand-brown sm:text-3xl">
            {{ product ? 'Редактировать проект' : 'Добавить проект' }}
          </h2>
        </div>
        <button type="button" @click="emit('close')" class="rounded-lg p-2 transition-colors hover:bg-brand-gray" aria-label="Закрыть">
          <LucideX :size="24" />
        </button>
      </header>

      <div class="space-y-6 overflow-y-auto p-4 sm:p-6">
        <div v-if="errorMessage" class="rounded-lg border border-red-100 bg-red-50 p-4 text-sm font-medium text-red-700">
          {{ errorMessage }}
        </div>

        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <div class="space-y-6">
            <div>
              <label class="mb-2 block text-sm font-semibold text-brand-brown">Название проекта</label>
              <input v-model="form.name" type="text" class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/30 px-4 py-3 outline-none transition-colors focus:border-brand-gold" placeholder="Кухня в скандинавском стиле">
            </div>
            <div>
              <label class="mb-2 block text-sm font-semibold text-brand-brown">Slug для URL</label>
              <input v-model="form.slug" type="text" class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/30 px-4 py-3 outline-none transition-colors focus:border-brand-gold" placeholder="scandi-white-kitchen">
            </div>
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label class="mb-2 block text-sm font-semibold text-brand-brown">Бюджет проекта, ₽</label>
                <input v-model.number="form.price" type="number" min="0" class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/30 px-4 py-3 outline-none transition-colors focus:border-brand-gold">
              </div>
              <div>
                <label class="mb-2 block text-sm font-semibold text-brand-brown">Старая цена, ₽</label>
                <input v-model.number="form.price_old" type="number" min="0" class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/30 px-4 py-3 outline-none transition-colors focus:border-brand-gold">
              </div>
            </div>
            <div>
              <label class="mb-2 block text-sm font-semibold text-brand-brown">Категория</label>
              <select v-model="form.project_category_id" class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/30 px-4 py-3 outline-none transition-colors focus:border-brand-gold">
                <option :value="undefined">Без категории</option>
                <option v-for="category in categories" :key="category.id" :value="category.id">{{ category.name }}</option>
              </select>
            </div>
          </div>

          <div class="space-y-6">
            <div>
              <label class="mb-2 block text-sm font-semibold text-brand-brown">Статус публикации</label>
              <div class="grid grid-cols-3 gap-1 rounded-lg bg-brand-gray/50 p-1">
                <button
                  v-for="status in projectStatuses"
                  :key="status"
                  type="button"
                  @click="form.status = status"
                  :class="[
                    'rounded-lg px-2 py-2 text-[10px] font-bold uppercase transition-all',
                    form.status === status ? 'bg-white text-brand-brown shadow-sm' : 'text-brand-brown/45 hover:text-brand-brown',
                  ]"
                >
                  {{ statusLabels[status] }}
                </button>
              </div>
            </div>
            <div>
              <label class="mb-2 block text-sm font-semibold text-brand-brown">Описание проекта</label>
              <textarea v-model="form.description" rows="4" class="w-full rounded-lg border border-brand-brown/10 bg-brand-gray/30 px-4 py-3 outline-none transition-colors focus:border-brand-gold" placeholder="Коротко опишите задачу, материалы и особенности проекта"></textarea>
            </div>
            <div>
              <label class="mb-2 block text-sm font-semibold text-brand-gold">AI-теги через запятую</label>
              <input v-model="form.ai_tags" type="text" class="w-full rounded-lg border border-brand-gold/20 bg-brand-gold/5 px-4 py-3 outline-none transition-colors focus:border-brand-gold" placeholder="сканди, белая кухня, светлый интерьер">
            </div>
          </div>
        </div>

        <div>
          <div class="mb-4 flex items-center justify-between gap-3">
            <label class="block text-sm font-semibold text-brand-brown">Фотографии проекта</label>
            <input ref="fileInput" type="file" class="hidden" accept="image/*" @change="handleFileUpload">
          </div>
          <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-5">
            <div v-for="(image, idx) in form.images" :key="`${image.url}-${idx}`" class="group relative aspect-square overflow-hidden rounded-lg border border-brand-brown/10 bg-brand-gray">
              <img :src="image.url" class="h-full w-full object-cover" alt="">
              <div class="absolute inset-0 flex items-center justify-center gap-1 bg-black/35 px-1 opacity-100 transition-opacity sm:opacity-0 sm:group-hover:opacity-100">
                <button v-if="idx > 0" type="button" @click="moveImage(idx, 'left')" class="rounded-lg bg-white/20 p-1.5 text-white transition-colors hover:bg-white/40" aria-label="Переместить левее">
                  <LucideArrowLeft :size="14" />
                </button>
                <button type="button" @click="removeImage(idx)" class="rounded-lg bg-red-500 p-1.5 text-white transition-colors hover:bg-red-600" aria-label="Удалить изображение">
                  <LucideTrash2 :size="14" />
                </button>
                <button v-if="idx < form.images.length - 1" type="button" @click="moveImage(idx, 'right')" class="rounded-lg bg-white/20 p-1.5 text-white transition-colors hover:bg-white/40" aria-label="Переместить правее">
                  <LucideArrowRight :size="14" />
                </button>
                <button v-if="!image.is_main" type="button" @click="setMainImage(idx)" class="rounded-lg bg-brand-gold p-1.5 text-white transition-colors hover:bg-brand-gold/80" aria-label="Сделать обложкой">
                  <LucideCheck :size="14" />
                </button>
              </div>
              <div v-if="image.is_main" class="absolute bottom-2 left-2 rounded bg-brand-gold px-2 py-0.5 text-[10px] font-bold uppercase text-white">
                Обложка
              </div>
            </div>

            <button
              type="button"
              @click="triggerUpload"
              :disabled="isUploading"
              class="flex aspect-square flex-col items-center justify-center rounded-lg border-2 border-dashed border-brand-brown/10 text-brand-brown/35 transition-all hover:border-brand-gold hover:text-brand-gold disabled:cursor-not-allowed disabled:opacity-50"
            >
              <LucideUpload v-if="!isUploading" :size="24" class="mb-2" />
              <div v-else class="mb-2 h-6 w-6 animate-spin rounded-full border-2 border-brand-gold border-t-transparent"></div>
              <span class="text-center text-xs font-bold uppercase">{{ isUploading ? 'Загрузка...' : 'Добавить фото' }}</span>
            </button>
          </div>
        </div>

        <div>
          <label class="mb-4 block text-sm font-semibold text-brand-brown">Характеристики проекта</label>
          <div class="mb-4 grid grid-cols-1 gap-3 md:grid-cols-2">
            <div v-for="(value, key) in form.specs" :key="key" class="flex items-center gap-2 rounded-lg border border-brand-brown/10 bg-brand-gray/30 p-3">
              <span class="shrink-0 font-semibold text-brand-brown">{{ key }}:</span>
              <span class="min-w-0 flex-1 truncate text-brand-brown/60">{{ value }}</span>
              <button type="button" @click="removeSpec(key as string)" class="text-brand-brown/25 transition-colors hover:text-red-500" aria-label="Удалить характеристику">
                <LucideTrash2 :size="16" />
              </button>
            </div>
          </div>
          <div class="flex flex-col gap-3 rounded-lg border border-brand-gold/10 bg-brand-gold/5 p-4 sm:flex-row">
            <input v-model="newSpecKey" type="text" placeholder="Название, например материал" class="min-h-10 flex-1 bg-transparent text-sm outline-none">
            <input v-model="newSpecValue" type="text" placeholder="Значение, например МДФ эмаль" class="min-h-10 flex-1 bg-transparent text-sm outline-none">
            <button type="button" @click="addSpec" class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg bg-brand-gold px-4 py-2 text-sm font-bold text-white transition-colors hover:bg-brand-brown">
              <LucidePlus :size="16" />
              Добавить
            </button>
          </div>
        </div>
      </div>

      <footer class="flex shrink-0 flex-col-reverse gap-2 border-t border-brand-brown/10 bg-brand-gray/20 p-4 sm:flex-row sm:justify-end sm:p-6">
        <button type="button" @click="emit('close')" class="rounded-lg px-6 py-3 font-bold text-brand-brown/55 transition-colors hover:text-brand-brown">Отмена</button>
        <button
          type="button"
          @click="save"
          :disabled="isSaving || isUploading"
          class="rounded-lg bg-brand-brown px-6 py-3 font-bold text-white shadow-lg shadow-brand-brown/10 transition-all hover:bg-brand-gold disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ isSaving ? 'Сохраняем...' : 'Сохранить' }}
        </button>
      </footer>
    </div>
  </div>
</template>
