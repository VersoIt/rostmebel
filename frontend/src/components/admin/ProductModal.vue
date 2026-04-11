<script setup lang="ts">
import { ref, watch } from 'vue';
import { LucideX, LucidePlus, LucideTrash2, LucideUpload, LucideArrowLeft, LucideArrowRight } from 'lucide-vue-next';
import type { Product, Category } from '@/types';
import api from '@/api/client';
import { useNotificationStore } from '@/stores/notifications';

const props = defineProps<{
  product?: Product | null;
  categories: Category[];
}>();

const emit = defineEmits(['close', 'saved']);
const notificationStore = useNotificationStore();

const form = ref<Partial<Product>>({
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

const newSpecKey = ref('');
const newSpecValue = ref('');

watch(() => props.product, (val) => {
  if (val) {
    form.value = { ...val };
  } else {
    form.value = {
      name: '',
      slug: '',
      description: '',
      price: 0,
      status: 'draft',
      images: [],
      specs: {},
    };
  }
}, { immediate: true });

const addSpec = () => {
  if (newSpecKey.value && newSpecValue.value) {
    if (!form.value.specs) form.value.specs = {};
    form.value.specs[newSpecKey.value] = newSpecValue.value;
    newSpecKey.value = '';
    newSpecValue.value = '';
  }
};

const removeSpec = (key: string) => {
  if (form.value.specs) {
    delete form.value.specs[key];
  }
};

const moveImage = (idx: number, direction: 'left' | 'right') => {
  if (!form.value.images) return;
  const newIdx = direction === 'left' ? idx - 1 : idx + 1;
  if (newIdx < 0 || newIdx >= form.value.images.length) return;
  
  const images = [...form.value.images];
  [images[idx], images[newIdx]] = [images[newIdx], images[idx]];
  form.value.images = images;
};

const fileInput = ref<HTMLInputElement | null>(null);
const isUploading = ref(false);

const handleFileUpload = async (e: Event) => {
  const target = e.target as HTMLInputElement;
  if (!target.files?.length) return;

  const file = target.files[0];
  const formData = new FormData();
  formData.append('image', file);

  isUploading.value = true;
  errorMessage.value = '';
  try {
    const { data } = await api.post('/admin/upload', formData);
    
    if (!form.value.images) form.value.images = [];
    form.value.images.push({ 
      url: data.url, 
      is_main: form.value.images.length === 0 
    });
    notificationStore.show('Изображение загружено', 'success');
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Ошибка при загрузке изображения';
    notificationStore.show(errorMessage.value, 'error');
  } finally {
    isUploading.value = false;
    target.value = '';
  }
};

const triggerUpload = () => {
  fileInput.value?.click();
};

const errorMessage = ref('');

const save = async () => {
  errorMessage.value = '';
  try {
    const payload = { ...form.value };
    if (!payload.price_old) payload.price_old = undefined;

    if (props.product?.id) {
      await api.put(`/admin/projects/${props.product.id}`, payload);
    } else {
      await api.post('/admin/projects', payload);
    }
    emit('saved');
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Ошибка при сохранении';
  }
};
</script>

<template>
  <div class="fixed inset-0 z-[60] flex items-center justify-center p-4">
    <div class="absolute inset-0 bg-brand-brown/80 backdrop-blur-sm" @click="$emit('close')"></div>
    <div class="relative bg-white w-full max-w-4xl rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
      <header class="p-8 border-b border-brand-brown/5 flex items-center justify-between shrink-0">
        <h2 class="font-serif text-3xl">{{ product ? 'Редактировать' : 'Добавить' }} проект</h2>
        <button @click="$emit('close')" class="p-2 hover:bg-brand-gray rounded-full transition-colors">
          <LucideX :size="24" />
        </button>
      </header>

      <div class="p-8 overflow-y-auto space-y-8">
        <div v-if="errorMessage" class="bg-red-50 text-red-600 p-4 rounded-2xl border border-red-100 text-sm font-medium">
          {{ errorMessage }}
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          <div class="space-y-6">
            <div>
              <label class="block text-sm font-semibold mb-2">Название проекта</label>
              <input v-model="form.name" type="text" class="w-full px-4 py-3 rounded-xl border border-brand-brown/10 bg-brand-gray/30" placeholder="Кухня в стиле лофт">
            </div>
            <div>
              <label class="block text-sm font-semibold mb-2">Slug (URL)</label>
              <input v-model="form.slug" type="text" class="w-full px-4 py-3 rounded-xl border border-brand-brown/10 bg-brand-gray/30" placeholder="kitchen-loft">
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-semibold mb-2">Бюджет проекта (₽)</label>
                <input v-model.number="form.price" type="number" class="w-full px-4 py-3 rounded-xl border border-brand-brown/10 bg-brand-gray/30">
              </div>
              <div>
                <label class="block text-sm font-semibold mb-2">Ориент. бюджет</label>
                <input v-model.number="form.price_old" type="number" class="w-full px-4 py-3 rounded-xl border border-brand-brown/10 bg-brand-gray/30">
              </div>
            </div>
            <div>
              <label class="block text-sm font-semibold mb-2">Категория</label>
              <select v-model="form.project_category_id" class="w-full px-4 py-3 rounded-xl border border-brand-brown/10 bg-brand-gray/30">
                <option :value="undefined">Без категории</option>
                <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
              </select>
            </div>
          </div>

          <div class="space-y-6">
            <div>
              <label class="block text-sm font-semibold mb-2">Статус публикации</label>
              <div class="flex gap-2 p-1 bg-brand-gray/50 rounded-xl">
                <button 
                  v-for="s in ['published', 'draft', 'archived']" 
                  :key="s"
                  @click="form.status = s as any"
                  :class="['flex-1 py-2 rounded-lg text-[10px] font-bold uppercase tracking-widest transition-all', form.status === s ? 'bg-white text-brand-brown shadow-sm' : 'text-brand-brown/40']"
                >
                  {{ s === 'published' ? 'Опубликован' : s === 'draft' ? 'Черновик' : 'Архив' }}
                </button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-semibold mb-2">Описание проекта</label>
              <textarea v-model="form.description" rows="4" class="w-full px-4 py-3 rounded-xl border border-brand-brown/10 bg-brand-gray/30"></textarea>
            </div>
            <div>
              <label class="block text-sm font-semibold mb-2 text-brand-gold">AI Теги (через запятую)</label>
              <input v-model="form.ai_tags" type="text" class="w-full px-4 py-3 rounded-xl border border-brand-gold/20 bg-brand-gold/5 outline-none focus:border-brand-gold" placeholder="лофт, дерево, коричневый, кухня">
            </div>
          </div>
        </div>

        <div>
          <div class="flex items-center justify-between mb-4">
            <label class="block text-sm font-semibold">Фотографии проекта</label>
            <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="handleFileUpload">
          </div>
          <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
            <div v-for="(img, idx) in form.images" :key="idx" class="relative aspect-square rounded-2xl overflow-hidden border border-brand-brown/10 group">
              <img :src="img.url" class="w-full h-full object-cover">
              <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-1 px-1">
                <button v-if="idx > 0" @click="moveImage(idx, 'left')" class="p-1.5 bg-white/20 text-white rounded-lg hover:bg-white/40">
                  <LucideArrowLeft :size="14" />
                </button>
                <button @click="form.images?.splice(idx, 1)" class="p-1.5 bg-red-500 text-white rounded-lg hover:bg-red-600">
                  <LucideTrash2 :size="14" />
                </button>
                <button v-if="idx < (form.images?.length || 0) - 1" @click="moveImage(idx, 'right')" class="p-1.5 bg-white/20 text-white rounded-lg hover:bg-white/40">
                  <LucideArrowRight :size="14" />
                </button>
                <button v-if="!img.is_main" @click="form.images?.forEach((m, i) => m.is_main = i === idx)" class="p-1.5 bg-brand-gold text-white rounded-lg hover:bg-brand-gold/80">
                  <LucidePlus :size="14" />
                </button>
              </div>
              <div v-if="img.is_main" class="absolute bottom-2 left-2 bg-brand-gold text-white text-[10px] px-2 py-0.5 rounded font-bold">ОБЛОЖКА</div>
            </div>
            
            <button @click="triggerUpload" :disabled="isUploading" class="aspect-square rounded-2xl border-2 border-dashed border-brand-brown/10 flex flex-col items-center justify-center text-brand-brown/20 hover:border-brand-gold hover:text-brand-gold transition-all disabled:opacity-50">
              <LucideUpload v-if="!isUploading" :size="24" class="mb-2" />
              <div v-else class="w-6 h-6 border-2 border-brand-gold border-t-transparent animate-spin rounded-full mb-2"></div>
              <span class="text-xs font-bold uppercase tracking-widest">{{ isUploading ? 'Загрузка...' : 'Добавить фото' }}</span>
            </button>
          </div>
        </div>

        <div>
          <label class="block text-sm font-semibold mb-4">Детали проекта</label>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div v-for="(value, key) in form.specs" :key="key" class="flex items-center gap-2 bg-brand-gray/30 p-3 rounded-xl border border-brand-brown/5">
              <span class="font-semibold shrink-0">{{ key }}:</span>
              <span class="flex-1 text-brand-brown/60">{{ value }}</span>
              <button @click="removeSpec(key as string)" class="text-brand-brown/20 hover:text-red-500">
                <LucideTrash2 :size="16" />
              </button>
            </div>
          </div>
          <div class="flex gap-4 p-4 bg-brand-gold/5 rounded-2xl border border-brand-gold/10">
            <input v-model="newSpecKey" type="text" placeholder="Название (Материал)" class="flex-1 bg-transparent border-none outline-none text-sm">
            <input v-model="newSpecValue" type="text" placeholder="Значение (Дуб массив)" class="flex-1 bg-transparent border-none outline-none text-sm">
            <button @click="addSpec" class="bg-brand-gold text-white px-4 py-2 rounded-lg text-sm font-bold">ОК</button>
          </div>
        </div>
      </div>

      <footer class="p-8 border-t border-brand-brown/5 bg-brand-gray/10 flex justify-end gap-4 shrink-0">
        <button @click="$emit('close')" class="px-8 py-3 rounded-xl font-bold text-brand-brown/40 hover:text-brand-brown">Отмена</button>
        <button @click="save" class="px-8 py-3 bg-brand-brown text-white rounded-xl font-bold hover:bg-brand-gold shadow-lg transition-all">Сохранить</button>
      </footer>
    </div>
  </div>
</template>
