<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { LucideArrowRight, LucideFilterX, LucideLoader2, LucideSearch } from 'lucide-vue-next';
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
  <div class="min-h-screen bg-brand-cream pt-28">
    <div class="ui-container ui-section">
      <header class="mb-10">
        <p class="ui-eyebrow mb-3">Портфолио</p>
        <h1 class="ui-title-xl">Проекты</h1>
        <p class="ui-copy-lg mt-4 max-w-2xl">
          Реальные кухни, шкафы, столы, прихожие и другая мебель по размеру. Фильтруйте по категории или найдите проект по названию, материалу и стилю.
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
              @input="queueSearch"
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
