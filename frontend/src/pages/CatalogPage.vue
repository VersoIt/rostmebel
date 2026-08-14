<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  LucideArrowRight,
  LucideFilterX,
  LucideLoader2,
  LucideRotateCcw,
  LucideSearch,
  LucideSparkles,
} from 'lucide-vue-next';
import ProductCard from '@/components/catalog/ProductCard.vue';
import { useProductStore } from '@/stores/products';
import { absoluteUrl, compactDescription, removeJsonLd, setJsonLd, setPageSeo } from '@/utils/seo';

const productStore = useProductStore();
const route = useRoute();
const router = useRouter();

const selectedCategory = ref('');
const searchQuery = ref('');
let searchTimer: ReturnType<typeof setTimeout> | null = null;

const selectedCategoryEntity = computed(() =>
  productStore.categories.find((item) => item.slug === selectedCategory.value) || null,
);

const trimmedSearchQuery = computed(() => searchQuery.value.trim());

const catalogPath = computed(() => {
  if (trimmedSearchQuery.value) {
    return '/catalog';
  }

  if (selectedCategory.value) {
    const params = new URLSearchParams({ category: selectedCategory.value });
    return `/catalog?${params.toString()}`;
  }

  return '/catalog';
});

const catalogTitle = computed(() => {
  if (trimmedSearchQuery.value) {
    return `Поиск проектов по запросу «${trimmedSearchQuery.value}» — РОСТ Мебель`;
  }

  if (selectedCategoryEntity.value) {
    return `${selectedCategoryEntity.value.name} на заказ — проекты РОСТ Мебель`;
  }

  return 'Проекты кухонь и корпусной мебели — РОСТ Мебель';
});

const catalogDescription = computed(() => {
  if (trimmedSearchQuery.value) {
    return compactDescription(
      `Результаты поиска по запросу «${trimmedSearchQuery.value}» в портфолио РОСТ Мебель. Подберите подходящие проекты мебели по материалам, стилю и бюджету.`,
      155,
    );
  }

  if (selectedCategoryEntity.value) {
    return compactDescription(
      `Подборка категории «${selectedCategoryEntity.value.name}» от РОСТ Мебель: реальные проекты, фото, бюджеты, материалы и детали исполнения.`,
      155,
    );
  }

  return compactDescription(
    'Портфолио реализованных кухонь, шкафов и систем хранения: фотографии, бюджеты, материалы и детали проектов РОСТ Мебель.',
    155,
  );
});

const catalogRobots = computed(() => (trimmedSearchQuery.value ? 'noindex,follow' : 'index,follow'));

const resultLabel = computed(() => {
  if (trimmedSearchQuery.value) {
    return `Найдено ${productStore.products.length} по запросу «${trimmedSearchQuery.value}»`;
  }

  if (selectedCategoryEntity.value) {
    return `${productStore.products.length} в категории «${selectedCategoryEntity.value.name}»`;
  }

  return `${productStore.products.length} проектов`;
});

const hasActiveFilters = computed(() => Boolean(selectedCategory.value || trimmedSearchQuery.value));

const syncFiltersFromRoute = () => {
  selectedCategory.value = route.query.category?.toString() || '';
  searchQuery.value = route.query.search?.toString() || '';
};

const updateCatalogSeo = () => {
  setPageSeo({
    title: catalogTitle.value,
    description: catalogDescription.value,
    path: catalogPath.value,
    robots: catalogRobots.value,
  });
};

