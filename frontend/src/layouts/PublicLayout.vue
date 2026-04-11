<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { LucideMenu, LucideX, LucideUser, LucideHeart } from 'lucide-vue-next';
import { useFavorites } from '@/composables/useFavorites';

const { favorites } = useFavorites();
const isMenuOpen = ref(false);
const isScrolled = ref(false);

const handleScroll = () => {
  isScrolled.value = window.scrollY > 50;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<template>
  <div class="font-sans text-brand-brown">
    <nav 
      :class="[
        'fixed top-0 left-0 right-0 z-50 transition-all duration-300 px-6 py-4',
        isScrolled || isMenuOpen ? 'bg-white/90 backdrop-blur-md shadow-sm py-3' : 'bg-transparent'
      ]"
    >
      <div class="max-w-7xl mx-auto flex items-center justify-between">
        <router-link to="/" class="flex items-center gap-2">
          <div class="w-10 h-10 bg-brand-gold rounded-xl flex items-center justify-center text-white font-serif text-2xl">Р</div>
          <span class="font-serif text-2xl tracking-tighter uppercase font-bold text-brand-brown">РОСТ <span class="text-brand-gold">Мебель</span></span>
        </router-link>

        <div class="hidden md:flex items-center gap-10">
          <router-link to="/" class="font-medium hover:text-brand-gold transition-colors">Главная</router-link>
          <router-link to="/catalog" class="font-medium hover:text-brand-gold transition-colors">Проекты</router-link>
          <router-link to="/contact" class="font-medium hover:text-brand-gold transition-colors">Контакты</router-link>
        </div>

        <div class="flex items-center gap-4">
          <router-link to="/favorites" class="relative p-2 text-brand-brown/60 hover:text-brand-gold transition-colors">
            <LucideHeart :size="24" />
            <span v-if="favorites.length > 0" class="absolute top-0 right-0 bg-brand-gold text-white text-[10px] font-bold w-4 h-4 rounded-full flex items-center justify-center">
              {{ favorites.length }}
            </span>
          </router-link>
          <button @click="isMenuOpen = !isMenuOpen" class="md:hidden p-2 text-brand-brown">
            <component :is="isMenuOpen ? LucideX : LucideMenu" :size="28" />
          </button>
        </div>
      </div>

      <transition 
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="opacity-0 -translate-y-4"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-4"
      >
        <div v-if="isMenuOpen" class="md:hidden absolute top-full left-0 right-0 bg-white border-t border-brand-brown/5 p-6 shadow-xl">
          <div class="flex flex-col gap-6 text-center">
            <router-link to="/" @click="isMenuOpen = false" class="text-xl font-medium">Главная</router-link>
            <router-link to="/catalog" @click="isMenuOpen = false" class="text-xl font-medium">Проекты</router-link>
            <router-link to="/contact" @click="isMenuOpen = false" class="text-xl font-medium">Контакты</router-link>
          </div>
        </div>
      </transition>
    </nav>

    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>

    <footer class="bg-brand-gray/50 border-t border-brand-brown/5 py-16 px-6">
      <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-12">
        <div class="md:col-span-2">
          <router-link to="/" class="flex items-center gap-2 mb-6">
            <div class="w-8 h-8 bg-brand-gold rounded-lg flex items-center justify-center text-white font-serif text-xl">Р</div>
            <span class="font-serif text-xl uppercase font-bold">РОСТ Мебель</span>
          </router-link>
          <p class="text-brand-brown/60 max-w-sm mb-6 leading-relaxed">
            Производство премиальной мебели по индивидуальным заказам.
          </p>
        </div>
        <div>
          <h4 class="font-serif text-lg mb-6 uppercase tracking-widest text-brand-brown/40">Разделы</h4>
          <ul class="space-y-4">
            <li><router-link to="/" class="hover:text-brand-gold transition-colors">Главная</router-link></li>
            <li><router-link to="/catalog" class="hover:text-brand-gold transition-colors">Проекты</router-link></li>
            <li><router-link to="/contact" class="hover:text-brand-gold transition-colors">Контакты</router-link></li>
          </ul>
        </div>
        <div>
          <h4 class="font-serif text-lg mb-6 uppercase tracking-widest text-brand-brown/40">Контакты</h4>
          <ul class="space-y-4 text-brand-brown/80">
            <li>+7 (978) 763-16-03</li>
            <li>rost.salon2003@mail.ru</li>
          </ul>
        </div>
      </div>
    </footer>
  </div>
</template>
