<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue';
import { useRoute } from 'vue-router';
import {
  LucideCalculator,
  LucideHeart,
  LucideMenu,
  LucideMessageCircle,
  LucidePhone,
  LucideX,
} from 'lucide-vue-next';
import { useFavorites } from '@/composables/useFavorites';
import BrandMark from '@/components/common/BrandMark.vue';

const { favorites } = useFavorites();
const route = useRoute();
const isMenuOpen = ref(false);
const isScrolled = ref(false);

const navLinks = [
  { name: 'Проекты', to: '/catalog' },
  { name: 'Контакты', to: '/contact' },
];

const mobileLinks = [
  { name: 'Главная', to: '/' },
  ...navLinks,
  { name: 'Избранное', to: '/favorites' },
];

const quickActions = [
  { name: 'Позвонить', href: 'tel:+79787631603', icon: LucidePhone },
  { name: 'WhatsApp', href: 'https://wa.me/79787631603', icon: LucideMessageCircle },
  { name: 'MAX', href: 'https://max.ru', icon: LucideMessageCircle },
];

const isHomePage = computed(() => route.name === 'home');
const isHeaderActive = computed(() => !isHomePage.value || isScrolled.value || isMenuOpen.value);

const handleScroll = () => {
  isScrolled.value = window.scrollY > 20;
};

onMounted(() => {
  handleScroll();
  window.addEventListener('scroll', handleScroll, { passive: true });
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});

watch(isMenuOpen, (open) => {
  document.body.style.overflow = open ? 'hidden' : '';
});

onUnmounted(() => {
  document.body.style.overflow = '';
});
</script>