const updateCatalogSchema = () => {
  setJsonLd('schema-catalog', {
    '@context': 'https://schema.org',
    '@type': 'CollectionPage',
    name: catalogTitle.value,
    description: catalogDescription.value,
    url: absoluteUrl(catalogPath.value),
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

const fetchCatalog = async () => {
  const params: Record<string, string | number> = { status: 'published' };

  if (selectedCategoryEntity.value) {
    params.project_category_id = selectedCategoryEntity.value.id;
  }

  if (trimmedSearchQuery.value) {
    params.search = trimmedSearchQuery.value;
  }

  await productStore.fetchProducts(params);
  updateCatalogSeo();
  updateCatalogSchema();
};

const applyRouteFilters = async (replace = true) => {
  const query: Record<string, string> = {};
  if (selectedCategory.value) {
    query.category = selectedCategory.value;
  }
  if (trimmedSearchQuery.value) {
    query.search = trimmedSearchQuery.value;
  }

  await router[replace ? 'replace' : 'push']({ query });
};

const queueSearch = () => {
  if (searchTimer) {
    clearTimeout(searchTimer);
  }

  searchTimer = setTimeout(() => {
    void applyRouteFilters(true);
  }, 250);
};

const selectCategory = async (slug: string) => {
  selectedCategory.value = slug;
  await applyRouteFilters(true);
};

const resetFilters = async () => {
  selectedCategory.value = '';
  searchQuery.value = '';
  await router.push({ query: {} });
};

watch(
  () => [route.query.category, route.query.search],
  async () => {
    syncFiltersFromRoute();
    if (!productStore.categories.length) return;
    await fetchCatalog();
  },
);

onMounted(async () => {
  syncFiltersFromRoute();
  await productStore.fetchCategories();
  await fetchCatalog();
});

onUnmounted(() => {
  if (searchTimer) {
    clearTimeout(searchTimer);
  }
  removeJsonLd('schema-catalog');
});
</script>

<template>
  <div class="min-h-screen bg-transparent pt-28">
    <div class="ui-container py-8 sm:py-10 lg:py-12">
      <section class="ui-surface mb-6 p-4 sm:p-5 lg:p-6">
        <div class="grid gap-5 lg:grid-cols-[minmax(0,1fr)_360px] lg:items-start">
          <div>
            <div class="mb-3 inline-flex items-center gap-2 rounded-full bg-brand-gold/10 px-3 py-1 text-[11px] font-black uppercase tracking-[0.18em] text-brand-gold">
              <LucideSparkles :size="14" />
              Портфолио
            </div>
            <h1 class="font-serif text-3xl font-bold leading-tight text-brand-brown sm:text-4xl">Проекты</h1>
            <p class="mt-3 max-w-2xl text-sm leading-7 text-brand-brown/62 sm:text-base">
              Реальные кухни, шкафы и другие проекты с фотографиями, бюджетом и деталями исполнения.
            </p>
          </div>

          <div class="space-y-3">
            <div class="relative">
              <input
                v-model="searchQuery"
                type="text"
                class="ui-input pl-11 pr-11"
                placeholder="Название, материал, стиль, город"
                @input="queueSearch"
              >
              <LucideSearch class="absolute left-4 top-1/2 -translate-y-1/2 text-brand-brown/35" :size="19" />
            </div>

            <div class="flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-brand-brown/8 bg-white/72 px-4 py-3">
              <span class="text-sm font-semibold text-brand-brown/62">{{ resultLabel }}</span>
              <button
                v-if="hasActiveFilters"
                type="button"
                class="inline-flex items-center gap-2 rounded-full bg-brand-gray px-3 py-1.5 text-xs font-bold text-brand-brown transition-colors hover:bg-brand-brown hover:text-white"
                @click="resetFilters"
              >
                <LucideRotateCcw :size="14" />
                Сбросить
              </button>
            </div>
          </div>
        </div>

        <div class="mt-4 border-t border-brand-brown/8 pt-4">
          <div class="-mx-1 overflow-x-auto px-1 no-scrollbar">
            <div class="flex min-w-max gap-2 pb-1">
              <button
                type="button"
                :class="['ui-chip', !selectedCategory ? 'ui-chip-active' : '']"
                @click="selectCategory('')"
              >
                Все проекты
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

          <div class="mt-8 flex items-center justify-between gap-4 border-t border-brand-brown/10 pt-5">
            <span class="text-sm font-medium text-brand-brown/45">{{ resultLabel }}</span>
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
            Попробуйте другой запрос или сбросьте фильтры.
          </p>
          <button type="button" class="ui-button ui-button-primary" @click="resetFilters">
            Сбросить фильтры
          </button>
        </div>
      </main>
    </div>
  </div>
</template>
