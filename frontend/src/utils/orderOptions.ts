export const projectTypeOptions = [
  'Кухня с техникой',
  'Кухня',
  'Шкаф или гардеробная',
  'Мебель для всей квартиры',
  'Коммерческий объект',
  'Пока не знаю',
] as const;

export const budgetOptions = [
  'До 200 000 ₽',
  '200 000-400 000 ₽',
  '400 000-700 000 ₽',
  'От 700 000 ₽',
  'Нужен расчет',
] as const;

export const contactOptions = [
  { value: 'phone', label: 'Звонок' },
  { value: 'whatsapp', label: 'WhatsApp' },
  { value: 'telegram', label: 'Telegram' },
  { value: 'max', label: 'MAX' },
  { value: 'email', label: 'Email' },
] as const;

export const contactMethodLabels = contactOptions.reduce<Record<string, string>>((acc, option) => {
  acc[option.value] = option.label;
  return acc;
}, {});