<template>
  <div class="flex min-h-screen flex-col bg-brand-cream font-sans text-brand-brown selection:bg-brand-gold selection:text-white">
    <header
      :class="[
        'fixed inset-x-0 top-0 z-[100] border-b transition-all duration-300',
        isHeaderActive
          ? 'border-brand-brown/10 bg-white/92 py-3 shadow-sm backdrop-blur-xl'
          : 'border-transparent bg-transparent py-5'
      ]"
    >
      <div class="ui-container flex items-center justify-between gap-4">
        <router-link to="/" class="flex min-w-0 items-center gap-2" aria-label="РОСТ Мебель">
          <BrandMark
            :class="[
              'shrink-0 transition-colors duration-200',
              isHeaderActive ? 'text-brand-gold' : 'text-white'
            ]"
            size="3rem"
          />
          <div class="min-w-0">
            <div :class="['truncate font-serif text-2xl font-black uppercase leading-none', isHeaderActive ? 'text-brand-brown' : 'text-white']">
              РОСТ
            </div>
            <div class="mt-1 text-[10px] font-black uppercase tracking-widest text-brand-gold">Мебель</div>
          </div>
        </router-link>

        <nav class="hidden items-center gap-8 md:flex">
          <router-link
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            :class="[
              'text-xs font-black uppercase tracking-widest transition-colors hover:text-brand-gold',
              isHeaderActive ? 'text-brand-brown' : 'text-white'
            ]"
          >
            {{ link.name }}
          </router-link>
        </nav>

        <div class="flex items-center gap-2 sm:gap-3">
          <a
            href="tel:+79787631603"
            :class="[
              'hidden h-11 items-center gap-2 rounded-lg border px-4 text-xs font-black uppercase tracking-widest transition-colors lg:flex',
              isHeaderActive
                ? 'border-brand-brown/10 text-brand-brown hover:bg-brand-brown hover:text-white'
                : 'border-white/25 text-white hover:bg-white hover:text-brand-brown'
            ]"
          >
            <LucidePhone :size="15" />
            +7 (978) 763-16-03
          </a>

          <router-link
            to="/favorites"
            class="relative flex h-11 w-11 items-center justify-center rounded-lg transition-colors hover:bg-brand-gold/12"
            aria-label="Избранное"
          >
            <LucideHeart :class="[isHeaderActive ? 'text-brand-brown' : 'text-white']" :size="23" />
            <span
              v-if="favorites.length"
              class="absolute -right-1 -top-1 flex h-5 min-w-5 items-center justify-center rounded-lg bg-brand-gold px-1 text-[10px] font-black text-white shadow-md"
            >
              {{ favorites.length }}
            </span>
          </router-link>

          <button
            type="button"
            class="flex h-11 w-11 items-center justify-center rounded-lg transition-colors hover:bg-brand-gold/12 md:hidden"
            aria-label="Открыть меню"
            @click="isMenuOpen = true"
          >
            <LucideMenu :class="isHeaderActive ? 'text-brand-brown' : 'text-white'" :size="28" />
          </button>
        </div>
      </div>
    </header>

    <Teleport to="body">
      <transition name="menu">
        <div v-if="isMenuOpen" class="fixed inset-0 z-[210] flex flex-col bg-brand-brown p-5 text-white">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <BrandMark class="shrink-0 text-brand-gold" size="3rem" />
              <div>
                <div class="text-[11px] font-black uppercase tracking-widest text-brand-gold">РОСТ Мебель</div>
                <div class="font-serif text-3xl font-bold">Меню</div>
              </div>
            </div>
            <button
              type="button"
              class="flex h-12 w-12 items-center justify-center rounded-lg bg-white/10 text-white transition-colors hover:bg-white hover:text-brand-brown"
              aria-label="Закрыть меню"
              @click="isMenuOpen = false"
            >
              <LucideX :size="28" />
            </button>
          </div>

          <nav class="mt-14 flex flex-col gap-4">
            <router-link
              v-for="link in mobileLinks"
              :key="link.to"
              :to="link.to"
              class="rounded-lg border border-white/10 px-4 py-4 font-serif text-3xl transition-colors hover:border-brand-gold hover:text-brand-gold"
              @click="isMenuOpen = false"
            >
              {{ link.name }}
            </router-link>
          </nav>

          <div class="mt-auto border-t border-white/10 pt-6">
            <p class="mb-3 text-xs font-black uppercase tracking-widest text-white/45">Связаться</p>
            <a href="tel:+79787631603" class="font-serif text-3xl font-bold text-white">+7 (978) 763-16-03</a>
          </div>
        </div>
      </transition>
    </Teleport>

    <main class="flex-1">
      <router-view />
    </main>

    <nav
      class="fixed inset-x-3 bottom-3 z-[180] grid grid-cols-4 gap-1 rounded-lg border border-brand-brown/10 bg-white/95 p-1.5 shadow-2xl shadow-black/20 backdrop-blur-xl md:hidden"
      aria-label="Быстрые действия"
    >
      <a
        v-for="action in quickActions"
        :key="action.name"
        :href="action.href"
        class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg text-[11px] font-black text-brand-brown/70 transition-colors hover:bg-brand-gray hover:text-brand-gold"
        :target="action.href.startsWith('http') ? '_blank' : undefined"
        :rel="action.href.startsWith('http') ? 'noopener noreferrer' : undefined"
      >
        <component :is="action.icon" :size="18" />
        {{ action.name }}
      </a>
      <router-link
        :to="{ path: '/contact', hash: '#quote-quiz' }"
        class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-lg bg-brand-brown text-[11px] font-black text-white transition-colors hover:bg-brand-gold"
      >
        <LucideCalculator :size="18" />
        Расчет
      </router-link>
    </nav>

    <footer class="border-t border-brand-brown/10 bg-brand-cream px-4 pb-28 pt-12 sm:px-6 md:pb-12 lg:py-16">
      <div class="mx-auto grid max-w-7xl grid-cols-1 gap-10 md:grid-cols-4">
        <div class="md:col-span-2">
          <div class="mb-6 flex items-center gap-2">
            <BrandMark class="shrink-0 text-brand-gold" size="2.35rem" />
            <span class="font-serif text-2xl font-black uppercase text-brand-brown">РОСТ <span class="text-brand-gold">Мебель</span></span>
          </div>
          <p class="max-w-md leading-7 text-brand-brown/62">
            Проектируем, производим и устанавливаем кухни, шкафы и системы хранения по Крыму. Работаем с понятной сметой, договором и монтажом под ключ.
          </p>
        </div>

        <div>
          <h4 class="mb-5 text-xs font-black uppercase tracking-widest text-brand-gold">Навигация</h4>
          <ul class="space-y-3 text-sm font-bold">
            <li><router-link to="/catalog" class="text-brand-brown/70 transition-colors hover:text-brand-gold">Проекты</router-link></li>
            <li><router-link to="/contact" class="text-brand-brown/70 transition-colors hover:text-brand-gold">Контакты</router-link></li>
            <li><router-link to="/favorites" class="text-brand-brown/70 transition-colors hover:text-brand-gold">Избранное</router-link></li>
          </ul>
        </div>

        <div>
          <h4 class="mb-5 text-xs font-black uppercase tracking-widest text-brand-gold">Контакты</h4>
          <ul class="space-y-3 text-sm font-bold text-brand-brown/70">
            <li><a href="tel:+79787631603" class="transition-colors hover:text-brand-gold">+7 (978) 763-16-03</a></li>
            <li><a href="mailto:rost.salon2003@mail.ru" class="break-all transition-colors hover:text-brand-gold">rost.salon2003@mail.ru</a></li>
            <li>Симферополь, выезд по Крыму</li>
          </ul>
        </div>
      </div>

      <div class="mx-auto mt-10 flex max-w-7xl flex-col gap-3 border-t border-brand-brown/10 pt-6 text-xs font-semibold text-brand-brown/35 sm:flex-row sm:items-center sm:justify-between">
        <p>© 2026 РОСТ Мебель. Все права защищены.</p>
        <a href="#" class="transition-colors hover:text-brand-gold">Политика конфиденциальности</a>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.menu-enter-active,
.menu-leave-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.menu-enter-from,
.menu-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
