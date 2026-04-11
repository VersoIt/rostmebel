<script setup lang="ts">
import { ref } from 'vue';
import { 
  LucidePhone, 
  LucideMail, 
  LucideMapPin, 
  LucideMessageCircle, 
  LucideSend,
  LucidePlus,
  LucideMinus
} from 'lucide-vue-next';
import OrderForm from '@/components/order/OrderForm.vue';

const faqs = ref([
  { 
    q: 'Сколько времени занимает изготовление кухни?', 
    a: 'В среднем от 35 до 60 рабочих дней, в зависимости от сложности изделия (эмаль, шпон), материалов и эксклюзивной фурнитуры.',
    isOpen: false 
  },
  { 
    q: 'Выезжаете ли вы на замер по Крыму?', 
    a: 'Да, наш дизайнер-замерщик выезжает в любую точку Крыма. При заключении договора замер и проект — бесплатно.',
    isOpen: false 
  },
  { 
    q: 'Какую фурнитуру вы используете?', 
    a: 'Мы работаем только с проверенными брендами: Blum (Австрия), Hettich (Германия) и Grass. Это гарантирует плавность хода и долговечность.',
    isOpen: false 
  }
]);

const toggleFaq = (idx: number) => {
  faqs.value[idx].isOpen = !faqs.value[idx].isOpen;
};
</script>

<template>
  <div class="bg-white min-h-screen pt-32 pb-24">
    <div class="max-w-7xl mx-auto px-6">
      <header class="mb-20">
        <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.3em] mb-4 block text-center">Связь с нами</span>
        <h1 class="font-serif text-5xl md:text-7xl text-brand-brown text-center">Давайте обсудим <br> ваш проект</h1>
      </header>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-24 items-start mb-40">
        <!-- Contact Info -->
        <div class="space-y-12">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-10">
            <div class="p-8 bg-brand-cream/30 rounded-3xl border border-brand-brown/5">
              <LucidePhone class="text-brand-gold mb-6" :size="32" />
              <h4 class="font-bold text-lg mb-2">Позвонить</h4>
              <a href="tel:+79787631603" class="text-xl font-medium text-brand-brown hover:text-brand-gold transition-colors">+7 (978) 763-16-03</a>
            </div>
            <div class="p-8 bg-brand-cream/30 rounded-3xl border border-brand-brown/5">
              <LucideMail class="text-brand-gold mb-6" :size="32" />
              <h4 class="font-bold text-lg mb-2">Написать</h4>
              <a href="mailto:info@rostmebel.ru" class="text-xl font-medium text-brand-brown hover:text-brand-gold transition-colors">info@rostmebel.ru</a>
            </div>
          </div>

          <div class="p-10 bg-brand-brown text-white rounded-[2.5rem] shadow-2xl relative overflow-hidden">
            <div class="absolute top-0 right-0 w-32 h-32 bg-white/5 rounded-full blur-3xl"></div>
            <LucideMapPin class="text-brand-gold mb-6" :size="32" />
            <h4 class="font-serif text-3xl mb-4">Наш офис и шоурум</h4>
            <p class="text-white/60 text-lg leading-relaxed mb-8">
              г. Севастополь, ул. Фиолентовское шоссе, 1/2 <br>
              Работаем Пн-Сб с 10:00 до 19:00
            </p>
            <div class="flex gap-4">
              <a href="https://wa.me/79787631603" class="flex-1 bg-green-600 hover:bg-green-700 text-white py-4 rounded-2xl flex items-center justify-center gap-2 font-bold transition-all">
                <LucideMessageCircle :size="20" /> WhatsApp
              </a>
              <a href="https://t.me/rostmebel" class="flex-1 bg-blue-500 hover:bg-blue-600 text-white py-4 rounded-2xl flex items-center justify-center gap-2 font-bold transition-all">
                <LucideSend :size="20" /> Telegram
              </a>
            </div>
          </div>
        </div>

        <!-- Order Form -->
        <div class="bg-white p-12 rounded-[3rem] shadow-2xl border border-brand-brown/5">
          <h3 class="font-serif text-3xl mb-2 text-brand-brown">Оставить заявку</h3>
          <p class="text-brand-brown/60 mb-10">Мы перезвоним вам в течение 15 минут</p>
          <OrderForm />
        </div>
      </div>

      <!-- FAQ Section -->
      <div class="max-w-4xl mx-auto">
        <h2 class="font-serif text-4xl text-brand-brown mb-12 text-center">Часто задаваемые вопросы</h2>
        <div class="space-y-4">
          <div 
            v-for="(faq, idx) in faqs" :key="idx"
            class="border border-brand-brown/10 rounded-2xl overflow-hidden"
          >
            <button 
              @click="toggleFaq(idx)"
              class="w-full p-6 text-left flex items-center justify-between hover:bg-brand-cream/20 transition-colors"
            >
              <span class="font-bold text-lg text-brand-brown">{{ faq.q }}</span>
              <component :is="faq.isOpen ? LucideMinus : LucidePlus" class="text-brand-gold" :size="20" />
            </button>
            <transition name="faq">
              <div v-if="faq.isOpen" class="p-6 pt-0 text-brand-brown/60 leading-relaxed border-t border-brand-brown/5 bg-brand-cream/5">
                {{ faq.a }}
              </div>
            </transition>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.faq-enter-active, .faq-leave-active { transition: all 0.3s ease; max-height: 200px; }
.faq-enter-from, .faq-leave-to { opacity: 0; max-height: 0; transform: translateY(-10px); }
</style>
