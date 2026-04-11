<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { 
  LucidePhone, 
  LucideMail, 
  LucideMapPin, 
  LucideMessageCircle, 
  LucideSend,
  LucidePlus,
  LucideMinus,
  LucideClock,
  LucideInstagram,
  LucideCompass
} from 'lucide-vue-next';
import OrderForm from '@/components/order/OrderForm.vue';
import { gsap } from 'gsap';

const faqs = ref([
  { 
    q: 'Сколько времени занимает изготовление кухни?', 
    a: 'В среднем от 35 до 60 рабочих дней. Каждый проект уникален: мы заказываем фурнитуру и материалы под ваши размеры напрямую у производителей в Европе и России.',
    isOpen: false 
  },
  { 
    q: 'Выезжаете ли вы на замер по Крыму?', 
    a: 'Безусловно. Наш инженер-дизайнер выезжает в любую точку полуострова: Севастополь, Симферополь, Ялта, Евпатория. При заключении договора замер — в подарок.',
    isOpen: false 
  },
  { 
    q: 'Предоставляете ли вы технические чертежи строителям?', 
    a: 'Да, это наше преимущество. Мы выдаем полные карты разводки электрики и сантехники, а также схему раскладки плитки, чтобы ваш ремонт прошел без ошибок.',
    isOpen: false 
  }
]);

const toggleFaq = (idx: number) => {
  faqs.value[idx].isOpen = !faqs.value[idx].isOpen;
};

onMounted(() => {
  gsap.from('.reveal', {
    y: 40,
    opacity: 0,
    duration: 1,
    stagger: 0.2,
    ease: 'power4.out'
  });
});
</script>

