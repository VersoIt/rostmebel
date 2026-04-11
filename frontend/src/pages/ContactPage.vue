<script setup lang="ts">
import OrderForm from '@/components/order/OrderForm.vue';
import { LucidePhone, LucideMail, LucideMapPin } from 'lucide-vue-next';
</script>

<template>
  <div class="bg-brand-cream min-h-screen pt-32 pb-24 px-4">
    <div class="max-w-7xl mx-auto">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-16">
        <div>
          <h1 class="font-serif text-5xl text-brand-brown mb-8">Контакты</h1>
          <p class="text-xl text-brand-brown/60 mb-12">Мы всегда рады обсудить ваши идеи и помочь с выбором идеальной мебели.</p>
          
          <div class="space-y-8">
            <div class="flex items-start gap-6">
              <div class="w-14 h-14 bg-brand-gold/10 rounded-2xl flex items-center justify-center text-brand-gold">
                <LucidePhone :size="28" />
              </div>
              <div>
                <h3 class="font-medium text-brand-brown/40 uppercase tracking-widest text-sm mb-1">Телефон</h3>
                <p class="text-2xl font-serif">+7 (978) 763-16-03</p>
              </div>
            </div>

            <div class="flex items-start gap-6">
              <div class="w-14 h-14 bg-brand-gold/10 rounded-2xl flex items-center justify-center text-brand-gold">
                <LucideMail :size="28" />
              </div>
              <div>
                <h3 class="font-medium text-brand-brown/40 uppercase tracking-widest text-sm mb-1">Email</h3>
                <p class="text-2xl font-serif">rost.salon2003@mail.ru</p>
              </div>
            </div>

            <div class="flex items-start gap-6">
              <div class="w-14 h-14 bg-brand-gold/10 rounded-2xl flex items-center justify-center text-brand-gold">
                <LucideMapPin :size="28" />
              </div>
              <div>
                <h3 class="font-medium text-brand-brown/40 uppercase tracking-widest text-sm mb-1">Адрес</h3>
                <p class="text-2xl font-serif">г. Симферополь</p>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white p-12 rounded-3xl shadow-xl border border-brand-brown/5">
          <h2 class="font-serif text-3xl mb-2">Напишите нам</h2>
          <p class="text-brand-brown/60 mb-8">Оставьте заявку, и мы перезвоним вам для консультации</p>
          <OrderForm />
        </div>
      </div>
    </div>
  </div>
</template>
