<script setup lang="ts">
import { EMAIL_DISPLAY, EMAIL_HREF, PHONE_DISPLAY, PHONE_HREF } from '@/constants/contacts';

const updatedAt = '28 сентября 2026 года';
const operatorName = 'Индивидуальный предприниматель Лященко Александр Юрьевич';

const personalData = [
  'имя',
  'номер телефона',
  'email, если пользователь указывает его в форме',
  'город или район',
  'комментарий к заявке',
  'тип проекта, ориентир по бюджету и предпочтительный способ связи',
  'технические данные: IP-адрес, cookie, данные браузера и сведения о действиях на сайте',
];

const purposes = [
  'обработка заявок и обратная связь с пользователем',
  'предварительный расчет мебели и консультация по проекту',
  'ведение истории обращений и улучшение качества обслуживания',
  'аналитика посещаемости сайта и улучшение интерфейса',
  'защита сайта от спама, злоупотреблений и технических атак',
];

const rights = [
  'получать информацию об обработке своих персональных данных',
  'требовать уточнения, блокирования или удаления персональных данных',
  'отозвать согласие на обработку персональных данных',
  'обжаловать действия оператора в уполномоченный орган или суд',
];
</script>

<template>
  <div class="min-h-screen bg-brand-cream pt-28">
    <section class="ui-container py-10 sm:py-14">
      <div class="mx-auto max-w-5xl">
        <div class="mb-8 rounded-[2rem] bg-brand-brown p-6 text-white shadow-[0_24px_70px_rgba(23,33,29,0.14)] sm:p-8">
          <p class="mb-3 text-xs font-black uppercase tracking-widest text-brand-gold">Документы сайта</p>
          <h1 class="font-serif text-4xl font-bold leading-tight sm:text-5xl">
            Политика обработки персональных данных
          </h1>
          <p class="mt-4 max-w-3xl leading-8 text-white/72">
            В этой политике описано, какие данные собирает сайт РОСТ Мебель, для чего они используются и как пользователь может обратиться по вопросам обработки данных.
          </p>
          <p class="mt-5 text-sm font-semibold text-white/52">Дата обновления: {{ updatedAt }}</p>
        </div>

        <div class="grid gap-5 lg:grid-cols-[250px_minmax(0,1fr)]">
          <aside class="self-start rounded-2xl border border-brand-brown/10 bg-white/80 p-5 text-sm font-semibold text-brand-brown/62 lg:sticky lg:top-28">
            <p class="mb-3 text-xs font-black uppercase tracking-widest text-brand-gold">Кратко</p>
            <p class="leading-7">
              Данные нужны для заявок, консультаций, связи с клиентом и аналитики сайта. Мы не публикуем персональные данные пользователей и не продаем их третьим лицам.
            </p>
          </aside>

          <article class="space-y-5">
            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">1. Общие положения</h2>
              <div class="space-y-3 leading-8 text-brand-brown/68">
                <p>
                  Настоящая Политика определяет порядок обработки персональных данных пользователей сайта
                  <span class="font-bold text-brand-brown">rostmebel.shop</span>.
                </p>
                <p>
                  Оператор персональных данных: <span class="font-bold text-brand-brown">{{ operatorName }}</span>. Связаться с оператором можно по телефону
                  <a :href="PHONE_HREF" class="font-bold text-brand-gold hover:text-brand-brown">{{ PHONE_DISPLAY }}</a>
                  или по email
                  <a :href="EMAIL_HREF" class="font-bold text-brand-gold hover:text-brand-brown">{{ EMAIL_DISPLAY }}</a>.
                </p>
                <p>
                  Политика подготовлена с учетом Федерального закона РФ № 152-ФЗ «О персональных данных» и применяется ко всем данным, которые пользователь передает через формы сайта или которые автоматически собираются при посещении сайта.
                </p>
              </div>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">2. Какие данные обрабатываются</h2>
              <ul class="grid gap-2 leading-7 text-brand-brown/68 sm:grid-cols-2">
                <li v-for="item in personalData" :key="item" class="rounded-2xl bg-brand-gray/70 px-4 py-3">
                  {{ item }}
                </li>
              </ul>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">3. Цели обработки данных</h2>
              <ul class="space-y-3 leading-8 text-brand-brown/68">
                <li v-for="item in purposes" :key="item" class="flex gap-3">
                  <span class="mt-3 h-1.5 w-1.5 shrink-0 rounded-full bg-brand-gold"></span>
                  <span>{{ item }}</span>
                </li>
              </ul>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">4. Правовые основания</h2>
              <p class="leading-8 text-brand-brown/68">
                Обработка данных осуществляется на основании согласия пользователя, которое пользователь выражает при отправке формы на сайте, а также в случаях, когда обработка необходима для ответа на обращение, подготовки расчета, исполнения договоренностей с пользователем или соблюдения требований закона.
              </p>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">5. Передача третьим лицам</h2>
              <div class="space-y-3 leading-8 text-brand-brown/68">
                <p>
                  Персональные данные могут передаваться сервисам, которые помогают работе сайта и обработке заявок: хостинг-провайдеру, сервисам аналитики, мессенджерам, почтовым и техническим сервисам.
                </p>
                <p>
                  Передача выполняется только в объеме, необходимом для работы сайта, связи с пользователем, защиты от спама и анализа посещаемости. Персональные данные не продаются и не публикуются.
                </p>
              </div>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">6. Cookie и аналитика</h2>
              <p class="leading-8 text-brand-brown/68">
                Сайт может использовать cookie и сервисы веб-аналитики, включая Яндекс Метрику, чтобы понимать посещаемость, источники переходов, клики и технические ошибки. Пользователь может ограничить использование cookie в настройках браузера, но часть функций сайта может работать некорректно.
              </p>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">7. Срок хранения</h2>
              <p class="leading-8 text-brand-brown/68">
                Данные хранятся столько, сколько необходимо для обработки заявки, связи с пользователем, ведения истории обращений, выполнения договоренностей и соблюдения требований закона. После достижения целей обработки данные удаляются или обезличиваются, если закон не требует их дальнейшего хранения.
              </p>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">8. Права пользователя</h2>
              <ul class="space-y-3 leading-8 text-brand-brown/68">
                <li v-for="item in rights" :key="item" class="flex gap-3">
                  <span class="mt-3 h-1.5 w-1.5 shrink-0 rounded-full bg-brand-gold"></span>
                  <span>{{ item }}</span>
                </li>
              </ul>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">9. Защита данных</h2>
              <p class="leading-8 text-brand-brown/68">
                Оператор принимает организационные и технические меры для защиты персональных данных от неправомерного доступа, изменения, раскрытия, блокирования, копирования, распространения и иных неправомерных действий.
              </p>
            </section>

            <section class="ui-card p-5 sm:p-7">
              <h2 class="ui-title-md mb-4">10. Обращения по персональным данным</h2>
              <p class="leading-8 text-brand-brown/68">
                Чтобы задать вопрос, уточнить данные, отозвать согласие или запросить удаление данных, напишите на
                <a :href="EMAIL_HREF" class="font-bold text-brand-gold hover:text-brand-brown">{{ EMAIL_DISPLAY }}</a>
                или позвоните по номеру
                <a :href="PHONE_HREF" class="font-bold text-brand-gold hover:text-brand-brown">{{ PHONE_DISPLAY }}</a>.
              </p>
            </section>
          </article>
        </div>
      </div>
    </section>
  </div>
</template>