<template>
  <div class="bg-white min-h-screen pt-40 pb-32 overflow-hidden selection:bg-brand-gold selection:text-white">
    <div class="max-w-7xl mx-auto px-6">
      
      <!-- Section 1: Header -->
      <header class="mb-32 relative">
        <div class="absolute -top-20 -left-20 w-64 h-64 bg-brand-gold/5 rounded-full blur-3xl animate-pulse"></div>
        <span class="reveal text-brand-gold font-bold text-xs uppercase tracking-[0.5em] mb-6 block">Communication</span>
        <h1 class="reveal font-serif text-6xl md:text-9xl text-brand-brown leading-[0.9] tracking-tighter mb-8">
          Давайте создадим <br> 
          <span class="text-brand-gold italic ml-12 md:ml-32">вашу историю</span>
        </h1>
        <p class="reveal text-brand-brown/40 text-lg md:text-xl max-w-2xl font-light leading-relaxed">
          Мы всегда открыты для новых амбициозных проектов. Запишитесь на встречу в нашем шоуруме или пригласите дизайнера на объект.
        </p>
      </header>

      <div class="grid grid-cols-1 lg:grid-cols-12 gap-16 items-start mb-40">
        <!-- Section 2: Contact Methods -->
        <div class="lg:col-span-7 space-y-12">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <a href="tel:+79787631603" class="reveal group p-10 bg-brand-cream/20 rounded-[2rem] border border-brand-brown/5 hover:border-brand-gold/20 transition-all duration-500 hover:shadow-2xl hover:shadow-brand-gold/5">
              <div class="w-14 h-14 bg-white rounded-2xl flex items-center justify-center text-brand-gold shadow-sm mb-8 group-hover:scale-110 transition-transform">
                <LucidePhone :size="24" />
              </div>
              <h4 class="font-serif text-2xl text-brand-brown mb-2">Позвонить</h4>
              <p class="text-brand-brown/40 text-sm mb-6 uppercase tracking-widest font-bold">Ежедневно с 10:00 до 19:00</p>
              <span class="text-2xl font-medium text-brand-brown border-b-2 border-brand-gold/20 group-hover:border-brand-gold transition-colors">+7 (978) 763-16-03</span>
            </a>

            <a href="mailto:info@rostmebel.ru" class="reveal group p-10 bg-brand-cream/20 rounded-[2rem] border border-brand-brown/5 hover:border-brand-gold/20 transition-all duration-500 hover:shadow-2xl hover:shadow-brand-gold/5">
              <div class="w-14 h-14 bg-white rounded-2xl flex items-center justify-center text-brand-gold shadow-sm mb-8 group-hover:scale-110 transition-transform">
                <LucideMail :size="24" />
              </div>
              <h4 class="font-serif text-2xl text-brand-brown mb-2">Написать</h4>
              <p class="text-brand-brown/40 text-sm mb-6 uppercase tracking-widest font-bold">Для ТЗ и предложений</p>
              <span class="text-2xl font-medium text-brand-brown border-b-2 border-brand-gold/20 group-hover:border-brand-gold transition-colors">info@rostmebel.ru</span>
            </a>
          </div>

          <!-- Showroom Dark Card -->
          <div class="reveal p-12 bg-brand-brown text-white rounded-[3rem] shadow-2xl relative overflow-hidden group">
            <div class="absolute top-0 right-0 w-96 h-96 bg-brand-gold/10 rounded-full blur-[100px] -translate-y-1/2 translate-x-1/2 group-hover:scale-125 transition-transform duration-1000"></div>
            
            <div class="relative z-10 grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
              <div>
                <LucideCompass class="text-brand-gold mb-8" :size="48" />
                <h4 class="font-serif text-4xl mb-6 leading-tight">Посетите наш <br> шоурум</h4>
                <div class="space-y-4 text-white/60 mb-10 text-lg">
                  <p class="flex items-center gap-3">
                    <LucideMapPin :size="20" class="text-brand-gold" />
                    Севастополь, Фиолентовское шоссе, 1/2
                  </p>
                  <p class="flex items-center gap-3">
                    <LucideClock :size="20" class="text-brand-gold" />
                    Пн-Сб: 10:00 — 19:00
                  </p>
                </div>
                <div class="flex gap-4">
                  <a href="https://wa.me/79787631603" class="flex-1 bg-white/5 hover:bg-green-600/20 border border-white/10 text-white py-4 rounded-2xl flex items-center justify-center gap-2 font-bold transition-all backdrop-blur-md">
                    <LucideMessageCircle :size="20" class="text-green-500" /> WhatsApp
                  </a>
                  <a href="https://t.me/rostmebel" class="flex-1 bg-white/5 hover:bg-blue-500/20 border border-white/10 text-white py-4 rounded-2xl flex items-center justify-center gap-2 font-bold transition-all backdrop-blur-md">
                    <LucideSend :size="20" class="text-blue-400" /> Telegram
                  </a>
                </div>
              </div>
              <div class="hidden md:block">
                <div class="aspect-square rounded-3xl overflow-hidden grayscale opacity-50 group-hover:grayscale-0 group-hover:opacity-100 transition-all duration-700">
                  <img src="https://images.unsplash.com/photo-1618221195710-dd6b41faaea6?q=80&w=800&auto=format&fit=crop" class="w-full h-full object-cover">
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Section 3: Consultation Form -->
        <div class="lg:col-span-5 sticky top-32">
          <div class="reveal bg-white p-12 rounded-[3.5rem] shadow-[0_50px_100px_-20px_rgba(44,24,16,0.15)] border border-brand-brown/5 relative">
            <div class="absolute -top-6 left-12 bg-brand-gold text-brand-brown px-6 py-2 rounded-full text-[10px] font-black uppercase tracking-widest shadow-xl">
              Personal Meeting
            </div>
            <h3 class="font-serif text-4xl mb-3 text-brand-brown">Консультация</h3>
            <p class="text-brand-brown/40 mb-10 font-medium">Оставьте заявку на встречу с ведущим дизайнером ателье</p>
            <OrderForm />
            
            <p class="mt-8 text-[10px] text-center text-brand-brown/30 uppercase tracking-widest font-bold">
              Нажимая кнопку, вы соглашаетесь с политикой конфиденциальности
            </p>
          </div>
        </div>
      </div>

      <!-- Section 4: FAQ -->
      <div class="reveal max-w-5xl mx-auto py-32 border-t border-brand-brown/5">
        <div class="text-center mb-20">
          <span class="text-brand-gold font-bold text-xs uppercase tracking-[0.5em] mb-4 block">Information</span>
          <h2 class="font-serif text-5xl text-brand-brown">Вопросы и ответы</h2>
        </div>
        
        <div class="grid grid-cols-1 gap-6">
          <div 
            v-for="(faq, idx) in faqs" :key="idx"
            class="group transition-all duration-500"
          >
            <button 
              @click="toggleFaq(idx)"
              :class="[
                'w-full p-10 text-left flex items-center justify-between rounded-3xl transition-all duration-500 border',
                faq.isOpen ? 'bg-brand-brown text-white border-brand-brown shadow-2xl' : 'bg-white text-brand-brown border-brand-brown/5 hover:border-brand-gold/30 shadow-sm'
              ]"
            >
              <span class="font-serif text-2xl" :class="{ 'text-brand-gold': faq.isOpen }">{{ faq.q }}</span>
              <div :class="['w-10 h-10 rounded-full flex items-center justify-center transition-all duration-500', faq.isOpen ? 'bg-brand-gold text-brand-brown rotate-180' : 'bg-brand-cream text-brand-gold']">
                <LucidePlus v-if="!faq.isOpen" :size="20" />
                <LucideMinus v-else :size="20" />
              </div>
            </button>
            <transition name="faq-expand">
              <div v-if="faq.isOpen" class="px-10 py-8 text-white/60 leading-relaxed text-lg font-light bg-brand-brown rounded-b-[2rem] -mt-4 pt-12 relative z-[-1]">
                {{ faq.a }}
              </div>
            </transition>
          </div>
        </div>
      </div>

      <!-- Final Tagline -->
      <div class="reveal py-20 text-center relative">
        <div class="font-serif text-[15vw] leading-none text-brand-cream/50 absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 -z-10 select-none whitespace-nowrap">
          ROST MEBEL
        </div>
        <h2 class="font-serif text-4xl md:text-6xl text-brand-brown mb-8 relative z-10 italic">
          Превращаем мечты в <br> <span class="text-brand-gold not-italic">совершенную форму.</span>
        </h2>
      </div>

    </div>
  </div>
</template>

<style scoped>
.faq-expand-enter-active, .faq-expand-leave-active { 
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1); 
  max-height: 500px;
  overflow: hidden;
}
.faq-expand-enter-from, .faq-expand-leave-to { 
  opacity: 0; 
  max-height: 0; 
  transform: translateY(-20px); 
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: scale(1); }
  50% { opacity: 0.8; transform: scale(1.1); }
}

.animate-pulse {
  animation: pulse 8s infinite ease-in-out;
}
</style>
