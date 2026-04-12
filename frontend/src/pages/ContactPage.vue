<script setup lang="ts">
import { ref } from 'vue';
import {
  LucideClock,
  LucideMail,
  LucideMapPin,
  LucideMessageCircle,
  LucideMinus,
  LucidePhone,
  LucidePlus,
  LucideSend,
} from 'lucide-vue-next';
import OrderForm from '@/components/order/OrderForm.vue';

const faqs = ref([
  {
    q: 'Сколько занимает изготовление мебели?',
    a: 'Обычно 35-60 рабочих дней после согласования проекта и материалов. Точный срок зависит от фурнитуры, фасадов и объема монтажа.',
    isOpen: false,
  },
  {
    q: 'Можно начать с планировки без выезда?',
    a: 'Да. Пришлите размеры, фото помещения и список техники. Мы быстро оценим конфигурацию, бюджет и спорные места.',
    isOpen: false,
  },
  {
    q: 'Вы даете чертежи для ремонта?',
    a: 'Да. Для кухни готовим карты розеток, выводов воды, подсветки и примыканий, чтобы строители подготовили помещение под мебель.',
    isOpen: false,
  },
]);

const toggleFaq = (idx: number) => {
  faqs.value[idx].isOpen = !faqs.value[idx].isOpen;
};
</script>

<template>
  <div class="bg-brand-cream min-h-screen pt-28 pb-16 text-brand-brown selection:bg-brand-gold selection:text-white">
    <div class="mx-auto max-w-7xl px-5">
      <header class="grid grid-cols-1 gap-8 border-b border-brand-brown/10 pb-10 lg:grid-cols-[1.1fr_0.9fr] lg:items-end">
        <div>
          <p class="mb-3 text-sm font-semibold uppercase text-brand-gold">Контакты</p>
          <h1 class="font-serif text-4xl leading-tight md:text-6xl">Обсудим мебель без долгих презентаций</h1>
          <p class="mt-5 max-w-2xl text-lg leading-8 text-brand-brown/68">
            Расскажите, что нужно сделать: кухня, шкафы, гардеробная или мебель для всего объекта. Мы зададим точные вопросы и предложим следующий шаг.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <a href="tel:+79787631603" class="border border-brand-brown/10 bg-white p-5 transition hover:border-brand-gold">
            <LucidePhone class="mb-4 text-brand-gold" :size="24" />
            <div class="text-sm text-brand-brown/55">Телефон</div>
            <div class="mt-1 text-xl font-semibold">+7 (978) 763-16-03</div>
          </a>
          <a href="mailto:rost.salon2003@mail.ru" class="border border-brand-brown/10 bg-white p-5 transition hover:border-brand-gold">
            <LucideMail class="mb-4 text-brand-gold" :size="24" />
            <div class="text-sm text-brand-brown/55">Почта</div>
            <div class="mt-1 break-all text-xl font-semibold">rost.salon2003@mail.ru</div>
          </a>
        </div>
      </header>

      <section class="grid grid-cols-1 gap-8 py-10 lg:grid-cols-[0.9fr_1.1fr] lg:items-start">
        <div class="space-y-4">
          <div class="border border-brand-brown/10 bg-white p-6">
            <h2 class="font-serif text-3xl">Как быстрее получить расчет</h2>
            <div class="mt-6 space-y-4">
              <div class="flex gap-3">
                <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-brand-brown text-white">1</div>
                <p class="leading-7 text-brand-brown/68">Пришлите фото помещения, примерные размеры и список техники.</p>
              </div>
              <div class="flex gap-3">
                <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-brand-brown text-white">2</div>
                <p class="leading-7 text-brand-brown/68">Напишите, что важно: бюджет, сроки, фасады, хранение, встроенная техника.</p>
              </div>
              <div class="flex gap-3">
                <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-brand-brown text-white">3</div>
                <p class="leading-7 text-brand-brown/68">Мы предложим конфигурацию и скажем, где нужен точный замер.</p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <a href="https://wa.me/79787631603" class="flex items-center justify-between border border-brand-brown/10 bg-white p-5 transition hover:border-brand-gold">
              <span>
                <span class="block text-sm text-brand-brown/55">WhatsApp</span>
                <span class="mt-1 block font-semibold">Написать</span>
              </span>
              <LucideMessageCircle class="text-brand-gold" :size="24" />
            </a>
            <a href="https://t.me/rostmebel" class="flex items-center justify-between border border-brand-brown/10 bg-white p-5 transition hover:border-brand-gold">
              <span>
                <span class="block text-sm text-brand-brown/55">Telegram</span>
                <span class="mt-1 block font-semibold">Открыть чат</span>
              </span>
              <LucideSend class="text-brand-gold" :size="24" />
            </a>
          </div>

          <div class="overflow-hidden border border-brand-brown/10 bg-white">
            <img src="/assets/images/showroom.jpg" class="h-64 w-full object-cover" alt="Шоурум РОСТ Мебель">
            <div class="grid grid-cols-1 gap-4 p-6 sm:grid-cols-2">
              <p class="flex items-center gap-3 text-brand-brown/70">
                <LucideMapPin class="text-brand-gold" :size="20" />
                Симферополь, выезд по Крыму
              </p>
              <p class="flex items-center gap-3 text-brand-brown/70">
                <LucideClock class="text-brand-gold" :size="20" />
                Пн-Сб: 10:00-19:00
              </p>
            </div>
          </div>
        </div>

        <div class="lg:sticky lg:top-28">
          <div class="border border-brand-brown/10 bg-white p-6 md:p-8">
            <p class="mb-3 text-sm font-semibold uppercase text-brand-gold">Заявка</p>
            <h2 class="font-serif text-3xl">Получить консультацию</h2>
            <p class="mt-3 leading-7 text-brand-brown/60">
              Ответим, какой формат работы подойдет, и подскажем ориентир по бюджету до встречи.
            </p>
            <div class="mt-6">
              <OrderForm />
            </div>
          </div>
        </div>
      </section>

      <section class="border-t border-brand-brown/10 pt-10">
        <div class="grid grid-cols-1 gap-8 lg:grid-cols-[0.8fr_1.2fr]">
          <div>
            <p class="mb-3 text-sm font-semibold uppercase text-brand-gold">FAQ</p>
            <h2 class="font-serif text-4xl">Частые вопросы</h2>
          </div>

          <div class="space-y-3">
            <div v-for="(faq, idx) in faqs" :key="faq.q" class="border border-brand-brown/10 bg-white">
              <button
                @click="toggleFaq(idx)"
                class="flex w-full items-center justify-between gap-4 p-5 text-left"
              >
                <span class="font-serif text-2xl">{{ faq.q }}</span>
                <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-brand-cream text-brand-gold">
                  <LucideMinus v-if="faq.isOpen" :size="18" />
                  <LucidePlus v-else :size="18" />
                </span>
              </button>
              <div v-if="faq.isOpen" class="border-t border-brand-brown/10 px-5 pb-5 pt-4 leading-7 text-brand-brown/68">
                {{ faq.a }}
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
