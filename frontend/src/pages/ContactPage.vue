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
import QuoteQuiz from '@/components/order/QuoteQuiz.vue';

const maxMessengerUrl = 'https://max.ru';

const faqs = ref([
  {
    q: 'Сколько занимает изготовление мебели?',
    a: 'Обычно 35-60 рабочих дней после согласования проекта и материалов. Точный срок зависит от фурнитуры, фасадов и объема монтажа.',
    isOpen: false,
  },
  {
    q: 'Можно начать с планировки без выезда?',
    a: 'Да. Пришлите размеры, фото помещения и список техники. Мы оценим конфигурацию, бюджет и места, которые нужно проверить на замере.',
    isOpen: false,
  },
  {
    q: 'Вы даете чертежи для ремонта?',
    a: 'Да. Для кухни готовим карты розеток, выводов воды, подсветки и примыканий, чтобы строители подготовили помещение под мебель.',
    isOpen: false,
  },
  {
    q: 'Можно заказать кухню вместе с техникой?',
    a: 'Да. Можем заложить технику в проект и смету: варочную поверхность, духовой шкаф, вытяжку, посудомойку, холодильник, мойку и смеситель. Так проще сразу проверить размеры, розетки и посадочные места.',
    isOpen: false,
  },
]);

const toggleFaq = (idx: number) => {
  faqs.value[idx].isOpen = !faqs.value[idx].isOpen;
};
</script>

<template>
  <div class="min-h-screen bg-brand-cream pt-28 text-brand-brown">
    <div class="ui-container ui-section">
      <header class="grid grid-cols-1 gap-8 border-b border-brand-brown/10 pb-10 lg:grid-cols-[1.1fr_0.9fr] lg:items-end">
        <div>
          <p class="ui-eyebrow mb-3">Контакты</p>
          <h1 class="ui-title-xl">Обсудим мебель без длинных презентаций</h1>
          <p class="ui-copy-lg mt-5 max-w-2xl">
            Расскажите, что нужно сделать: кухня, шкафы, гардеробная или мебель для всего объекта. Мы зададим точные вопросы и предложим следующий шаг.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <a href="tel:+79787631603" class="ui-card ui-card-hover p-5">
            <LucidePhone class="mb-4 text-brand-gold" :size="24" />
            <div class="text-sm text-brand-brown/55">Телефон</div>
            <div class="mt-1 text-xl font-semibold">+7 (978) 763-16-03</div>
          </a>
          <a href="mailto:rost.salon2003@mail.ru" class="ui-card ui-card-hover p-5">
            <LucideMail class="mb-4 text-brand-gold" :size="24" />
            <div class="text-sm text-brand-brown/55">Почта</div>
            <div class="mt-1 break-all text-xl font-semibold">rost.salon2003@mail.ru</div>
          </a>
        </div>
      </header>

      <section class="grid grid-cols-1 gap-8 py-10 lg:grid-cols-[0.9fr_1.1fr] lg:items-start">
        <div class="space-y-4">
          <div class="ui-card p-5 sm:p-6">
            <h2 class="ui-title-md">Как быстрее получить расчет</h2>
            <div class="mt-6 space-y-4">
              <div v-for="(text, idx) in [
                'Пришлите фото помещения, примерные размеры и список техники.',
                'Напишите, что важно: бюджет, сроки, фасады, хранение, встроенная техника.',
                'Мы предложим конфигурацию и скажем, где нужен точный замер.'
              ]" :key="text" class="flex gap-3">
                <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-brand-brown text-sm font-bold text-white">{{ idx + 1 }}</div>
                <p class="ui-copy">{{ text }}</p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <a href="https://wa.me/79787631603" class="ui-card ui-card-hover flex items-center justify-between p-5">
              <span>
                <span class="block text-sm text-brand-brown/55">WhatsApp</span>
                <span class="mt-1 block font-semibold">Написать</span>
              </span>
              <LucideMessageCircle class="text-brand-gold" :size="24" />
            </a>
            <a href="https://t.me/rostmebel" class="ui-card ui-card-hover flex items-center justify-between p-5">
              <span>
                <span class="block text-sm text-brand-brown/55">Telegram</span>
                <span class="mt-1 block font-semibold">Открыть чат</span>
              </span>
              <LucideSend class="text-brand-gold" :size="24" />
            </a>
            <a :href="maxMessengerUrl" target="_blank" rel="noopener noreferrer" class="ui-card ui-card-hover flex items-center justify-between p-5">
              <span>
                <span class="block text-sm text-brand-brown/55">MAX</span>
                <span class="mt-1 block font-semibold">Открыть</span>
              </span>
              <LucideMessageCircle class="text-brand-gold" :size="24" />
            </a>
          </div>

          <div class="ui-card overflow-hidden">
            <img src="/assets/images/showroom.jpg" class="h-64 w-full object-cover" alt="Шоурум РОСТ Мебель">
            <div class="grid grid-cols-1 gap-4 p-5 sm:grid-cols-2">
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

        <div id="quote-quiz" class="scroll-mt-28 lg:sticky lg:top-28">
          <div class="ui-card p-5 sm:p-8">
            <p class="ui-eyebrow mb-3">Быстрый старт</p>
            <h2 class="ui-title-md">Рассчитать проект</h2>
            <p class="ui-copy mt-3">
              Четыре коротких шага: поймем задачу, бюджет, город и удобный способ связи.
            </p>
            <div class="mt-6">
              <QuoteQuiz />
            </div>
          </div>
        </div>
      </section>

      <section class="border-t border-brand-brown/10 pt-10">
        <div class="grid grid-cols-1 gap-8 lg:grid-cols-[0.8fr_1.2fr]">
          <div>
            <p class="ui-eyebrow mb-3">FAQ</p>
            <h2 class="ui-title-lg">Частые вопросы</h2>
          </div>

          <div class="space-y-3">
            <div v-for="(faq, idx) in faqs" :key="faq.q" class="ui-card overflow-hidden">
              <button type="button" class="flex w-full items-center justify-between gap-4 p-5 text-left" @click="toggleFaq(idx)">
                <span class="font-serif text-xl font-bold sm:text-2xl">{{ faq.q }}</span>
                <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-brand-gray text-brand-gold">
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
