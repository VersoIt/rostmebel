<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { LucideMenu, LucideX, LucideUser, LucideHeart, LucidePhone } from 'lucide-vue-next';
import { useFavorites } from '@/composables/useFavorites';

const { favorites } = useFavorites();
const route = useRoute();
const isMenuOpen = ref(false);
const isScrolled = ref(false);

const isHomePage = computed(() => route.name === 'home');

const handleScroll = () => {
  isScrolled.value = window.scrollY > 20;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});

// Logic for header appearance
const isHeaderActive = computed(() => !isHomePage.value || isScrolled.value);
</script>

<template>
  <div class="min-h-screen flex flex-col font-sans selection:bg-brand-gold selection:text-white">
    <!-- Navigation -->
    <header 
      :class="[
        'fixed w-full z-[100] transition-all duration-500 border-b',
        isHeaderActive 
          ? 'py-4 bg-white/80 backdrop-blur-xl border-brand-brown/5 shadow-sm' 
          : 'py-8 bg-transparent border-transparent'
      ]"
    >
      <div class="max-w-7xl mx-auto px-6 flex items-center justify-between">
        <router-link to="/" class="flex items-center gap-3 group">
          <div class="w-12 h-12 bg-brand-brown text-brand-gold rounded-lg flex items-center justify-center font-serif text-2xl group-hover:bg-brand-gold group-hover:text-white transition-all duration-500 shadow-lg">Р</div>
          <div class="flex flex-col">
            <span :class="['font-serif text-2xl uppercase font-black leading-none transition-colors', isHeaderActive ? 'text-brand-brown' : 'text-white']">РОСТ</span>
            <span class="text-[10px] uppercase font-bold text-brand-gold">Мебель</span>
          </div>
        </router-link>

        <!-- Desktop Menu -->
        <nav class="hidden md:flex items-center gap-10">
          <router-link 
            v-for="link in [
              { name: 'Проекты', to: '/catalog' },
              { name: 'Контакты', to: '/contact' }
            ]" 
            :key="link.to"
            :to="link.to" 
            :class="['font-bold text-xs uppercase hover:text-brand-gold transition-colors', isHeaderActive ? 'text-brand-brown' : 'text-white']"
          >
            {{ link.name }}
          </router-link>
        </nav>

        <!-- Icons -->
        <div class="flex items-center gap-4">
          <a href="tel:+79787631603" :class="['hidden lg:flex items-center gap-2 px-5 py-2.5 rounded-lg border transition-all font-bold text-xs uppercase', isHeaderActive ? 'border-brand-brown/10 text-brand-brown hover:bg-brand-brown hover:text-white' : 'border-white/20 text-white hover:bg-white hover:text-brand-brown']">
            <LucidePhone :size="14" />
            +7 (978) 763-16-03
          </a>
          <router-link to="/favorites" class="relative p-2 group">
            <LucideHeart :class="[isHeaderActive ? 'text-brand-brown' : 'text-white', 'group-hover:text-brand-gold transition-colors']" :size="24" />
            <span v-if="favorites.length" class="absolute -top-1 -right-1 w-5 h-5 bg-brand-gold text-white text-[10px] font-black rounded-lg flex items-center justify-center shadow-md animate-bounce">{{ favorites.length }}</span>
          </router-link>
          <button @click="isMenuOpen = true" class="md:hidden p-2">
            <LucideMenu :class="isHeaderActive ? 'text-brand-brown' : 'text-white'" :size="28" />
          </button>
        </div>
      </div>
    </header>

    <!-- Mobile Menu Overlay -->
    <Teleport to="body">
      <transition name="menu">
        <div v-if="isMenuOpen" class="fixed inset-0 z-[200] bg-brand-brown flex flex-col p-10">
          <div class="flex justify-between items-center mb-20">
            <span class="font-serif text-2xl text-white font-bold uppercase tracking-widest">Меню</span>
            <button @click="isMenuOpen = false" class="w-12 h-12 bg-white/5 rounded-lg flex items-center justify-center text-white">
              <LucideX :size="32" />
            </button>
          </div>
          
          <nav class="flex flex-col gap-8">
            <router-link 
              v-for="link in [
                { name: 'Главная', to: '/' },
                { name: 'Проекты', to: '/catalog' },
                { name: 'Контакты', to: '/contact' },
                { name: 'Избранное', to: '/favorites' }
              ]" 
              :key="link.to"
              :to="link.to" 
              @click="isMenuOpen = false"
              class="text-4xl font-serif text-white hover:text-brand-gold transition-colors"
            >
              {{ link.name }}
            </router-link>
          </nav>

          <div class="mt-auto pt-10 border-t border-white/10">
            <p class="text-white/40 uppercase text-xs mb-4 font-bold">Связаться с нами</p>
            <a href="tel:+79787631603" class="text-2xl text-white font-bold">+7 (978) 763-16-03</a>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Content -->
    <main class="flex-1">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="bg-brand-cream border-t border-brand-brown/5 pt-20 pb-12 px-6">
      <div class="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-4 gap-16">
        <div class="col-span-1 md:col-span-2">
          <div class="flex items-center gap-3 mb-8">
            <div class="w-10 h-10 bg-brand-brown text-brand-gold rounded-lg flex items-center justify-center font-serif text-xl font-bold">Р</div>
            <span class="font-serif text-2xl uppercase font-black text-brand-brown">РОСТ <span class="text-brand-gold">Мебель</span></span>
          </div>
          <p class="text-brand-brown/60 max-w-sm mb-10 leading-relaxed font-medium">
            Проектирование и изготовление премиальной мебели по индивидуальным размерам по всему Крыму.
          </p>
        </div>
        
        <div>
          <h4 class="font-bold text-xs uppercase text-brand-gold mb-8">Навигация</h4>
          <ul class="space-y-4">
            <li><router-link to="/catalog" class="font-bold text-brand-brown/80 hover:text-brand-gold transition-colors uppercase text-[10px]">Все проекты</router-link></li>
            <li><router-link to="/contact" class="font-bold text-brand-brown/80 hover:text-brand-gold transition-colors uppercase text-[10px]">Контакты</router-link></li>
            <li><router-link to="/favorites" class="font-bold text-brand-brown/80 hover:text-brand-gold transition-colors uppercase text-[10px]">Избранное</router-link></li>
          </ul>
        </div>

        <div>
          <h4 class="font-bold text-xs uppercase text-brand-gold mb-8">Контакты</h4>
          <ul class="space-y-4 font-bold text-brand-brown/80">
            <li class="flex items-center gap-2">
              <span class="w-1.5 h-1.5 bg-brand-gold rounded-full"></span>
              +7 (978) 763-16-03
            </li>
            <li class="flex items-center gap-2">
              <span class="w-1.5 h-1.5 bg-brand-gold rounded-full"></span>
              rost.salon2003@mail.ru
            </li>
            <li class="text-[10px] uppercase text-brand-brown/40 mt-4">г. Симферополь</li>
          </ul>
        </div>
      </div>
      
      <div class="max-w-7xl mx-auto mt-16 pt-8 border-t border-brand-brown/5 flex flex-col md:flex-row justify-between items-center gap-4">
        <p class="text-brand-brown/30 text-xs font-bold uppercase">© 2024 РОСТ Мебель. Все права защищены.</p>
        <div class="flex gap-8">
          <a href="#" class="text-brand-brown/30 text-[10px] uppercase hover:text-brand-gold transition-colors font-bold">Политика конфиденциальности</a>
        </div>
      </div>
    </footer>
  </div>
</template>

<style>
.menu-enter-active, .menu-leave-active { transition: all 0.6s cubic-bezier(0.85, 0, 0.15, 1); }
.menu-enter-from { transform: translateY(-100%); }
.menu-leave-to { transform: translateY(-100%); }

.selection\:bg-brand-gold ::selection { background-color: #3E8C76; color: white; }
</style>
