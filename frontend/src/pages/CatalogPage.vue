<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useProductStore } from '@/stores/products';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { LucideArrowRight, LucideFilterX, LucideLoader2, LucideSearch } from 'lucide-vue-next';
import { absoluteUrl, removeJsonLd, setJsonLd } from '@/utils/seo';

const productStore = useProductStore();
const route = useRoute();
const router = useRouter();

const selectedCategory = ref(route.query.category?.toString() || '');
const searchQuery = ref('');

const updateCatalogSchema = () => {
  setJsonLd('schema-catalog', {
    '@context': 'https://schema.org',
    '@type': 'CollectionPage',
    name: 'Проекты кухонь и корпусной мебели',
    url: absoluteUrl('/catalog'),
    mainEntity: {
      '@type': 'ItemList',
      itemListElement: productStore.products.map((product, index) => ({
        '@type': 'ListItem',
        position: index + 1,
        url: absoluteUrl(`/product/${product.slug || product.id}`),
        name: product.name,
      })),
    },
  });
};

const fetch = async () => {
  const params: Record<string, string | number> = { status: 'published' };

  if (selectedCategory.value) {
    const category = productStore.categories.find((item) => item.slug === selectedCategory.value);
    if (category) params.project_category_id = category.id;
  }

  if (searchQuery.value.trim()) {
    params.search = searchQuery.value.trim();
  }

  await productStore.fetchProducts(params);
  updateCatalogSchema();
};

onMounted(async () => {
  await productStore.fetchCategories();
  await fetch();
});

onUnmounted(() => {
  removeJsonLd('schema-catalog');
});

const selectCategory = async (slug: string) => {
  selectedCategory.value = slug;
  await router.push({ query: { ...route.query, category: slug || undefined } });
  fetch();
};

const resetFilters = async () => {
  selectedCategory.value = '';
  searchQuery.value = '';
  await router.push({ query: {} });
  fetch();
};
</script>

<template>
  <div class="min-h-screen bg-brand-cream pt-28">
    <div class="ui-container ui-section">
      <header class="mb-10">
        <p class="ui-eyebrow mb-3">Портфолио</p>
        <h1 class="ui-title-xl">Проекты</h1>
        <p class="ui-copy-lg mt-4 max-w-2xl">
          Реальные кухни, шкафы и системы хранения. Фильтруйте по категории или найдите проект по названию.
        </p>
      </header>

      <section class="mb-8 space-y-4">
        <div class="-mx-4 overflow-x-auto px-4 no-scrollbar sm:mx-0 sm:px-0">
          <div class="flex min-w-max gap-2">
            <button
              type="button"
              :class="['ui-chip', !selectedCategory ? 'ui-chip-active' : '']"
              @click="selectCategory('')"
            >
              Все работы
            </button>
            <button
              v-for="category in productStore.categories"
              :key="category.id"
              type="button"
              :class="['ui-chip', selectedCategory === category.slug ? 'ui-chip-active' : '']"
              @click="selectCategory(category.slug)"
            >
              {{ category.name }}
            </button>
          </div>
        </div>

        <div class="ui-card p-4">
          <div class="relative max-w-xl">
            <input
              v-model="searchQuery"
              type="text"
              class="ui-input pl-11"
              placeholder="Название, материал, стиль"
              @input="fetch"
            >
            <LucideSearch class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/35" :size="19" />
          </div>
        </div>
      </section>

      <main>
        <div v-if="productStore.loading" class="ui-empty py-20">
          <LucideLoader2 class="mx-auto mb-4 animate-spin text-brand-gold" :size="40" />
          Загружаем проекты
        </div>

        <div v-else-if="productStore.products.length > 0" class="motion-fade-up">
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
            <ProductCard v-for="product in productStore.products" :key="product.id" :product="product" />
          </div>

          <div class="mt-10 flex items-center justify-between border-t border-brand-brown/10 pt-5">
            <span class="text-sm font-medium text-brand-brown/45">Показано {{ productStore.products.length }} проектов</span>
            <router-link to="/contact" class="hidden items-center gap-2 text-sm font-bold text-brand-gold hover:text-brand-brown sm:inline-flex">
              Обсудить похожий проект
              <LucideArrowRight :size="17" />
            </router-link>
          </div>
        </div>

        <div v-else class="ui-empty py-16">
          <LucideFilterX :size="56" class="mx-auto mb-5 text-brand-brown/12" />
          <h2 class="ui-title-md mb-2">Проекты не найдены</h2>
          <p class="mx-auto mb-6 max-w-md text-brand-brown/55">
            Попробуйте другой запрос или сбросьте фильтр категории.
          </p>
          <button type="button" class="ui-button ui-button-primary" @click="resetFilters">
            Сбросить фильтры
          </button>
        </div>
      </main>
    </div>
  </div>
</template>
